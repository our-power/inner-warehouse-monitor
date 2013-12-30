package util

import (
    "log"
    "os"
)

type ExceptionHandler struct {
    Logger *log.Logger
}

func (h *ExceptionHandler)HandleException(msg string){
    if err := recover(); err != nil {
        h.Logger.Println("****************************************************")
        //这里的err其实就是panic传入的内容
        h.Logger.Fatalln(err)
        h.Logger.Println(msg)
        h.Logger.Println("#####################################################")
    }
}

func InitHandler(logFileName string, prefix string)(exceptionHandler *ExceptionHandler) {
    fh, _ := os.OpenFile(logFileName, os.O_RDWR | os.O_APPEND | os.O_CREATE, 0777)
    defer fh.Close()

    l := log.New(fh, prefix, log.LstdFlags)

    exceptionHandler = &ExceptionHandler {
        Logger: l,
    }
    return
}