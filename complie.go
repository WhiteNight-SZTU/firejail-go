package main

import (
	"os"
	"strings"
	"os/exec"
	"fmt"
	"path/filepath"
)

//读取并编译代码文件

func isGoFile(filePath string) bool {
	return strings.HasSuffix(filePath, ".go")
}

func ComplieFile(filePath string){
	logger.Info("Sandbox","Complie File: "+filePath)
	if _, err := os.Stat(filePath); err != nil {
		logger.Error("Sandbox",err)
	}
	if isGoFile(filePath){
		compileGo(filePath)
	}
}

func getPathAndBaseName(filePath string)(string,string){
	dir := filepath.Dir(filePath)
	baseName := strings.TrimSuffix(filepath.Base(filePath), filepath.Ext(filePath))
	logger.Debug("Sandbox",fmt.Sprintf("Dir: %s, BaseName: %s",dir,baseName))
	return dir,baseName
}

func compileGo(filePath string){
	path,baseName := getPathAndBaseName(filePath)
	command := "go" + " build -o " + path + "/"+ baseName + " " + filePath
	cmd := exec.Command("sh", "-c", command)
	logger.Info("Sandbox","Execute Command: "+command)
	if err := cmd.Run(); err != nil {
		logger.Error("Sandbox",err)
	}
}
