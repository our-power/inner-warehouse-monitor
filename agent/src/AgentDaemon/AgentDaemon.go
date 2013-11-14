package main 

import (
	"os"
	"time"
	"runtime"
	"net"
	"net/rpc"
	"os/exec"
)

type Daemon int
var PIDs map[string] *os.Process

/*
	Only those who start processes have the permission to kill them.
	As a work around, we listen to rpc request and kill the required process.
*/
func (Daemon) Kill (pName, reply *string) error {
	*reply = "success"
	err := PIDs[*pName].Kill()
	if err != nil {
		*reply = err.Error()
	}
	return nil
}

func accept() {
	sys := new(Daemon)
	server := rpc.NewServer()
	server.Register(sys)
	l, _ := net.Listen("tcp", "127.0.0.1:8773")
	for {
		server.Accept(l)
	}
}

func daemon(programName string) {
	for {
		prefix := "./"
		postfix := ""
		if runtime.GOOS == "windows" {
			postfix = ".exe"
		}
		toExec := prefix + programName + postfix
		cmd := exec.Command(toExec)
		err := cmd.Start()
		if err != nil {
			time.Sleep(time.Second * 5)
			continue
		}
		// Only by storing these *os.Process objects started just now could we kill them later. FindProcess is useless.
		PIDs[programName] = cmd.Process
		err = cmd.Wait()
		time.Sleep(time.Second * 100)
	}
}

func main() {
	done := make(chan bool, 1)
	mainProgram := "EccReportAgent"
	updateProgram := "AgentUpdate"
	PIDs = make(map[string] *os.Process)
	go accept()
	go daemon(mainProgram)
	go daemon(updateProgram)
	<- done
}
