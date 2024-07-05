package database

import "time"

type RefreshToken struct {
	UserID    int       `json:"user_id"`
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
}


func (db *DB) SaveRefreshToken(userID int, token string) error {
	dbStructure, err := db.loadDB()
	if err != nil {
		return err
	}

	refreshToken := RefreshToken{
		UserID:    userID,
		Token:     token,
		ExpiresAt: time.Now().Add(time.Hour),
	}

	dbStructure.RefreshTokens[token] = refreshToken

	err = db.writeDB(dbStructure)
	if err != nil {
		return err
	}

	return nil
}

// 从 数据库中删除 refreshToken
func (db *DB) RevokeRefreshToken(token string) error {
	dbStructure, err := db.loadDB()
	if err != nil {
		return err
	}
	delete(dbStructure.RefreshTokens, token)
	err = db.writeDB(dbStructure)
	if err != nil {
		return err
	}
	return nil
}

// UserForRefreshToken 函数通过 refreshToken 找到 user
// db 是一个 DB 类型的接收器，表示数据库对象。
func (db *DB) UserForRefreshToken(token string) (User, error) {
	//  加载数据库结构
	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	// 从数据库结构中获取与 token 对应的 refreshToken
	refreshToken, ok := dbStructure.RefreshTokens[token]
	if !ok {
		return User{}, ErrNotExist
	}

	// 检查 refreshToken 是否已经过期
	// 检查刷新令牌的过期时间是否在当前时间之前。
	if refreshToken.ExpiresAt.Before(time.Now()) {
		return User{}, ErrNotExist
	}

	// 根据 refreshToken 中的 UserID 获取用户信息
	user, err := db.GetUser(refreshToken.UserID)
	if err != nil {
		return User{}, err
	}

	return user, nil

}
