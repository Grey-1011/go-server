package database

import (
	"encoding/json"
	"errors"
	"os"    // os 用于文件操作。
	"sync"
)

// 数据库结构体 DB
type DB struct {
	path string      // 数据库文件的路径。
	mu   *sync.RWMutex  // 读写锁（RWMutex），用于确保并发安全。
}

// 数据库的内部结构，包含一个 Chirps 映射
type DBStructure struct {
	Chirps map[int]Chirp `json:"chirps"`
}

// Chirp 结构体表示一个 chirp（类似 tweet），包含两个字段：
type Chirp struct {
	ID   int    `json:"id"`
	Body string `json:"body"`
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


// ==== 创建 Chirp ====
// CreateChirp 方法创建一个新的 chirp 并保存到数据库中。
/*
1) 首先加载当前的数据库结构。
2) 创建一个新的 Chirp，分配一个唯一 ID。
3) 将新的 Chirp 添加到 dbStructure.Chirps 映射中。
4) 将更新后的数据库结构写回文件。
*/
func (db *DB) CreateChirp(body string) (Chirp, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}

	id := len(dbStructure.Chirps) + 1
	chirp := Chirp{
		ID:   id,
		Body: body,
	}
	dbStructure.Chirps[id] = chirp

	err = db.writeDB(dbStructure)
	if err != nil {
		return Chirp{}, err
	}

	return chirp, nil
}


// ==== 获取所有 Chirps ====
/*
1) 加载数据库结构。
2) 遍历 dbStructure.Chirps 映射，将所有的 Chirp 添加到一个切片中。
3) 返回这个切片。
*/
func (db *DB) GetChirps() ([]Chirp, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return nil, err
	}

	chirps := make([]Chirp, 0, len(dbStructure.Chirps))
	for _, chirp := range dbStructure.Chirps {
		chirps = append(chirps, chirp)
	}

	return chirps, nil
}



// ==== 创建数据库文件 ====
/*
createDB 方法创建一个新的空数据库文件，包含一个空的 Chirps 映射，并将其写回文件。
*/
func (db *DB) createDB() error {
	dbStructure := DBStructure{
		Chirps: map[int]Chirp{},
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
	err = json.Unmarshal(dat, &dbStructure)
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

	dat, err := json.Marshal(dbStructure)
	if err != nil {
		return err
	}

	err = os.WriteFile(db.path, dat, 0600)
	if err != nil {
		return err
	}
	return nil
}