package main

import (
	"net/http"
	"time"

	"github.com/Grey-1011/go-server/internal/auth"
)

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Token string `json:"token"`
	}

	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't find token")
		return
	}

	// UserForRefreshToken 函数通过 refreshToken 找到 user
	user, err := cfg.DB.UserForRefreshToken(refreshToken)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't get user for refresh token")
		return
	}

	// 创建新的 accessToken
	accessToken, err := auth.MakeJWT(
		user.ID,
		cfg.jwtSecret,
		time.Hour,
	)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate token")
		return
	}

	// 响应返回 accessToken
	respondWithJSON(w, http.StatusOK, response{
		Token: accessToken,
	})

}


// 撤销与请求头中传递的 refreshToken 匹配的数据库中的 refreshToken
func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't find token")
		return
	}

	err = cfg.DB.RevokeRefreshToken(refreshToken)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't revoke session")
		return
	}

	// 响应 204 : 204状态码表示请求成功但不返回任何内容。
	w.WriteHeader(http.StatusNoContent)
}