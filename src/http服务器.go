package main

import (
	"fmt"
	"net/http"
)

//writer  给客户端回复数据
//request  读取客户端发送的数据
func HandConn(writer http.ResponseWriter, request *http.Request) {
	fmt.Println("request.Method = ", request.Method)		//请求方法
	fmt.Println("request.URL = ", request.URL)		//请求地址
	fmt.Println("request.Header = ", request.Header)		//请求头
	fmt.Println("request.Body = ", request.Body)		//请求主体

	writer.Write([]byte("你好 我是李歘歘"))	//给客户端回复数据
}

func main()  {
	//注册处理函数，用户连接，自动调用指定的处理函数
	http.HandleFunc("/",HandConn)

	//监听绑定
	http.ListenAndServe(":8000",nil)

}


