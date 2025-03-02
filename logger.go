package main

import (
    "log"
    "os"
    "time"
    "runtime"
    "fmt"
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
    _, file, line, ok := runtime.Caller(2)
    callerInfo := ""
    if ok {
        callerInfo = fmt.Sprintf("%s:%d ", file, line)
    }
    if level == "info"{
        l.logger.SetPrefix(fmt.Sprintf("[INFO] %s", callerInfo))
    }else if level == "error"{
        l.logger.SetPrefix(fmt.Sprintf("[ERROR] %s", callerInfo))
    }else if level == "debug"{
        l.logger.SetPrefix(fmt.Sprintf("[DEBUG] %s", callerInfo))
    }else if level == "warning"{
        l.logger.SetPrefix(fmt.Sprintf("[WARNING] %s", callerInfo))
    }else if level == "userOutput"{
        //std ouput
        l.logger.SetPrefix(fmt.Sprintf("[USEROUTPUT] %s", callerInfo))
    }

}

func (l *Logger)initFile(logType string){
    date := time.Now().Format("2006-01-02")
    if logType == "Sandbox"{
        file, err := os.OpenFile("Sandbox/Logs/"+date+".log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
        if err != nil{
            log.Fatal(err)
        }
        l.logger = log.New(file, "", log.Ldate|log.Ltime)
    }else if logType == "Database"{
        file, err := os.OpenFile("Database/Logs/"+date+".log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
        if err != nil{
            log.Fatal(err)
        }
        l.logger = log.New(file, "", log.Ldate|log.Ltime)
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

func (l *Logger) UserOutput(logType string,v... interface{}){
    l.initFile(logType)
    l.setPrefix("userOutput")
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