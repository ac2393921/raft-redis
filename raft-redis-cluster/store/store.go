package store

import (
	"context"
	"io"

	"github.com/bootjp/go-kvlib/store"
)

type Txn = store.Txn

// Txnは、トランザクション用のインターフェースを定義する
// Storeの実装を自身で用意する場合は、store.Txnと同等のインターフェースを定義しこちらは削除する
// type Txn interface {
// 	Get(ctx context.Context, key []byte) ([]byte, error)
// 	Put(ctx context.Context, key []byte, value []byte) error
// 	Delete(ctx context.Context, key []byte) error
// 	Exists(ctx context.Context, key []byte) (bool, error)
// }

// Storeは、キーバリューストアのインターフェースを定義する
// Get, Put, Delete, Exists, Snapshot, Restore, Txn, Closeの関数を提供する
// このインターフェースを実装することで、任意のキーバリューストアを利用できる
type Store interface {
	Get(ctx context.Context, key []byte) ([]byte, error)
	Put(ctx context.Context, key []byte, value []byte) error
	Delete(ctx context.Context, key []byte) error
	Exists(ctx context.Context, key []byte) (bool, error)
	Snapshot() (io.ReadWriter, error)
	Restore(buf io.Reader) error
	// Txn トランザクション用の関数を提供する
	// トランザクションないで複数の操作をまとめて実行するために使用する
	// トランザクション内でエラーが発生した場合、トランザクションはコミットされる
	// トランザクション内でエラーが発生しなかった場合、トランザクションはロールバックされる
	// トランザクション内で発生したエラーは呼び出し元に返される
	Txn(ctx context.Context, f func(ctx context.Context, txn Txn) error) error
	Close() error
}

var ErrKeyNotFound = store.ErrKeyNotFound

func NewMemoryStore() Store {
	return store.NewMemoryStore()
}
