package main

import (
	"os"
	"fmt"
	"flag"
	"net/http"
	"time"
	"utils"
	"strings"
)

var settings utils.Settings
var netClient *http.Client

var (
	topic = flag.String("topic", "", "The topic you want to send a message.")
	message = flag.String("message", "", "Message body.")
)

func main() {
	netClient = utils.BuildClient()
	var err error
	settings, err = utils.LoadSettings()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	flag.Parse()
	if *topic == "" || *message == "" {
		fmt.Println()
		os.Exit(1)
	}
	outModules(prepareOutput(*topic, *message, 0))
}

func prepareOutput(topic string, output string, interval int) (request string, content string) {

	timeStamp := time.Now().Unix()
	timeIdx := int(timeStamp)
	dateNow := time.Now().Format("20060102")
	ip, hostName, macAddress, err := utils.GetLocalInfo()
	if err != nil {
		fmt.Println(err.Error())
	}
	h, m, s := time.Now().Clock()
	if interval != 0 {
		timeIdx = (3600*h + 60*m + s) / interval
	}
	line := ""
	lines := strings.Split(output, "\n")
	if len(lines) > 3 && topic != "accessibility" {
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
		fmt.Println(err.Error())
		return
	}
	fmt.Printf("%s\n", resp)
}

/*
	Output the final content to all the output modules configured in the setting file,
*/
func outModules(srcString string, content string) {
	for i, _ := range settings.OutModules {
		if strings.HasSuffix(srcString, "error") {
			fmt.Printf("%s\n", srcString)
			continue
		}
		urlString := settings.OutModules[i]["url"] + "?" + srcString
		hostHeader := settings.OutModules[i]["host"]
		runOutModule(urlString, content, hostHeader)
	}
}
