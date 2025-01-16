package main

import (
	"os"
	"strings"
	"os/exec"
	"fmt"
	"path/filepath"
	"bufio"
)

//读取并编译代码文件

func isGoFile(filePath string) bool {
	return strings.HasSuffix(filePath, ".go")
}

func ComplieFile(filePath string)(string,string){
	logger.Info("Sandbox","Complie File: "+ filePath)
	if _, err := os.Stat(filePath); err != nil {
		logger.Error("Sandbox",err)
	}
	if isGoFile(filePath){
		path,baseName := compileGo(filePath)
		return path,baseName
	}else{
		err := fmt.Errorf("Undefined File Type")
		logger.Error("Sandbox",err)
		return "",""
	}
}

func getPathAndBaseName(filePath string)(string,string){
	dir := filepath.Dir(filePath)
	baseName := strings.TrimSuffix(filepath.Base(filePath), filepath.Ext(filePath))
	logger.Debug("Sandbox",fmt.Sprintf("Dir: %s, BaseName: %s",dir,baseName))
	return dir,baseName
}

func compileGo(filePath string)(string,string){
	path,baseName := getPathAndBaseName(filePath)
	command := "go" + " build -o " + path + "/"+ baseName + " " + filePath
	_,err,outputErr := RunCommand(command)
	logger.Info("Sandbox","Execute Command: "+command)
	if err != nil {
		logger.Error("Sandbox",fmt.Errorf("Compile Error: %s",outputErr))
		logger.Error("Sandbox",err)
		return "",""
	}
	logger.Info("Sandbox","Compile Success")
	return path,baseName
}

//TODO: 多线程计算CPU时长和内存占用；超时处理
func RunCommand(command string)(string,error,string){
	cmd := exec.Command("sh", "-c", command)
	var output string
	var outputErr string
	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		logger.Error("Sandbox","Error creating StdoutPipe: ", err)
		return "",err,""
	}
	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		logger.Error("Sandbox","Error creating StderrPipe: ", err)
		return "",err,""
	}
	if err := cmd.Start(); err != nil {
		logger.Error("Sandbox","Error starting command: ", err)
		return "",err,""
	}

	go func() {
		scanner := bufio.NewScanner(stdoutPipe)
		for scanner.Scan() {
			line := scanner.Text()
			output += line + "\n"
		}
	}()
	go func() {
		scanner := bufio.NewScanner(stderrPipe)
		for scanner.Scan() {
			line := scanner.Text()
			outputErr += line + "  "
		}
	}()

	cmd.Wait()
	return output,nil,outputErr
}