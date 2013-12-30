package util

import (
    "log"
    "os"
)

func HandleException(logFileName string, msg string) {
    if err := recover(); err != nil {
        fh, _ := os.OpenFile(logFileName, os.O_RDWR | os.O_APPEND | os.O_CREATE, 0666)
        defer fh.Close()

        logger := log.New(fh, "\r\n", log.LstdFlags)
        logger.Println("****************************************************")
        //这里的err其实就是panic传入的内容
        logger.Println(err)
        logger.Println(msg)
        logger.Println("#####################################################")
    }

}