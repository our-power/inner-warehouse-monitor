package main

import (
	"html/template"
	"archive/zip"
	"net/http"
	"net/url"
	"path"
	"strings"
	"strconv"
	"utils"
	"time"
	"fmt"
	"log"
	"os"
	"io"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

var logger *log.Logger
var db *sql.DB

func dealRequest(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	action := r.FormValue("action")
	switch action {
	case "get_list":
		ip := r.FormValue("ip")
		getList(ip, w)
	case "get_file":
		filename := r.FormValue("name")
		version := r.FormValue("v")
		getFile(filename, version, w)
	case "set_done":
		ip := r.FormValue("ip")
		setDone(ip, w)
	}
}

func getList(ip string, w http.ResponseWriter) {
	rows, err := db.Query("select v,files from version where v=(select v from machine where ip=? and done=0)", ip)
	if err != nil {
		logger.Println("cannot query table version")
		fmt.Fprintf(w, "")
		return
	}
	for rows.Next() {
		var v, files string
		err = rows.Scan(&v, &files)
		if err == nil {
			fmt.Fprintf(w, "%s;%s", v, files)
			return
		}
	}
	fmt.Fprintf(w, "")
	return
}

func getFile(filename string, version string, w http.ResponseWriter) {
	name, _ := url.QueryUnescape(filename)
	_, err := strconv.Atoi(version)
	if strings.Contains(name, "..") || err != nil {
		logger.Println("invalid query.")
		fmt.Fprintf(w, "")
		return
	}
	outFile := "../up/" + version + "/" + name
	f, err := os.Open(outFile)
	if nil != err && !os.IsExist(err) {
		logger.Println(err.Error())
		fmt.Fprintf(w, "")
		return
	}
	defer f.Close()
	io.Copy(w, f)
}

func setDone(ip string, w http.ResponseWriter) {
	_, err := db.Exec("update machine set done=1 where ip=? and done=0", ip)
	if err != nil {
		logger.Println("cannot update table machine")
		fmt.Fprintf(w, "")
		return
	}
	fmt.Fprintf(w, "success")
}

type FilesAndVersion struct {
	Files   string
	Version int
}

func uploadFile(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		t, err := template.ParseFiles("../templates/upload.html")
		if err != nil {
			fmt.Println("Template not found.")
		} else {
			t.Execute(w, "")
		}
	} else {
		r.ParseMultipartForm(32<<20)
		file, handler, err := r.FormFile("uploadfile")
		if err != nil {
			fmt.Println(err)
			return
		}
		defer file.Close()
		f, err := os.OpenFile("../up/" + handler.Filename, os.O_WRONLY | os.O_CREATE, 0666)
		if err != nil {
			fmt.Fprintf(w, err.Error())
			return
		}
		defer f.Close()
		io.Copy(f, file)
		v := getNextVersion()
		files, err := unzipAndMove("../up/" + handler.Filename, v)
		if err != nil {
			logger.Println(err)
			fmt.Fprintf(w, err.Error())
			return
		}
		fv := FilesAndVersion {files, v}
		t, err := template.ParseFiles("../templates/iplist.html")
		if err != nil {
			fmt.Println("Template not found.")
		} else {
			t.Execute(w, fv)
		}
	}
}

func store(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		files := r.FormValue("files")
		version := r.FormValue("version")
		ips := r.FormValue("ips")
		ips = strings.Replace(ips, "\r\n", "\n", -1)
		ips = strings.Replace(ips, "\r", "\n", -1)
		ipList := strings.Split(ips, "\n")
		_, err := db.Exec("replace into version (v, files, time) values (?,?,?)", version, files, time.Now().Format("2006-01-02 15:04:05"))
		if err != nil {
			logger.Println("cannot update table version", err)
			fmt.Fprintf(w, "failed")
			return
		}
		for _, ip := range ipList {
			_, err := db.Exec("replace into machine (ip, v, done) values (?,?,0)", ip, version)
			if err != nil {
				logger.Println("cannot update table machine", err)
				fmt.Fprintf(w, "failed")
				return
			}
		}
		fmt.Fprintf(w, "Update complete. All agents will update themselves in the following one hour.")
	}
}

func getNextVersion() int {
	var v int
	rows, err := db.Query("select v from version order by v desc limit 1")
	if err != nil {
		logger.Println("cannot query table version")
		return 0
	}
	for rows.Next() {
		err = rows.Scan(&v)
		if err == nil {
			break
		}
	}
	v += 1
	return v
}

func unzipAndMove(filename string, version int) (string, error) {
	var files string
	os.RemoveAll("../up/" + strconv.Itoa(version) + "/")
	rd, err := zip.OpenReader(filename);
	if err != nil {
		return "", err
	}
	filenames := make([]string,0, 100)
	for _, f := range rd.File {
		fname := f.FileInfo().Name()
		// exclude dir
		if !strings.HasSuffix(fname, "/") {
			filenames = append(filenames, fname)
		}
		rc, err := f.Open()
		if err != nil {
			return "", err
		}
		os.MkdirAll("../up/" + strconv.Itoa(version) + "/" + path.Dir(fname), os.ModePerm)
		fw, _ := os.Create("../up/" + strconv.Itoa(version) + "/" + fname)
		if err != nil {
			return "", err
		}
		_, err = io.Copy(fw, rc)
		if err != nil {
			return "", err
		}
		if fw != nil {
			fw.Close()
		}
	}
	files = strings.Join(filenames, ";")
	//logger.Println(files)
	defer rd.Close()
	return files, nil
}

func main() {
	var err error
	os.Mkdir("../log/", 0666)
	os.Mkdir("../up/", 0666)
	logger = utils.InitLogger("../log/server.log")
	db, err = sql.Open("mysql", "root:@tcp(localhost:3306)/agentserver?charset=utf8")
	if err != nil {
		logger.Fatalln("cannot open database")
		os.Exit(1)
	}
	defer db.Close()

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("../static/"))))
	http.HandleFunc("/update", dealRequest)
	http.HandleFunc("/upload", uploadFile)
	http.HandleFunc("/store", store)
	err = http.ListenAndServe(":9090", nil)
	if err != nil {
		logger.Fatalln("ListenAndServe: ", err)
	}
}

