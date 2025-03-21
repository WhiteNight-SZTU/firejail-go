package main

import (
	"fmt"
	"strings"
	"io/ioutil"
)

// 0 AC Accepted
// 1 WA Wrong Answer
// 2 RE Runtime Error
// 3 TLE Time Limit Exceeded

func CompareAnswer(output string,filename string)(int,string,error){
	//将output写入 filename.output
	output_path := fmt.Sprintf("Sandbox/Output/%s.output",filename)
	err := ioutil.WriteFile(output_path,[]byte(output),0644)
	if err != nil {
		logger.Error("Sandbox","Error writing output file"+err.Error())
		return 2,"RE",err
	}
	answer_path := fmt.Sprintf("Sandbox/Answer/%s.answer",filename)
	answer,err := ioutil.ReadFile(answer_path)
	if err != nil {
		logger.Error("Sandbox","Error reading answer file"+err.Error())
		return 2,"RE",err
	}
	logger.Debug("Sandbox","Output:"+output)
	logger.Debug("Sandbox","Answer:"+string(answer))
	if strings.Compare(output,string(answer)) == 0 {
		return 0,"AC",nil
	}else{
		return 1,"WA",nil
	}
	return 2,"RE",nil
}
