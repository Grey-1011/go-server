package main

import (
	"net/http"
	"sort"
	"strconv"
)

func (cfg *apiConfig) handlerChirpsGet(w http.ResponseWriter, r *http.Request) {
	chirpIdString := r.PathValue("chirpID")
	chirpId, err := strconv.Atoi(chirpIdString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid chirp ID")
		return
	}

	dbChirp, err := cfg.DB.GetChirp(chirpId)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't get chirp")
		return
	}

	respondWithJSON(w, http.StatusOK, Chirp{
		ID:   dbChirp.ID,
		Body: dbChirp.Body,
	})
}


func (cfg *apiConfig) handlerChirpsRetrieve(w http.ResponseWriter, r *http.Request) {
	// 获取所有 Chirps
	dbChirps, err := cfg.DB.GetChirps()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't retrieve chirps")
		return
	}

	// 初始化 authorID 为 -1（表示没有指定特定作者）。
	authorID := -1
	authorIDString := r.URL.Query().Get("author_id")
	// 如果 author_id 存在
	if authorIDString != "" {
		authorID, err = strconv.Atoi(authorIDString)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid author ID")
			return
		}
	}

	chirps := []Chirp{}
	for _, dbChirp := range dbChirps {
		// 如果指定了 author_id 且 chirp 的 AuthorID 不匹配，跳过该 chirp。
		if authorID != -1 && dbChirp.AuthorID != authorID {
			continue
		}

		chirps = append(chirps, Chirp{
			ID:   dbChirp.ID,
			AuthorID:  dbChirp.AuthorID,
			Body: dbChirp.Body,
		})
	}

	// 按 id 升序排列
	sort.Slice(chirps, func(i, j int) bool {
		return chirps[i].ID < chirps[j].ID
	})

	respondWithJSON(w, http.StatusOK, chirps)
}
