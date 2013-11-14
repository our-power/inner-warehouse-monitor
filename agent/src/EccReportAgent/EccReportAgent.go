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
	"net/url"
	"net/http"
	"os/signal"
	"utils"
)

const APPVERSION = "1.0"

var settings utils.Settings
var logger *log.Logger
var stop bool
var netClient *http.Client

/* Send heartbeat signal to server */
func heartbeat() {
	outModules(prepareOutput("1007", "alive", 0))
	c := time.Tick(time.Duration(settings.Hb) * time.Second)
	for _ = range c {
		if stop { break }
		outModules(prepareOutput("1007", "alive", 0))
	}
}

func runInModule(name string, bid string, interval string) {
	itv, _ := strconv.Atoi(interval)
	c := time.Tick(time.Duration(itv) * time.Second)
	for _ = range c {
		if stop { break }
		var result string
		split := strings.Split(name, " ")
		cmd := exec.Command(split[0], split[1:]...)
		buf, err := cmd.Output()
		if err != nil {
			result = err.Error()
		}
		result = string(buf)
		output := prepareOutput(bid, result, itv)
		outModules(output)
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
		switch runtime.GOOS {
			case "windows":
				if mod["windows"] != "1" {
					continue 
				}
				prepend = "c:/windows/system32/cscript.exe /nologo "
				ext = ".vbs"
			case "linux":
				if mod["linux"] != "1" { 
					continue 
				}
				prepend = "/bin/bash "
				ext = ".sh"
		}
		modPathName := modPath + modName + ext
		_, err := os.Stat(modPathName)
		if nil != err && !os.IsExist(err) {
			logger.Println(modPathName, "doesn't exist. Ignored.")
			continue
		}
		command := prepend + modPathName
		bid := mod["bid"]
		interval := mod["interval"]
		go runInModule(command, bid, interval)
	}
}

/*
	Format the output string to a standard one
	Param: bid, of the monitoring content; output, original content string.
	Return: Formatted string
*/
func prepareOutput(bid string, output string, interval int) string {
	
	timeStamp := time.Now().Unix()
	timeIdx := int(timeStamp)
	dateNow := time.Now().Format("20060102")
	ip, hostName, err := utils.GetLocalInfo()
	if err != nil {
		logger.Println(err.Error())
	}
	h, m, s := time.Now().Clock()
	if interval != 0 {
		timeIdx = (3600*h + 60*m + s) / interval
	}
	content := fmt.Sprintf("%s\t%d\t%s\t%s\t%s", dateNow, timeIdx, ip, hostName, output)
	appVersion := runtime.GOOS + APPVERSION
	result := fmt.Sprintf("app=%s&bid=%s&time=%d&content=%s", appVersion, bid, timeStamp, url.QueryEscape(content))
	return result
}

/*
	Call ReadRemote and print for debug purpose
*/
func runOutModule(urlString string, hostHeader string) {
	resp, err := utils.ReadRemote(urlString, hostHeader, netClient)
	if err != nil {
	    logger.Println(err.Error())
	    return
	}
	fmt.Printf("%s\n", urlString)
	fmt.Printf("%s\n", resp)
}

/*
	Output the final content to all the output modules configured in the setting file,
*/
func outModules(srcString string) {
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
		go runOutModule(urlString, hostHeader)
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
	go heartbeat()
	go inModules()
	<- done
}
