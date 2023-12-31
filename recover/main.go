package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"net/http"
	"runtime/debug"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/panic/", panicDemo)
	mux.HandleFunc("/panic-after/", panicAfterDemo)
	mux.HandleFunc("/", hello)
	// recoverMnのログへの標準出力
	// /panic/ のパスの指定で  わざとpanicをおこす
	log.Fatal(http.ListenAndServe(":3002", recoverMw(mux, true)))
}

func recoverMw(app http.Handler, dev bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Println(err)
				stack := debug.Stack()
				log.Println(string(stack))
				if !dev { // 開発の準備などにぱにくった時の対策
					http.Error(w, "Something went wrong :(", http.StatusInternalServerError)
					return
				}
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintf(w, "<h1>panic: %v</h1><pre>%s</pre>", err, string(stack)  )
			}
		}()
		// 標準の http.ResponseWriter をラップして拡張するカスタム responseWriter を使用して、
		// HTTPレスポンスをカスタマイズするためのもの
		// です。このカスタム responseWriter によって、標準のレスポンスライターにはない追加の機能や情報を提供
		nw := &responseWriter{ResponseWriter: w}
		app.ServeHTTP(nw, r)
		// panic のとき にここは実行されない

		nw.flush()
	}
}

// type ResponseWriter interface {
// 	Header() Header
// 	Write([]byte) (int, error)
// 	WriteHeader(statusCode int)
// }

type responseWriter struct {
	http.ResponseWriter
	writes [][]byte // スライスのスライス 複数対応
	status int
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	rw.writes = append(rw.writes, b)
	return len(b), nil
}

 // status コードを変える カスタムインターフェイス
func (rw *responseWriter) WriteHeader(statusCode int) {
	rw.status = statusCode
}

func (rw *responseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	hijacker, ok := rw.ResponseWriter.(http.Hijacker)
	if !ok {
		return nil, nil, fmt.Errorf("the ResponseWriter does not support the Hijacker interface")
	}
	return hijacker.Hijack()
}

func (rw *responseWriter) Flush() {
	flusher, ok := rw.ResponseWriter.(http.Flusher)
	if !ok {
		return
	}
	flusher.Flush()
}

func (rw *responseWriter) flush() error {
	if rw.status != 0 {
		rw.ResponseWriter.WriteHeader(rw.status)  // パニックからの回復のためのステータスコードのかきかえ
	}
	for _, write := range rw.writes {
		_, err := rw.ResponseWriter.Write(write)
		if err != nil {
			return err
		}
	}
	return nil
}

func panicDemo(w http.ResponseWriter, r *http.Request) {
	funcThatPanics()
}

func panicAfterDemo(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "<h1>Hellgredfgo!</h1>")
	funcThatPanics()
}

func funcThatPanics() {
	panic("Oh no!")
}

func hello(w http.ResponseWriter, r *http.Request) {
	// オーバライドwite が 	n, err = w.Write(p.buf)で実行されて  
	// &つけたresponseWriterの writeがかきかえられる
	fmt.Fprintln(w, "<h1>fggfdgHello!</h1>")
}
