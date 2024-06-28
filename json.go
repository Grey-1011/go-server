package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func respondWithError(w http.ResponseWriter, code int, msg string) {
	if code > 499 {
		log.Printf("Responding with 5XX error: %s", msg)
	}
	type errorResponse struct {
		Error string `json:"error"`
	}
	// w.Header().Set("Content-Type", "application/json")
	// w.WriteHeader(code)
	// // 将错误消息编码为 JSON 并写入响应
	// json.NewEncoder(w).Encode(errorResponse{Error: msg})

	respondWithJSON(w, code, errorResponse{
		Error: msg,
	})
}



func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	dat, err := json.Marshal(payload) //go数据 编码为 JSON 数据
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return		
	}
	// // 将成功消息编码为 JSON 并写入响应
	// json.NewEncoder(w).Encode(payload)
	w.WriteHeader(code)
	w.Write(dat)
}
