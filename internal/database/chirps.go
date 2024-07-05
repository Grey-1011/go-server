package database

// Chirp 结构体表示一个 chirp（类似 tweet），包含两个字段：
type Chirp struct {
	ID   int    `json:"id"`
	Body string `json:"body"`
	AuthorID int `json:"author_id"`
}


// ==== 创建 Chirp ====
// CreateChirp 方法创建一个新的 chirp 并保存到数据库中。
/*
1) 首先加载当前的数据库结构。
2) 创建一个新的 Chirp，分配一个唯一 ID。
3) 将新的 Chirp 添加到 dbStructure.Chirps 映射中。
4) 将更新后的数据库结构写回文件。
*/
func (db *DB) CreateChirp(body string, authorID int) (Chirp, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}

	id := len(dbStructure.Chirps) + 1
	chirp := Chirp{
		ID:   id,
		Body: body,
		AuthorID: authorID,
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



// 获取指定 ID 的 Chirp
func (db *DB) GetChirp(id int) (Chirp, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}

	chirp, ok := dbStructure.Chirps[id]
	if !ok {
		return Chirp{}, ErrNotExist
	}

	return chirp, nil
}
