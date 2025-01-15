package main

import (
    "log"
    "os"
    "time"
)

type Logger struct{
    logger *log.Logger
    level string
}
var logger *Logger

func (l *Logger) initFolder(){
    if _, err := os.Stat("Sandbox/Logs"); os.IsNotExist(err){
        os.Mkdir("Sandbox/Logs", 0755)
    }
    if _, err := os.Stat("Database/Logs"); os.IsNotExist(err){
        os.Mkdir("Database/Logs", 0755)
    }
}

func (l *Logger) setLogLevel(level string){
    l.level = level
    l.setPrefix(level)
}

func (l *Logger) setPrefix(level string){
    if level == "info"{
        l.logger.SetPrefix("[INFO] ")
    }else if level == "error"{
        l.logger.SetPrefix("[ERROR] ")
    }else if level == "debug"{
        l.logger.SetPrefix("[DEBUG] ")
    }else if level == "warning"{
        l.logger.SetPrefix("[WARNING] ")
    }
}

func (l *Logger)initFile(logType string){
    date := time.Now().Format("2006-01-02")
    if logType == "Sandbox"{
        file, err := os.OpenFile("Sandbox/Logs/Sandbox"+date+".log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
        if err != nil{
            log.Fatal(err)
        }
        l.logger = log.New(file, "", log.Ldate|log.Ltime|log.Lshortfile)
    }else if logType == "Database"{
        file, err := os.OpenFile("Database/Logs/Database"+date+".log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
        if err != nil{
            log.Fatal(err)
        }
        l.logger = log.New(file, "", log.Ldate|log.Ltime|log.Lshortfile)
    }
}

func (l *Logger) Debug(logType string,v... interface{}){
    l.initFile(logType)
    if l.level == "debug"{
        l.setPrefix("debug")
        l.logger.Println(v...)
        l.setPrefix(l.level)
    }
}

func (l *Logger) Info(logType string,v... interface{}){
    l.initFile(logType)
    if l.level == "info" || l.level == "debug"{
        l.setPrefix("info")
        l.logger.Println(v...)
        l.setPrefix(l.level)
    }
}

func (l *Logger) Error(logType string,v... interface{}){
    l.initFile(logType)
    l.setPrefix("error")
    l.logger.Println(v...)
    l.setPrefix(l.level)
}

func (l *Logger) Warning(logType string,v... interface{}){
    l.initFile(logType)
    l.setPrefix("warning")
    l.logger.Println(v...)
    l.setPrefix(l.level)
}

func InitLogger(level string)(*Logger){
    logger = &Logger{
        level: "info",
        logger: new(log.Logger),
    }
    logger.initFolder()
    logger.setLogLevel(level)
    return logger
}

func (l *Logger)test(){
    logger.Debug("Sandbox","Debug message")
    logger.Info("Sandbox","Info message")
    logger.Error("Sandbox","Error message")
    logger.Warning("Sandbox","Warning message")
}