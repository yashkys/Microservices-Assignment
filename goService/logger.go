package main

import (
	"log"
	"os"
)

var (
	infoLog    *log.Logger
	successLog *log.Logger
	warningLog *log.Logger
	errorLog   *log.Logger
)

func initLogger() {
	infoLog = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	successLog = log.New(os.Stdout, "\033[32mSUCCESS\033[0m: ", log.Ldate|log.Ltime|log.Lshortfile)
	warningLog = log.New(os.Stdout, "\033[33mWARNING\033[0m: ", log.Ldate|log.Ltime|log.Lshortfile)
	errorLog = log.New(os.Stdout, "\033[31mError\033[0m: ", log.Ldate|log.Ltime|log.Lshortfile)
}
