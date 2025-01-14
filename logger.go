package main

import (
    "log"
    "os"
    "time"
	"fmt"
)

// 日志管理
type Logger struct {
	logger *log.Logger
	level string 
	/* 日志级别 
	debug info warn error
	*/
}

func InitFolder() error{
	if _, err := os.Stat("Database"); os.IsNotExist(err) {
		err := os.Mkdir("Database", os.ModePerm)
		if err != nil {
			fmt.Println("创建日志文件夹失败:", err)
			return err
		}
	}
	if _, err := os.Stat("Sandbox"); os.IsNotExist(err) {
		err := os.Mkdir("Sandbox", os.ModePerm)
		if err != nil {
			fmt.Println("创建日志文件夹失败:", err)
		}
	}
	return nil
}

func InitLogger(logType string)*Logger{
	InitFolder()
	currentDate := time.Now().Format("2006-01-02")
	logPath := logType +"/" + "Logs/" + currentDate + ".log"
	logFile, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("创建日志失败:", err)
	}
	logger := log.New(logFile, "", log.LstdFlags|log.Lshortfile)
	l := &Logger{logger: logger}
	l.SetLoggerLevel("info")
	return l
}

func (l *Logger)SetLoggerLevel(level string) {
	switch level {
	case "debug":
		l.logger.SetPrefix("[DEBUG]")
	case "info":
		l.logger.SetPrefix("[INFO]")
	case "warn":
		l.logger.SetPrefix("[WARN]")
	case "error":
		l.logger.SetPrefix("[ERROR]")
	default:
		l.logger.SetPrefix("[INFO]")
	}
	l.level = level
}

func (l *Logger)Debug(v ...interface{}) {
	if l.level == "debug" {
		l.logger.SetPrefix("[DEBUG]")
		l.logger.Println(v ...)
		l.SetLoggerLevel(l.level)
	}
}
func (l *Logger)Info(v ...interface{}) {
	l.logger.SetPrefix("[INFO]")
	l.logger.Println(v ...)
	l.SetLoggerLevel(l.level)
}
func (l *Logger)Warn(v ...interface{}) {
	l.logger.SetPrefix("[WARN]")
	l.logger.Println(v ...)
	l.SetLoggerLevel(l.level)
}
func (l *Logger)Error(v ...interface{}) {
	l.logger.SetPrefix("[ERROR]")
	l.logger.Println(v ...)
	l.SetLoggerLevel(l.level)
}

func main(){
	logger := InitLogger("Database")
	logger.SetLoggerLevel("debug")
	logger.Debug("This is a debug message.")
	logger.Info("This is an info message.")
	logger.Warn("This is a warn message.")
	logger.Error("This is an error message.")
	logger.Debug("This is a debug message.")
	logger.Info("This is an info message.")
	logger.Warn("This is a warn message.")
	logger.Error("This is an error message.")
}
