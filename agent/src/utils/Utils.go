package utils

import (
	"os"
	"fmt"
	"log"
	"net"
	"time"
	"errors"
	"strings"
	"net/http"
	"io/ioutil"
	"encoding/json"
)

type Settings struct {
	InModules []map[string] string
	OutModules []map[string] string
	UpdateServer []map[string] string
	Hb int
	Update int
	Role string
}

/* Load settings */
func LoadSettings() (settings Settings, err error) {
	bytes, err := ioutil.ReadFile("../etc/settings.json")
	if err != nil {
		return
	}
	err = json.Unmarshal(bytes, &settings)
	if err != nil {
		return
	}
	err = checkSettings(settings)
	if err != nil {
		return
	}
	return
}

/*
	Check settings obtained from file, some items must exist.
	Otherwise throw an error.
*/
func checkSettings(settings Settings) (err error) {
	// input check
	if len(settings.InModules)<=0||len(settings.InModules)>15 {
		err = errors.New("Too many or too few Input modules")
		return
	}
	for _, mod := range settings.InModules {
		_, ok1 := mod["name"]
		_, ok2 := mod["topic"]
		_, ok3 := mod["interval"]
		_, ok4 := mod["type"]
		if !ok1 || !ok2 || !ok3 || !ok4 {
			err = errors.New("Some Input modules are not properly configured.")
			return
		}
	}
	// output check
	if len(settings.OutModules)<=0||len(settings.OutModules)>15 {
		err = errors.New("Too many or too few output modules")
		return
	}
	for _, mod := range settings.OutModules {
		_, ok1 := mod["url"]
		_, ok2 := mod["host"]
		if !ok1 || !ok2 {
			err = errors.New("Some output modules are not properly configured.")
			return
		}
	}
	// update server check
	if len(settings.UpdateServer)!=1 {
		err = errors.New("One and only one UpdateServer can be configured")
		return
	}
	_, ok1 := settings.UpdateServer[0]["url"]
	_, ok2 := settings.UpdateServer[0]["host"]
	if !ok1 || !ok2 {
		err = errors.New("UpdateServer is not properly configured.")
		return
	}
	return
}

/*
	Get specific IP address (exclude ip starts with "127.0.0" and "0.0.0.0") and the hostname.
	Return: IP and hostname
*/
func GetLocalInfo() (ip string, hostName string, macAddress string, err error) {
	addrs, err := net.InterfaceAddrs()
    for _, ad := range addrs {
    	if tmp := strings.Split(ad.String(),"/")[0]; !strings.HasPrefix(tmp, "127.0.0") && !strings.HasPrefix(tmp, "0.0.0") {
    		ip = tmp
    		break
    	}
    }
	hostName, _ = os.Hostname()
	ifs, err := net.Interfaces()
	if err != nil || len(ifs) <= 0 {
		macAddress = "no_nic"
	} else {
		macAddress = ifs[0].HardwareAddr.String()
	}
	return
}

/*
	Do request. Returns a slice of byte.
	If the hostHeader string for a module is "" then we use no hostHeader for it.
*/
func ReadRemote(method string, urlString string, content string, hostHeader string, client *http.Client) (b []byte, err error) {
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

func BuildClient() *http.Client {
	var myTransport http.RoundTripper = &http.Transport {
		// Timeout is set to 10 seconds
		ResponseHeaderTimeout: time.Second * 10,
	}
	client := &http.Client{ Transport: myTransport }
	return client
}

/* Initiate and return a logger by the filename passed in */
func InitLogger(filename string) (logger *log.Logger) {
	logfile,err := os.OpenFile(filename, os.O_CREATE | os.O_RDWR | os.O_APPEND, 0666) 
	if err!=nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
	logger = log.New(logfile,"\r\n",log.Ldate|log.Ltime|log.Lshortfile)
	return
}
