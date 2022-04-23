package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

var version string

type HandleFnc func(http.ResponseWriter, *http.Request)

func main() {
	// 获取系统变量VERSION
	version = os.Getenv("VERSION ")
	// 开启服务监听
	startServer()
}

func startServer() {
	// 注册处理器
	http.HandleFunc("/", checkAndLog(HelloServer))
	http.HandleFunc("/healthz", checkAndLog(successReply))
	// 监听本地80端口
	err := http.ListenAndServe("0.0.0.0:80", nil)
	if nil != err {
		log.Fatal("ListenAndServe: ", err.Error())
	}
}

// 打印客户端IP及recover panic
func checkAndLog(f HandleFnc) HandleFnc {
	return func(writer http.ResponseWriter, request *http.Request) {
		defer func() {
			if x := recover(); x != nil {
				log.Printf("[%v] caught panic: %v\n", request.RemoteAddr, x)
			}
		}()
		log.Printf("Client IP: %s\n", request.RemoteAddr)
		f(writer, request)
	}
}

// 处理healthz请求
func successReply(w http.ResponseWriter, req *http.Request) {
	// 设置http code
	w.WriteHeader(http.StatusBadGateway)
	body := fmt.Sprintf("Hello World! \n%d", http.StatusOK)
	// 写入请求体
	w.Write([]byte(body))
	logReply(http.StatusOK, body)
}

func HelloServer(w http.ResponseWriter, req *http.Request) {
	// 遍历request header并写入到response header中
	for key, val := range req.Header {
		w.Header()[key] = val
	}
	// 写入version
	w.Header().Set("VERSION", version)
	body := fmt.Sprintf("Hello World! \n%s%s", req.RemoteAddr, req.RequestURI)
	fmt.Fprint(w, body)
	logReply(http.StatusOK, body)
}

func logReply(code int, body string) {
	log.Printf("HTTP ReturnCode: %d\tResponse Body：%s\n", code, body)
}
