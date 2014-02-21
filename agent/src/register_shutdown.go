package main

import (
	"os"
	"fmt"
	"flag"
	"net"
	"net/http"
	"time"
	"strings"
	"io/ioutil"
	"log"
)

var netClient *http.Client
var logger *log.Logger

var (
	reg = flag.Bool("register", false, "register or shutdown?")
)

func main() {
	netClient = buildClient()
	flag.Parse()
	b, err := ioutil.ReadFile("C:/salt/var/register_shutdown.conf")
	if err != nil {
		logger.Println(err.Error())
		os.Exit(1)
	}
	logger = initLogger("C:/salt/var/register_shutdown.log")
	lines := strings.Split(string(b), ";")
	if len(lines) != 2 {
		logger.Println("too many or too few parameters")
		os.Exit(1)
	}
	url := lines[0] + "register"
	role := lines[1]
	var content string
	if *reg {
		content = "windows1.1," + role
	} else {
		content = "shutdown"
	}
	output := prepareOutput(content)
	resp, err := readRemote("POST", url, output, "", netClient)
	if err != nil {
		logger.Println(err.Error())
		os.Exit(1)
	}
	logger.Printf("%s\n", resp)
	logger.Println(url + "\r\n" + output)
}

func prepareOutput(output string) (content string) {

	timeStamp := time.Now().Unix()
	timeIdx := int(timeStamp)
	dateNow := time.Now().Format("20060102")
	ip, hostName, macAddress, err := getLocalInfo()
	if err != nil {
		logger.Println(err.Error())
	}
	line := strings.Trim(output, "\r")
	content = fmt.Sprintf("%s\r\n%d\r\n%s\r\n%s\r\n%s\r\n%s", dateNow, timeIdx, ip, hostName, macAddress, line)
	return
}

func buildClient() *http.Client {
	var myTransport http.RoundTripper = &http.Transport {
		// Timeout is set to 10 seconds
		ResponseHeaderTimeout: time.Second * 10,
	}
	client := &http.Client{ Transport: myTransport }
	return client
}

func getLocalInfo() (ip string, hostName string, macAddress string, err error) {
	addrs, err := net.InterfaceAddrs()
	var index int
    for i, ad := range addrs {
    	if tmp := strings.Split(ad.String(),"/")[0]; !strings.HasPrefix(tmp, "127.0.0") && !strings.HasPrefix(tmp, "0.0.0") {
    		ip = tmp
    		index = i
    		break
    	}
    }
	hostName, _ = os.Hostname()
	ifs, err := net.Interfaces()
	if err != nil || len(ifs) <= index {
		macAddress = "no_nic"
	} else {
		macAddress = ifs[index].HardwareAddr.String()
	}
	return
}

func readRemote(method string, urlString string, content string, hostHeader string, client *http.Client) (b []byte, err error) {
	req, _ := http.NewRequest(method, urlString, strings.NewReader(content))
	if hostHeader != "" {
		req.Header.Set("Host", hostHeader)
	} 
	res, err := client.Do(req)
	if err != nil {
	    return
	}
	resp, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
	    return
	}
	b = resp
	return
}

func initLogger(filename string) (logger *log.Logger) {
	logfile, err := os.OpenFile(filename, os.O_CREATE | os.O_RDWR | os.O_APPEND, 0666)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
	logger = log.New(logfile, "\r\n", log.Ldate | log.Ltime | log.Lshortfile)
	return
}