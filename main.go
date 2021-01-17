package main

import (
	"context"
	"flag"
	"fmt"
	"go-todo-app/internal/todo"
	"os"
	"os/signal"
	"syscall"
)

const (
	ExitOK    = 0
	ExitError = 1
)

func main() {
	//flag with specified name, default value, and usage string.
	flagAddr := flag.String("addr", ":8080", "host:port")
	flag.Parse()
	os.Exit(run(*flagAddr))
}

func run(addr string) int {
	//os.Signal型の容量1のチャンネルを生成
	sigCh := make(chan os.Signal, 1)
	// 引数に指定したSignal発生時にチャネルに通知する。
	// Signalを指定しない場合、全てにシグナルに対して通知を発生させる。
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	errCh := make(chan error)

	db := todo.NewMemoryDB()
	// 使用するポートとDBオブジェクトを引数に与えてサーバーを生成
	// テスト時にはテスト用DBを渡すことができる。
	s := todo.NewServer(addr, db)

	// 別スレッドでサーバー始動。
	go func() {
		errCh <- s.Start()
	}()

	// サーバースレッドでエラーが発生した場合、またはシグナルを検知した場合
	// エラーを返してプログラムを終了する。
	select {
	case err := <-errCh:
		// エラーChに受信し、かつ値がnilでないときにログ出力して終了。
		// errChにnil返すときってどういうとき？？
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return ExitError
		}
	case <-sigCh:
		if err := s.Stop(context.Background()); err != nil {
			return ExitError
		}
	}
	return ExitOK
}
