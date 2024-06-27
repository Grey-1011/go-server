package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

type apiConfig struct {
	fileserverHits int
}

/*
在 Go 中，只有导出的结构体字段（即以大写字母开头的字段）
才可以使用 encoding/json 包进行编码或解码。
*/
type validateChirRequest struct {
	Body string  `json:"body"`  // 注意json 后没有空格
}
type errorResponse struct {
	Error string `json:"error"`
}
// type validResponse struct {
// 	Valid bool `json:"valid"`
// }

var profaneWords = []string{"kerfuffle", "sharbert", "fornax"}

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

	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	// 注册 /metrics 处理程序
	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)
	// 注册 /reset 处理程序
	mux.HandleFunc("GET /api/reset", apiCfg.handlerReset)

	// 我们定义了一个路由规则，将 POST 请求映射到 /api/validate_chirp 处理函数 handlerValidateChirp：
	mux.HandleFunc("POST /api/validate_chrip", handlerValidateChirp)
	
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
	w.Header().Add("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`
		<html>

			<body>
    		<h1>Welcome, Chirpy Admin</h1>
    		<p>Chirpy has been visited %d times!</p>
			</body>

		</html>`, cfg.fileserverHits)))
}

// middlewareMetricsInc 中间件方法，递增计数器
func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits++
		next.ServeHTTP(w, r)
	})
}



func handlerValidateChirp(w http.ResponseWriter, r *http.Request) {
	// 创建一个 JSON 解码器来解析请求体
	decoder := json.NewDecoder(r.Body)
	req := validateChirRequest{}
	
	// 解析请求体中的 JSON 数据并存储到 `req` 结构体中
	err := decoder.Decode(&req)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}
	// 检查 Chirp 是否超过 140 个字符
	if len(req.Body) > 140 {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}

	cleanedBody := cleanProfaneWords(req.Body)
	// 如果 Chirp 合法，则返回成功响应
	respondWithJSON(w, http.StatusOK, cleanedBody)
}



func respondWithError(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	// 将错误消息编码为 JSON 并写入响应
	json.NewEncoder(w).Encode(errorResponse{Error: message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	// 将成功消息编码为 JSON 并写入响应
	json.NewEncoder(w).Encode(payload)
}


// 定义了 cleanProfaneWords 函数来替换不雅词语。
func cleanProfaneWords(body string) string {
	words := strings.Split(body, " ")
	for i, word := range words {
		cleanedWord := strings.Trim(word, ".,!?:;") // 去除标点符号
		for _, profaneWord := range profaneWords {
			if strings.ToLower(cleanedWord) == profaneWord {
				words[i] = strings.Replace(words[i], cleanedWord, "****", -1)
			}
		}
	}
	return strings.Join(words, " ")
}