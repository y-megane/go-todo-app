package todo

import (
	"context"
	"sync"
	"time"

	"github.com/google/uuid"
)

// DBインタフェース定義
// PutメソッドとGetALLメソッドを実装したStructは「DBインタフェースを実装」したことになる
type DB interface {
	Put(ctx context.Context, todo *TODO) error
	GetAll(ctx context.Context) ([]*TODO, error)
}
type MemoryDB struct {
	// 排他制御用のMutex
	sync.RWMutex
	// 文字列をキー、TODO型を値とするmap
	// 扱う値はポインタで定義する
	m map[string]*TODO
}

var _ DB = (*MemoryDB)(nil)

// ファクトリメソッド
func NewMemoryDB() *MemoryDB {
	return &MemoryDB{m: map[string]*TODO{}}
}

func (db *MemoryDB) Put(ctx context.Context, todo *TODO) error {
	// 引数のTODOのIDが空の場合＝新規登録の場合、UUIDを生成してタスクのIDとする
	if todo.ID == "" {
		id, err := uuid.NewRandom()
		if err != nil {
			return err
		}
		todo.ID = id.String()
	}
	// 引数のTODOの作成日が空の場合＝新規登録の場合、日付時刻を取得する
	if todo.CreatedAt.IsZero() {
		todo.CreatedAt = time.Now()
	}
	// ロック＞ 登録 ＞アンロック
	db.Lock()
	db.m[todo.ID] = todo
	db.Unlock()

	return nil
}

func (db *MemoryDB) GetAll(ctx context.Context) ([]*TODO, error) {
	// 読み取りロック
	// RLock同士はブロックせずLock（読書ロック）のみブロックされる。
	db.RLock()
	// メソッドの最後にロックを解除する
	defer db.RUnlock()

	// DB登録済のデータ数と同じ長さのTODO配列を生成して、DBの値を追加する
	todos := make([]*TODO, 0, len(db.m))
	for _, todo := range db.m {
		todos = append(todos, todo)
	}
	return todos, nil
}
