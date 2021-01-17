package todo

import (
	"context"
	"encoding/json"
	"net/http"
)

// httpサーバーとDBを要素にもつstruct
type Server struct {
	server *http.Server
	db     DB
}

// サーバー生成用関数
// アドレスとDBを受け取る
func NewServer(addr string, db DB) *Server {
	return &Server{
		server: &http.Server{Addr: addr},
		db:     db,
	}
}

// サーバー始動
func (s *Server) Start() error {

	// Handler設定
	s.initHandlers()
	// 始動
	err := s.server.ListenAndServe()
	// ErrServerClosed以外のエラーが帰ってきた場合errを返す
	// 呼び出し側(main)でログ出力してサーバー終了する
	if err != nil && err != http.ErrServerClosed {
		return err
	}
	// それ以外の場合nil（正常終了）
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	// 現在アクティブな接続が完了するまで待機してからサーバーを終了する。
	// 待機中にcontextがタイムアウトした場合、contextエラーを返す。
	// 正常にサーバーをShutdownできた場合、ErrServerClosedを返す。
	err := s.server.Shutdown(ctx)
	if err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

// ここから ＝＝＝＝＝＝＝＝＝＝＝＝＝＝＝＝＝＝＝＝＝＝＝＝＝＝＝＝＝＝＝＝＝＝＝＝＝＝＝＝＝＝＝＝＝＝＝＝＝＝＝＝＝＝＝＝＝＝＝＝＝＝＝＝＝＝＝＝＝＝＝＝＝＝＝＝＝＝＝＝＝＝
func (s *Server) initHandlers() {
	mux := http.NewServeMux()
	s.server.Handler = mux

	// httpサーバにRouter設定
	mux.HandleFunc("/create", s.HandleCreate)
	mux.HandleFunc("/getall", s.HandleGetAll)
}

// TODO作成ハンドラ
func (s *Server) HandleCreate(w http.ResponseWriter, r *http.Request) {
	var todo TODO
	// BodyのJSONをデコードしてTODO structの変数に格納する。
	// エラーだった場合はBadRequestレスポンスを返す
	if err := json.NewDecoder(r.Body).Decode(&todo); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// デコードがエラー出ない場合、DB#Putにコンテキストとデコードした値を渡してデータを登録する
	if err := s.db.Put(r.Context(), &todo); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// TODO取得ハンドラ
func (s *Server) HandleGetAll(w http.ResponseWriter, r *http.Request) {
	// DB＃GetALLで全データ取得してTODO配列に格納
	// エラー担った場合はInternalErrorを返す
	todos, err := s.db.GetAll(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// データを取得できた場合、JSONにエンコードしてResponseWriterに書き込む
	if err := json.NewEncoder(w).Encode(todos); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
