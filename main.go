package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/Grey-1011/go-server/internal/database"
	"github.com/joho/godotenv"
)

type apiConfig struct {
	fileserverHits int
	DB             *database.DB
	jwtSecret      string
}

func main() {
	const filepathRoot = "."
	const port = "8080"

	godotenv.Load(".env")

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET enviroment variable is not set")
	}

	// 创建新数据库
	db, err := database.NewDB("database.json")
	if err != nil {
		log.Fatal(err)
	}

	dbg := flag.Bool("debug", false, "Enable debug mode")
	flag.Parse()
	if dbg != nil && *dbg {
		err := db.ResetDB()
		if err != nil {
			log.Fatal(err)
		}
	}
	// Use: go build -o out && ./out --debug

	apiCfg := apiConfig{
		fileserverHits: 0,
		DB:             db,
		jwtSecret:      jwtSecret,
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
	// 根据 ID 获取 Chirps
	mux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.handlerChirpsGet)

	mux.HandleFunc("POST /api/users", apiCfg.handlerUsersCreate)
	// 更新用户的电子邮件和密码
	mux.HandleFunc("PUT /api/users", apiCfg.handlerUsersUpdate)
	mux.HandleFunc("POST /api/login", apiCfg.handlerLogin)

	mux.HandleFunc("POST /api/refresh", apiCfg.handlerRefresh)
	mux.HandleFunc("POST /api/revoke", apiCfg.handlerRevoke)

	mux.HandleFunc("DELETE /api/chirps/{chirpID}", apiCfg.handlerChirpsDelete)
	mux.HandleFunc("POST /api/polka/webhooks", apiCfg.handlerWebhook)

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
