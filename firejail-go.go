package main

import (
	"fmt"
	"strings"
)

type Firejail struct {
	cpu string
	mem string
	profile string
}

func (f Firejail)runFirejail(path string, filename string)error {
	command := "firejail --" + f.profile + " --rlimit-cpu=" + f.cpu + " --rlimit-as=" + f.mem
	command += " " + path + "/" + filename
	logger.Debug("Sandbox","Executing command: ", command)
	print_flag := false


	output,err,outputErr := RunCommand(command) //compile.go 

	if strings.Contains(output, "Executing command") {
		print_flag = true
	}
	if strings.Contains(output, "Answer") {
		logger.UserOutput("Sandbox","Output: ", output)
	}else if print_flag {
		logger.Debug("Sandbox","Output: ", output)
	}
	if err != nil || strings.Contains(outputErr,"error") {
		if strings.Contains(outputErr,"error") {
			logger.Error("Sandbox","Error executing firejail: ",outputErr)
			return err
		}
		logger.Error("Sandbox","Error executing firejail: ",err)
		return nil
	}
	return nil
}

func test_cpu(){
	path,filename := ComplieFile("Sandbox/test_cpu.go")
	if filename == "" {
		err := fmt.Errorf("Cant't find compiled file: %s in %s",filename,path)
		logger.Error("Sandbox",err)
		return 
	}
	firejail := Firejail{cpu: "10", mem: "700m", profile: "noprofile"}
	err := firejail.runFirejail(path,filename)
	if err != nil {
		logger.Error("Sandbox","Execute Firejail error: ",err)
		return
	}
}

func test_memory(){
	path,filename := ComplieFile("Sandbox/test_memory.go")
	if filename == "" {
		err := fmt.Errorf("Cant't find compiled file: %s in %s",filename,path)
		logger.Error("Sandbox",err)
		return 
	}
	firejail := Firejail{cpu: "512", mem: "1g", profile: "noprofile"}
	err := firejail.runFirejail(path,filename)
	if err != nil {
		logger.Error("Sandbox","Execute Firejail error: ",err)
		return
	}
}
func main() {
    InitLogger("debug")
    test_cpu()
    test_memory()
}
