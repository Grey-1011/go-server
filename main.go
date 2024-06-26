package main

import (
	"fmt"
	"log"
	"net/http"
)

type apiConfig struct {
	fileserverHits int
}

func main() {
	const port = "8080"

	apiCfg := apiConfig{
		fileserverHits: 0,
	}

	// create a  new http.ServeMux
	/*
		http.NewServeMux() 创建了一个新的 ServeMux 实例， 这是一个 HTTP 请求的路由器。
		路由器的作用是根据请求的 URL 将请求分配给不同的处理程序（handler）。
		你可以将不同的 URL 路径和对应的处理函数注册到这个路由器上。
	*/
	mux := http.NewServeMux()

	// 使用 middlewareMetricsInc 中间件包装文件服务器处理程序
	mux.Handle("/app/*", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(".")))))

	mux.HandleFunc("/healthz", handlerReadiness)
	// 注册 /metrics 处理程序
	mux.HandleFunc("/metrics", apiCfg.handlerMetrics)
	// 注册 /reset 处理程序
	mux.HandleFunc("/reset", apiCfg.handlerReset)

	/*
		使用 &符号创建一个指向 http.Server 结构体的指针。
		这允许在其他函数和方法中使用这个指针来引用和修改同一个服务器实例。
	*/
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	// start the server
	srv.ListenAndServe()

	log.Printf("Serving on port: %s\n", port)
	log.Fatal(srv.ListenAndServe())

}

// /metrics 处理程序, 返回请求计数
func (cfg *apiConfig) handlerMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Hits: %d", cfg.fileserverHits)))
}

// middlewareMetricsInc 中间件方法，递增计数器
func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits++
		next.ServeHTTP(w, r)
	})
}
