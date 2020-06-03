package main

import (
	"log"
	"os"
	"time"
)

//CreateDir  文件夹创建
func CreateDir(path string) bool {
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		log.Println(err)
	}
	return true
}

//IsExist  判断文件夹/文件是否存在  存在返回 true
func IsExist(f string) bool {
	_, err := os.Stat(f)
	return err == nil || os.IsExist(err)
}

func init() {
	times := time.Now().Format("20060102")
	dir := "./switcherlog"
	logdir := "./switcherlog/" + times + ".log"
	//创建文件夹
	if !IsExist(dir) {
		CreateDir(dir)
	}
	logFile, err := os.OpenFile(logdir, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Println(err)
	}

	if times == time.Now().Format("20060102") {
		log.SetOutput(logFile) // 将文件设置为log输出的文件
		log.SetPrefix("[switcher2.4]")
		log.SetFlags(log.LstdFlags | log.Lshortfile) // | log.LUTC)
	}

}
