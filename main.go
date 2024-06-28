package main

import (
	"log"
	"net/http"

	"github.com/Grey-1011/go-server/internal/database"
)

type apiConfig struct {
	fileserverHits int
	DB             *database.DB
}

func main() {
	const filepathRoot = "."
	const port = "8080"

	// 创建新数据库
	db, err := database.NewDB("database.json")
	if err != nil {
		log.Fatal(err)
	}

	apiCfg := apiConfig{
		fileserverHits: 0,
		DB:             db,
	}

	// create a  new http.ServeMux
	/*
		http.NewServeMux() 创建了一个新的 ServeMux 实例， 这是一个 HTTP 请求的路由器。
		路由器的作用是根据请求的 URL 将请求分配给不同的处理程序（handler）。
		你可以将不同的 URL 路径和对应的处理函数注册到这个路由器上。
	*/
	mux := http.NewServeMux()

	fsHandler := apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))))
	// 使用 middlewareMetricsInc 中间件包装文件服务器处理程序
	mux.Handle("/app/*", fsHandler)

	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	// 注册 /metrics 处理程序
	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)
	// 注册 /reset 处理程序
	mux.HandleFunc("GET /api/reset", apiCfg.handlerReset)

	// 我们定义了一个路由规则，将 POST 请求映射到 /api/validate_chirp 处理函数 handlerValidateChirp：
	mux.HandleFunc("POST /api/chirps", apiCfg.handlerChirpsCreate)

	// handlerChirpsRetrieve 获取所有 Chirps
	mux.HandleFunc("GET /api/chirps", apiCfg.handlerChirpsRetrieve)
	/*
		使用 &符号创建一个指向 http.Server 结构体的指针。
		这允许在其他函数和方法中使用这个指针来引用和修改同一个服务器实例。
	*/
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving on port: %s\n", port)
	// start the server
	log.Fatal(srv.ListenAndServe())

}
