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
        h.Logger.Println(err)
        h.Logger.Println(msg)
        h.Logger.Println("#####################################################")
    }
}