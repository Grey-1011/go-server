package main

import (
	"log"
	"net/http"
)

func main() {
	const port = "8080"

	// create a  new http.ServeMux
	/*
	http.NewServeMux() 创建了一个新的 ServeMux 实例， 这是一个 HTTP 请求的路由器。
	路由器的作用是根据请求的 URL 将请求分配给不同的处理程序（handler）。
	你可以将不同的 URL 路径和对应的处理函数注册到这个路由器上。
	*/
	mux := http.NewServeMux()
	
	/*
	使用 &符号创建一个指向 http.Server 结构体的指针。
	这允许在其他函数和方法中使用这个指针来引用和修改同一个服务器实例。
	*/
	srv := &http.Server{
		Addr:  ":" + port,
		Handler: mux,
	}
	// start the server
	srv.ListenAndServe()

	log.Printf("Serving on port: %s\n", port)
	log.Fatal(srv.ListenAndServe())

}