package main

import (
    "fmt"
    "io/ioutil"
    "net/http"
)

func uploadHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method == "POST" {
        // 处理代码输入
        code := r.FormValue("code")
        codeFilePath := "Sandbox/file.go"
        err := ioutil.WriteFile(codeFilePath, []byte(code), 0644)
        if err != nil {
            fmt.Println("保存代码失败:", err)
            return
        }

        fmt.Fprintln(w, "代码保存成功")
    } else {
        http.ServeFile(w, r, "Website/upload.html")
    }
}

func main() {
    http.HandleFunc("/upload", uploadHandler)
    http.ListenAndServe(":8080", nil)
}