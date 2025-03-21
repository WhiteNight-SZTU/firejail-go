package main

import (
	"os"
	"strings"
	"os/exec"
	"fmt"
	"path/filepath"
	"bufio"
	"github.com/creack/pty" 
	"bytes"
)

//读取并编译代码文件

func isGoFile(filePath string) bool {
	return strings.HasSuffix(filePath, ".go")
}

func isPythonFile(filePath string) bool {
	return strings.HasSuffix(filePath, ".py")
}

func isCPPFile(filePath string) bool {
	return strings.HasSuffix(filePath, ".cpp")
}

func ComplieFile(filePath string)(string,string){
	logger.Info("Sandbox","Complie File: "+ filePath)
	if _, err := os.Stat(filePath); err != nil {
		logger.Error("Sandbox",err)
	}
	if isGoFile(filePath){ //Golang
		path,baseName := compileGo(filePath)
		return path,baseName
	}else if isPythonFile(filePath){ //Python
		path,baseName := compliePython(filePath)
		return path,baseName+".py"
	}else if isCPPFile(filePath){ //C++
		path,baseName := compileCPP(filePath)
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
	logger.Info("Sandbox","Execute Command: "+command)
	_,err,outputErr := RunCommandInFirejail(command)
	if outputErr != "" {
		logger.Warning("Sandbox",fmt.Errorf("Compile Error: %s",outputErr))
	}else if err != nil {
		logger.Error("Sandbox",fmt.Errorf("Compile Error: %s",outputErr))
		logger.Error("Sandbox",err)
		return "",""
	}

	logger.Info("Sandbox","Compile Success")
	return path,baseName
}

func compileCPP(filePath string)(string,string){
	path,baseName := getPathAndBaseName(filePath)
	command := "g++ " + filePath +" -o " + path + "/" + baseName
	_,err,outputErr := CompileInFirejail(command)
	if err != nil || strings.Contains(outputErr,"error") { 
		logger.Error("Sandbox",fmt.Errorf("Runtime Error: %s",outputErr))
		return "",""
	}
	return path,baseName
}

func compliePython(filePath string)(string,string){
	path,baseName := getPathAndBaseName(filePath)
	return path,baseName
}



/* 读取输入，执行程序并返回输出
 * @param command: 执行命令（命令可包含firejail指令本体）
 * filename: 输入文件的文件名（无后缀）
*/
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
			if strings.Contains(line, "The new log directory is") {
				outputErr += line + "\n"
			} else {
				output += line + "\n"
			}
		}
	}()
	go func() {
		scanner := bufio.NewScanner(stderrPipe)
		for scanner.Scan() {
			line := scanner.Text()
			outputErr += line + "  " + "\n"
		}
	}()

	cmd.Wait()
	logger.Debug("Sanbox","Command Execute Success")
	return output,nil,outputErr
}

//Run In Pty
func RunCommandInPTY(command string)(string,error,string){
	cmd := exec.Command("sh", "-c", command)
	pty, err := pty.Start(cmd)
	if err != nil {
		logger.Error("Sandbox","Error starting command: ", err)
		return "",err,""
	}
	defer pty.Close()
	var output string
	var outputErr string
	scanner := bufio.NewScanner(pty)
	for scanner.Scan() {
		line := scanner.Text()
		output += line + "\n"
	}
	cmd.Wait()
	return output,nil,outputErr
}

func RunCommandInFirejail(command string)(string,error,string){
	cmd := exec.Command("sh", "-c", "firejail " + command)
	logger.Debug("Sandbox","Execute Command: firejail " + command)
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

/* 读取输入，执行程序并返回输出
 * @param command: 执行命令（命令可包含firejail指令本体）
 * @param filename: 输入文件的文件名（无后缀）
 * 
 * filename: 输入文件的文件名（无后缀）
*/
func RunCommandWithInput(command string,filename string)(string,error,string){
	cmd := exec.Command("sh", "-c", command)
	var output string
	var outputErr string

	inputBytes, err := os.ReadFile("Sandbox/Input/" + filename + ".in")
	if err != nil {
		logger.Error("Sandbox",err)
	}
	var buf bytes.Buffer
	buf.Write(inputBytes)
	cmd.Stdin = &buf

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
		var lastLine string
		for scanner.Scan() {
			line := scanner.Text()
			logger.Debug("Sandbox", line)
			if strings.Contains(line, "The new log directory is") {
				outputErr += line + "\n"
			} else {
				if lastLine != "" {
					output += lastLine + "\n"
				}
				lastLine = line
			}
		}
		if lastLine != "" {
			output += lastLine
		}
	}()
	go func() {
		scanner := bufio.NewScanner(stderrPipe)
		for scanner.Scan() {
			line := scanner.Text()
			outputErr += line + "  " + "\n"
		}
	}()

	cmd.Wait()
	return output,nil,outputErr
}


/* 在Firejail中编译程序
 * @param command: 编译命令(不包含firejail指令本体
 * 指定 /etc/firejail/wn_compile.profile 作为配置文件
*/
func CompileInFirejail(command string)(string,error,string){
	command = "firejail --profile=/etc/firejail/wn_compile.profile " + command
	cmd := exec.Command("sh", "-c", command)
	logger.Debug("Sandbox","Execute Command:" + command)
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