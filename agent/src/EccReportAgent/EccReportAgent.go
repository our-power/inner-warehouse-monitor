package main 

import (
	"os"
	"fmt"
	"log"
	"time"
	"runtime"
	"strconv"
	"strings"
	"os/exec"
//	"net/url"
	"net/http"
	"os/signal"
	"utils"
)

const APPVERSION = "1.1"

var settings utils.Settings
var logger *log.Logger
var stop bool
var netClient *http.Client

/* Send heartbeat signal to server */
func heartbeat() {
	outModules(prepareOutput("heartbeat", "alive", settings.Hb))
	c := time.Tick(time.Duration(settings.Hb) * time.Second)
	for _ = range c {
		if stop { break }
		outModules(prepareOutput("heartbeat", "alive", settings.Hb))
	}
}

func register() {
	appVersion := runtime.GOOS + APPVERSION
	outModules(prepareOutput("register", appVersion + settings.Role, 0))
}

func runInModule(name string, topic string, interval string) {
	itv, _ := strconv.Atoi(interval)
	c := time.Tick(time.Duration(itv) * time.Second)
	for _ = range c {
		if stop { break }
		var result string
		split := strings.Split(name, "@")
		cmd := exec.Command(split[0], split[1:]...)
		buf, err := cmd.Output()
		if err != nil {
			result = err.Error()
		}
		result = string(buf)
		request, content := prepareOutput(topic, result, itv)
		// fmt.Println(request, content)
		outModules(request, content)
	}
}

/*
	Invoke scripts according to the OS type
	Return: A slice of strings, which includes contents returned by the scripts.
*/
func inModules() {
	for _, mod := range settings.InModules {
		var prepend string
		var ext string
		modPath := "../mod/"
		modName := mod["name"]
		fmt.Println(mod["type"])
		if mod["type"] == "perf" {
			prepend = "C:\\Windows\\System32\\typeperf.exe@-sc@2@-cf@"
			ext = ".txt"
		} else if mod["type"] == "exec" {
			prepend = "C:\\Windows\\System32\\cscript.exe@/nologo@"
			ext = ".vbs"
		}
		modPathName := modPath + modName + ext
		_, err := os.Stat(modPathName)
		if nil != err && !os.IsExist(err) {
			logger.Println(modPathName, "doesn't exist. Ignored.")
			continue
		}
		command := prepend + modPathName
		topic := mod["topic"]
		interval := mod["interval"]
		go runInModule(command, topic, interval)
	}
}

/*
	Format the output string to a standard one
	Param: topic, of the monitoring content; output, original content string.
	Return: Formatted string
*/
func prepareOutput(topic string, output string, interval int) (request string, content string) {
	
	timeStamp := time.Now().Unix()
	timeIdx := int(timeStamp)
	dateNow := time.Now().Format("20060102")
	ip, hostName, macAddress, err := utils.GetLocalInfo()
	if err != nil {
		logger.Println(err.Error())
	}
	h, m, s := time.Now().Clock()
	if interval != 0 {
		timeIdx = (3600*h + 60*m + s) / interval
	}
	line := ""
	lines := strings.Split(output, "\n")
	if len(lines) > 3 {
		line = lines[3]
		line = strings.Replace(line, "\"", "", -1)
	} else {
		line = output
	}
	line = strings.Trim(line, "\r")
	content = fmt.Sprintf("%s\r\n%d\r\n%s\r\n%s\r\n%s\r\n%s", dateNow, timeIdx, ip, hostName, macAddress, line)
	request = fmt.Sprintf("topic=%s", topic)
	return
}

/*
	Call ReadRemote and print for debug purpose
*/
func runOutModule(urlString string, content string, hostHeader string) {
	resp, err := utils.ReadRemote("POST", urlString, content, hostHeader, netClient)
	if err != nil {
	    logger.Println(err.Error())
	    return
	}
	fmt.Printf("%s\n", resp)
}

/*
	Output the final content to all the output modules configured in the setting file,
*/
func outModules(srcString string, content string) {
	if stop {
		return
	}
	for i, _ := range settings.OutModules {
		if strings.HasSuffix(srcString, "error") {
			logger.Printf("%s\n", srcString)
			continue
		}
		urlString := settings.OutModules[i]["url"] + "?" + srcString
		hostHeader := settings.OutModules[i]["host"]
		go runOutModule(urlString, content, hostHeader)
	}
}

func main() {
	done := make(chan bool, 1)
	logger = utils.InitLogger("../log/agent.log")
	// code snippet: capture Ctrl-C signal and handle it
	cc := make(chan os.Signal, 1)
	signal.Notify(cc, os.Interrupt, os.Kill)
	go func(){
	    for _ = range cc {
	    	stop = true
	    	time.Sleep(time.Second)
	        logger.Println("Agent stopped")
	        done <- true
	        os.Exit(0)
	    }
	}()
	logger.Println("Agent started")
	netClient = utils.BuildClient()
	var err error
	settings, err = utils.LoadSettings()
	if err != nil {
		logger.Fatalln(err.Error())
		os.Exit(1)
	}
	logger.Printf("Settings: %v", settings)
	register()
	time.Sleep(time.Second)
	go heartbeat()
	go inModules()
	<- done
}
