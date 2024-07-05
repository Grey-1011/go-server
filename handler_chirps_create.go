package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/Grey-1011/go-server/internal/auth"
)

/*
在 Go 中，只有导出的结构体字段（即以大写字母开头的字段）
才可以使用 encoding/json 包进行编码或解码。
*/
type Chirp struct {
	ID       int    `json:"id"`
	Body     string `json:"body"` // 注意json 后没有空格
	AuthorID int    `json:"author_id"`
}

func (cfg *apiConfig) handlerChirpsCreate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find JWT")
		return
	}
	subject, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Could't validate JWT")
		return
	}
	userID, err := strconv.Atoi(subject)
	if err != nil {
		respondWithJSON(w, http.StatusBadRequest, "Couldn't parse user ID")
		return
	}

	// 创建一个 JSON 解码器来解析请求体
	decoder := json.NewDecoder(r.Body)
	params := parameters{}

	// 解析请求体中的 JSON 数据并存储到 `params` 结构体中
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	cleaned, err := validateChirp(params.Body)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// 创建 Chirp ,  需要 userID
	chirp, err := cfg.DB.CreateChirp(cleaned, userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create chirp")
		return
	}

	respondWithJSON(w, http.StatusCreated, Chirp{
		ID:       chirp.ID,
		Body:     chirp.Body,
		AuthorID: chirp.AuthorID,
	})

	// 如果 Chirp 合法，则返回成功响应
	// respondWithJSON(w, http.StatusCreated, cleaned)
}

func validateChirp(body string) (string, error) {
	const maxChirpLength = 140
	if len(body) > maxChirpLength {
		return "", errors.New("Chirp is too long")
	}

	badWords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}
	cleaned := getCleanedBody(body, badWords)
	return cleaned, nil
}

// 定义了 cleanProfaneWords 函数来替换不雅词语。
func getCleanedBody(body string, badWords map[string]struct{}) string {
	words := strings.Split(body, " ")
	for i, word := range words {
		// cleanedWord := strings.Trim(word, ".,!?:;") // 去除标点符号
		loweredWord := strings.ToLower(word)
		if _, ok := badWords[loweredWord]; ok {
			words[i] = "****"
		}
	}
	return strings.Join(words, " ")
}
