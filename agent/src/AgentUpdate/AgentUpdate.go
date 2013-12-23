package main

import (
	"os"
	"log"
	"time"
	"utils"
	"strconv"
	"strings"
	"net"
	"net/url"
	"net/rpc"
	"net/http"
	"os/signal"
)

var settings utils.Settings
var logger *log.Logger
var stop bool
var client *rpc.Client
var netClient *http.Client

/*
	Check the update server to get update information.
	If the content is not empty, then split the file names split by a semicolon.
	File names are like these: mod/cpu_usage.vbs or bin/EccReportAgent
*/
func checkList() []string {
	var list []string
	ip, _, _, _ := utils.GetLocalInfo()
	urlString := settings.UpdateServer[0]["url"] + "?action=get_list&ip=" + ip
	hostString := settings.UpdateServer[0]["host"]
	resp, err := utils.ReadRemote("GET", urlString, "", hostString, netClient)
	if err != nil {
		logger.Println(err.Error())
		return list
	}
	if string(resp) == "" {
		return list
	}
	list = strings.Split(string(resp), ";")
	return list
}

/*
	If the file name prepend with "../" exists, rename it by appending an
	timestamp for now. And then download the file to the corresponding place.
*/
func downloadAndReplaceFile(filename string, version string) bool {
	theFile := "../" + filename
	if _, err := os.Stat(theFile); err == nil {
		// If local file exists.
		err = os.Rename(theFile, theFile + "." + strconv.FormatInt(time.Now().Unix(), 10))
		if err != nil {
			logger.Println(err.Error())
			return false
		}
	}
	urlString := settings.UpdateServer[0]["url"] + "?action=get_file&v=" + version + "&name=" + url.QueryEscape(filename)
	hostHeader := settings.UpdateServer[0]["host"]
	resp, err := utils.ReadRemote("GET", urlString, "", hostHeader, netClient)
	if err != nil {
		logger.Println(err.Error())
		return false
	}
	if len(resp) != 0 {
		// empty response means no such file exists, we should do nothing.
		f, err := os.OpenFile(theFile, os.O_CREATE | os.O_WRONLY | os.O_TRUNC, 0666)
		defer f.Close()
		if err != nil {
			logger.Println(err.Error())
			return false
		}
		f.Write(resp)
		return true
	}
	return false
}

func setDoneFlag() {
	ip, _, _, _ := utils.GetLocalInfo()
	urlString := settings.UpdateServer[0]["url"] + "?action=set_done&ip=" + ip
	hostHeader := settings.UpdateServer[0]["host"]
	_, err := utils.ReadRemote("GET", urlString, "", hostHeader, netClient)
	if err != nil {
		logger.Println(err.Error())
	}
}

/*
	Stop EccReportAgent by RPC the daemon if updates exist.
	Then download these files.
*/
func stopAndUpdate() {
	c := time.Tick(time.Duration(settings.Update)*time.Second)
	for _ = range c {
		if stop {
			return
		}
		if files := checkList(); len(files) > 1 {
			allDone := true
			// len(files) > 1 means "version + file list" or none
			// make a rpc call to daemon, let daemon kill the agent process
			var reply string
			err := client.Call("Daemon.Kill", "EccReportAgent", &reply)
			if err != nil {
				logger.Println("RPC error:", err.Error())
				return
			}
			if reply != "success" {
				logger.Println(reply)
				return
			}
			version := files[0]
			for _, filename := range files[1:] {
				succ := downloadAndReplaceFile(filename, version)
				allDone = allDone && succ
				if !succ {
					logger.Println("update failed for:", filename)
				}
				time.Sleep(time.Second*3)
			}
			if allDone {
				logger.Println("Complete updating:", files)
				setDoneFlag()
			}
			// reload settings in case that settings updated
			settings, err = utils.LoadSettings()
			if err != nil {
				logger.Fatalln(err)
				os.Exit(1)
			}
		}
	}
}

func main() {
	done := make(chan bool, 1)
	logger = utils.InitLogger("../log/update.log")
	cc := make(chan os.Signal, 1)
	signal.Notify(cc, os.Interrupt, os.Kill)
	go func() {
		for _ = range cc {
			stop = true
			time.Sleep(time.Second)
			logger.Println("Updater stopped")
			done <- true
			os.Exit(0)
		}
	}()
	// dial to the daemon, and create a rpc client
	conn, err := net.Dial("tcp", "127.0.0.1:8773")
	if err != nil {
		logger.Println(err.Error())
	}
	client = rpc.NewClient(conn)
	netClient = utils.BuildClient()
	logger.Println("Updater started")
	settings, err = utils.LoadSettings()
	if err != nil {
		logger.Fatalln(err)
		os.Exit(1)
	}
	stopAndUpdate()
	<-done
}
