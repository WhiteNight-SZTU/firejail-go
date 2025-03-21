package main

import (
	"fmt"
	"strings"
	//"github.com/creack/pty"
	"bufio"
	"os/exec"
	"regexp"
	"io"
	"os"
	"time"
    "strconv"
)

type Firejail struct {
	cpu string
	mem string
    fsize string //可创建的文件大小
	profile string
}

func (f Firejail)setCommand(path string,filename string)string{
    command := "firejail " + " --rlimit-cpu=" + f.cpu + " --rlimit-as=" + f.mem 
    //command += " --quiet"
    if (f.profile != "") {
        command += " --profile=" + f.profile
    }

    command += " /usr/bin/time -v "
    if (strings.HasSuffix(filename, ".py")) {
        command += " python3 "
    }
    command += " " + path + "/" + filename
    return command
}

/* 处理time指令的输出
* @param output:用户代码输出
* @return:进程PID 用户态时长 系统态时长 运行时长 内存占用量
*/
func monitorResult(output string) (string,string,string,string,string,error) {
    
    pidRe := regexp.MustCompile(`Parent pid \d+, child pid (\d+)`) 
    userTimeRe := regexp.MustCompile(`User time \(seconds\): (\d+\.\d+)`)
    sysTimeRe := regexp.MustCompile(`System time \(seconds\): (\d+\.\d+)`)
    elapsedRe := regexp.MustCompile(`Elapsed \(wall clock\) time \(h:mm:ss or m:ss\): (\d+:\d+\.\d+)`)
    maxMemRe := regexp.MustCompile(`Maximum resident set size \(kbytes\): (\d+)`)

    pidMatch := pidRe.FindStringSubmatch(output)
    userTimeMatch := userTimeRe.FindStringSubmatch(output)
    sysTimeMatch := sysTimeRe.FindStringSubmatch(output)
    elapsedMatch := elapsedRe.FindStringSubmatch(output)
    maxMemMatch := maxMemRe.FindStringSubmatch(output)

    maxMemKB,_ := strconv.Atoi(maxMemMatch[1])
    maxMemMB := float64(maxMemKB) / 1024.0

    return pidMatch[1],userTimeMatch[1],sysTimeMatch[1],elapsedMatch[1],strconv.FormatFloat(maxMemMB,'f', 2, 64),nil
}

//有标准输入的用户代码
func (f Firejail)runFirejailWithSTDIN(path string, filename string)(error) {
    command := f.setCommand(path,filename)
    logger.Debug("Sandbox","Executing command: ", command)
    startTime := time.Now()

    output,err,outputErr := RunCommandWithInput(command,filename)
    //If the start time greater than the f.cpu (cpu time),return TRL
    cpuLimitTime,err := strconv.Atoi(f.cpu)
    if err != nil {
        logger.Error("Sandbox","Execute Firejail error: ",err)
        return err
    }
    logger.Debug("Sandbox","cpuLimitTime: ",float64(cpuLimitTime))
    logger.Debug("Sandbox","startTime: ",time.Since(startTime).Seconds())
    if time.Since(startTime).Seconds() > float64(cpuLimitTime) {
        logger.Error("Sandbox",filename + " Result: TRL")
    }
    if err != nil {
        logger.Error("Sandbox","Execute Firejail error: ",err)
        return err
    }
    // firejail与time指令的完整输出
    //logger.Debug("Sandbox","Execute Firejail: ",outputErr)
    logger.UserOutput("Sandbox","Execute Firejail output: ",output)
    statuCode,result,err := CompareAnswer(output,filename)
    logger.Info("Sandbox",filename," StatuCode: ",statuCode," Result: ",result)
    if err != nil {
        logger.Info("Sandbox","Execute Firejail error: ",err)
        return err
    }
    pid,userTime,sysTime,elapsed,maxMem,err := monitorResult(outputErr)
    if err != nil {
        logger.Error("Sandbox","Monitor Firejail error: ",err)
        return err
    }
    logger.Info("Sandbox","PID: ",pid," UserTime: ",userTime," SysTime: ",sysTime," Elapsed: ",elapsed," MaxMem: ",maxMem)
    return nil
}

//无标准输入的用户代码
func (f Firejail)runFirejail(path string, filename string)error {
    command := f.setCommand(path,filename)
    logger.Debug("Sandbox","Executing command: ", command)
    startTime := time.Now()
    output,err,outputErr := RunCommand(command)
    //If the start time greater than the f.cpu (cpu time),return TRL
    cpuLimitTime,err := strconv.Atoi(f.cpu)
    if err != nil {
        logger.Error("Sandbox","Execute Firejail error: ",err)
        return err
    }
    logger.Debug("Sandbox","cpuLimitTime: ",float64(cpuLimitTime))
    logger.Debug("Sandbox","startTime: ",time.Since(startTime).Seconds())

    if time.Since(startTime).Seconds() > float64(cpuLimitTime) {
        logger.Error("Sandbox",filename + " Result: TRL")
        fmt.Println(filename + " Result: TRL")
    }
    if err != nil {
        logger.Error("Sandbox","Execute Firejail error: ",err)
        return err
    }
    // firejail与time指令的完整输出
    logger.Debug("Sandbox","Execute Firejail: ",outputErr)
    logger.UserOutput("Sandbox","Execute Firejail output: ",output)
    statuCode,result,err := CompareAnswer(output,filename)
    logger.Info("Sandbox",filename," StatuCode: ",statuCode," Result: ",result)
    if err != nil {
        logger.Info("Sandbox","Execute Firejail error: ",err)
    }
    pid,userTime,sysTime,elapsed,maxMem,err := monitorResult(outputErr)
    if err != nil {
        logger.Error("Sandbox","Monitor Firejail error: ",err)
        return err
    }
    logger.Info("Sandbox","PID: ",pid," UserTime: ",userTime," SysTime: ",sysTime," Elapsed: ",elapsed," MaxMem: ",maxMem)
    return nil
}

func (f Firejail) runFirejailWithCommunication(path_A string, filename_A string, path_B string, filename_B string) error {
    // 创建双向通信管道
    aToBReader, aToBWriter := io.Pipe()
    bToAReader, bToAWriter := io.Pipe()

    // 构建程序A的沙箱命令
    cmdA := exec.Command(
        "firejail",
        "--"+f.profile,
		"--ipc-namespace",
        path_A+"/"+filename_A,
    )
    cmdA.Stdin = bToAReader   // A从B的输出读取
    cmdA.Stdout = aToBWriter  // A的输出发送给B
    cmdA.Stderr = os.Stderr
    cmdA.Env = append(os.Environ(), "NOTE=TEST")

    // 构建程序B的沙箱命令
    cmdB := exec.Command(
        "firejail",
        "--"+f.profile,
		"--ipc-namespace",
        path_B+"/"+filename_B,
    )
    cmdB.Stdin = aToBReader   // B从A的输出读取
    cmdB.Stdout = bToAWriter  // B的输出发送给A
    cmdB.Stderr = os.Stderr
    cmdB.Env = append(os.Environ(), "NOTE=TEST")

    // 启动输出监控
    go logPipeOutput(aToBReader, "A→B") // 记录A发给B的消息
    go logPipeOutput(bToAReader, "B→A") // 记录B发给A的消息

    // 启动双进程
    logger.Debug("Sandbox", "正在启动进程A:", filename_A)
    if err := cmdA.Start(); err != nil {
        logger.Error("Sandbox", "进程A启动失败:", err)
        return fmt.Errorf("processA start failed: %v", err)
    }
    
    logger.Debug("Sandbox", "正在启动进程B:", filename_B)
    if err := cmdB.Start(); err != nil {
        logger.Error("Sandbox", "进程B启动失败:", err)
        cmdA.Process.Kill() // 清理已启动的进程A
        return fmt.Errorf("processB start failed: %v", err)
    }

    // 异步等待进程结束
    done := make(chan error, 2)
    go func() {
        err := cmdA.Wait()
        aToBWriter.Close() // 确保管道关闭
        done <- err
    }()
    go func() {
        err := cmdB.Wait()
        bToAWriter.Close() // 确保管道关闭
        done <- err
    }()

    // 带超时的等待机制
    timeout := time.After(30 * time.Second)
    var finalErr error
    
    select {
    case err := <-done:
        if err != nil {
            finalErr = fmt.Errorf("进程异常退出: %v", err)
            logger.Error("Sandbox", finalErr.Error())
        }
    case <-timeout:
        cmdA.Process.Kill()
        cmdB.Process.Kill()
        finalErr = fmt.Errorf("执行超时(30秒)")
        logger.Error("Sandbox", finalErr.Error())
    }

    // 最终状态记录
    logger.Debug("Sandbox", "通信会话结束",
        "status:", map[string]interface{}{
            "processA": cmdA.ProcessState,
            "processB": cmdB.ProcessState,
        })
    
    return finalErr
}

// 独立输出记录函数
func logPipeOutput(reader io.Reader, direction string) {
    scanner := bufio.NewScanner(reader)
    for scanner.Scan() {
        line := strings.TrimSpace(scanner.Text())
        if line != "" {
            logger.UserOutput("COMM", fmt.Sprintf("[%s] %s", direction, line))
        }
    }
    if err := scanner.Err(); err != nil {
        logger.Error("PIPE", "管道读取错误:", 
            "direction:", direction,
            "error:", err)
    }
}

func test_cpu(){
	path,filename := ComplieFile("Sandbox/test_cpu.go")
	if filename == "" {
		err := fmt.Errorf("Cant't find compiled file: %s in %s",filename,path)
		logger.Error("Sandbox",err)
		return 
	}
	firejail := Firejail{cpu: "30", mem: "2g", profile: "noprofile"}
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
	firejail := Firejail{cpu: "30", mem: "3g", profile: "noprofile"}
	err := firejail.runFirejail(path,filename)
	if err != nil {
		logger.Error("Sandbox","Execute Firejail error: ",err)
		return
	}
}

func test_communication(){
	path_A,filename_A := ComplieFile("Sandbox/test_A.go")
	path_B,filename_B := ComplieFile("Sandbox/test_B.go")
	if filename_A == "" {
		err := fmt.Errorf("Cant't find compiled file: %s in %s",filename_A,path_A)
		logger.Error("Sandbox",err)
		return 
	}
	if filename_B == "" {
		err := fmt.Errorf("Cant't find compiled file: %s in %s",filename_B,path_B)
		logger.Error("Sandbox",err)
		return 
	}
	firejail := Firejail{cpu: "60", mem: "null", profile: "noprofile"}
	err := firejail.runFirejailWithCommunication(path_A,filename_A,path_B,filename_B)
	if err != nil {
		logger.Error("Sandbox","Execute Firejail error: ",err)
		return
	}
}

func test_python(){
    path,filename := ComplieFile("Sandbox/test_py.py")
    if filename == "" {
        err := fmt.Errorf("Cant't find compiled file: %s in %s",filename,path)
        logger.Error("Sandbox",err)
        return 
    }
    firejail := Firejail{cpu: "30", mem: "2g", profile: "/etc/firejail/wn_run.profile "}
    err := firejail.runFirejail(path,filename)
    if err != nil {
        logger.Error("Sandbox","Execute Firejail error: ",err)
        return
    }
}


//C++：标准输入输出-单组输入测试
func test_cpp(){
    path,filename := ComplieFile("Sandbox/cpp/test_cpp.cpp")
    if filename == "" {
        err := fmt.Errorf("Cant't find compiled file: %s in %s",filename,path)
        logger.Error("Sandbox",err)
        return 
    }
    firejail := Firejail{cpu: "10", mem: "2g",fsize: "512k",
    profile: "/etc/firejail/wn_run.profile"}
    firejail.runFirejailWithSTDIN(path,filename)
}

//C++：无标准输入输出-死循环测试
func test_cpp_loop(){
    path,filename := ComplieFile("Sandbox/cpp/test_cpp_loop.cpp")
    if filename == "" {
        err := fmt.Errorf("Cant't find compiled file: %s in %s",filename,path)
        logger.Error("Sandbox",err)
        return 
    }
    firejail := Firejail{cpu: "10", mem: "2g",fsize: "512k",
    profile: "/etc/firejail/wn_run.profile"}
    firejail.runFirejail(path,filename)
}

//C++：无标准输入输出-文件读权限测试
func test_cpp_read(){
    path,filename := ComplieFile("Sandbox/cpp/test_cpp_read.cpp")
    if filename == "" {
        err := fmt.Errorf("Cant't find compiled file")
        logger.Error("Sandbox",err)
        return 
    }
    firejail := Firejail{cpu: "10", mem: "2g",fsize: "512k",
    profile: "/etc/firejail/wn_run.profile"}
    firejail.runFirejail(path,filename)
}

//C++：无标准输入输出-执行权限测试（常见指令，增删改查文件，重启）
func test_cpp_exec(){
    path,filename := ComplieFile("Sandbox/cpp/test_cpp_exec.cpp")
    if filename == "" {
        err := fmt.Errorf("Cant't find compiled file")
        logger.Error("Sandbox",err)
        return 
    }
    firejail := Firejail{cpu: "10", mem: "2g",fsize: "512k",
    profile: "/etc/firejail/wn_run.profile"}
    firejail.runFirejail(path,filename)
}

func main() {
    InitLogger("debug")
    //test_cpu()
    //test_memory()
	//test_communication()
    //test_python()
    
    //test_cpp()
    //test_cpp_loop()
    //test_cpp_read()
    test_cpp_exec()


    //clear()
    exit()
}

func clear(){
    files, _ := os.ReadDir("Sandbox/cpp")
    for _, file := range files {
        if !file.IsDir() && !strings.Contains(file.Name(), ".") {
            filePath := "Sandbox/cpp" + "/" + file.Name()
            os.Remove(filePath)
        }
    }
}

func exit() {
    logger.Info("Sandbox", "程序退出\n\n\n\n\n")
}