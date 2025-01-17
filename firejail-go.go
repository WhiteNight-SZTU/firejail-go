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
	output,err,outputErr := RunCommand(command)
	logger.Debug("Sandbox","Output: ", output)
	if err != nil || strings.Contains(outputErr,"error") {
		if strings.Contains(outputErr,"error") {
			err = fmt.Errorf("Error executing firejail: %s",outputErr)
		}
		logger.Error("Sandbox","Error executing firejail: ",err)
		return nil
	}
	return nil
}



func main(){
	InitLogger("debug")
	path,filename := ComplieFile("Sandbox/test.go")
	if filename == "" {
		err := fmt.Errorf("Cant't find compiled file: %s in %s",filename,path)
		logger.Error("Sandbox",err)
		return 
	}
	firejail := Firejail{cpu: "20", mem: "700m", profile: "noprofile"}
	err := firejail.runFirejail(path,filename)
	if err != nil {
		logger.Error("Sandbox","Execute Firejail error: ",err)
		return
	}
}
