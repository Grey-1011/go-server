package database

import (
	"encoding/json"
	"errors"
	"os" // os 用于文件操作。
	"sync"
)

var ErrNotExist = errors.New("resource does not exist")

// 数据库结构体 DB
type DB struct {
	path string        // 数据库文件的路径。
	mu   *sync.RWMutex // 读写锁（RWMutex），用于确保并发安全。
}

// 数据库的内部结构，包含一个 Chirps 映射
type DBStructure struct {
	Chirps        map[int]Chirp           `json:"chirps"`
	Users         map[int]User            `json:"users"`
	RefreshTokens map[string]RefreshToken `json:"refresh_tokens"`
}

// ==== 创建新数据库 ====
/*
NewDB 函数创建一个新的 DB 对象，确保数据库文件存在。
如果文件不存在，调用 ensureDB 函数创建它。
*/
func NewDB(path string) (*DB, error) {
	db := &DB{
		path: path,
		mu:   &sync.RWMutex{},
	}
	err := db.ensureDB()
	return db, err
}

// ==== 创建数据库文件 ====
/*
createDB 方法创建一个新的空数据库文件，包含一个空的 Chirps 映射，并将其写回文件。
*/
func (db *DB) createDB() error {
	dbStructure := DBStructure{
		Chirps:        map[int]Chirp{},
		Users:         map[int]User{},
		RefreshTokens: map[string]RefreshToken{},
	}
	return db.writeDB(dbStructure)
}

//  确保数据库文件存在
/*
ensureDB 方法检查数据库文件是否存在，如果不存在，则创建一个新的数据库文件。
*/
func (db *DB) ensureDB() error {
	_, err := os.ReadFile(db.path)
	if errors.Is(err, os.ErrNotExist) {
		return db.createDB()
	}
	return err
}

// 重置数据库
func (db *DB) ResetDB() error {
	err := os.Remove(db.path)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	return db.ensureDB()
}

// ==== 加载数据库 ====
/*
1) 使用读锁确保并发安全。
2) 读取文件内容。
3) 将文件内容解码为 DBStructure。
*/
func (db *DB) loadDB() (DBStructure, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	dbStructure := DBStructure{}
	dat, err := os.ReadFile(db.path)
	if errors.Is(err, os.ErrNotExist) {
		return dbStructure, err
	}
	err = json.Unmarshal(dat, &dbStructure) // 解码
	if err != nil {
		return dbStructure, err
	}

	return dbStructure, nil
}

// ==== 写入数据库 ====
/*
1) 使用写锁确保并发安全。
2) 将数据库结构编码为 JSON。
3) 将 JSON 数据写入文件。
*/
func (db *DB) writeDB(dbStructure DBStructure) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	dat, err := json.Marshal(dbStructure) // 编码
	if err != nil {
		return err
	}

	err = os.WriteFile(db.path, dat, 0600)
	if err != nil {
		return err
	}
	return nil
}
