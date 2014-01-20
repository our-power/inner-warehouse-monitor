package controllers

import (
	"strings"
	"big_brother/models"
	"time"
	"crypto/md5"
	"io"
	"fmt"
)

func HasTheRight(permission string, permissions interface{}) bool {
	hasTheRight := false
	permission_list := strings.Split(permissions.(string), "|")
	for _, value := range permission_list {
		if permission == value {
			hasTheRight = true
			break
		}
	}
	return hasTheRight
}


func StoreTrace(operation string, userName interface{}) error {
	o.Using("admin")
	/*
	type Trace struct {
		Id int
		User string
		Do_what string
		That_time string
	}
	*/
	trace := models.Trace{
		User: userName.(string),
		Do_what: operation,
		That_time: time.Now().Format("2006-01-02 15:04:05"),
	}
	//id, err := o.Insert(&trace)
	_, err := o.Insert(&trace)
	return err
}

func GenMd5Passwd(passwd string) string {
	h := md5.New()
	io.WriteString(h, passwd)
	return fmt.Sprintf("%x", h.Sum(nil))
}
