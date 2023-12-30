package main

import (
	"fmt"
	"net/http"

	urlshort "github.com/rensawamo/Gophercise"
)

func main() {
	mux := defaultMux()

	// Build the MapHandler using the mux as the fallback
	//リダイレクト先を指定
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "github.com/rensawamo/Gophercise/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}

	// serveを作成
	mapHandler := urlshort.MapHandler(pathsToUrls, mux)

	// Build the YAMLHandler using the mapHandler as the
	// fallback
	yaml := `
- path: /part2
  url: github.com/rensawamo/Gophercise/urlshort
- path: /urlshort-final
  url: github.com/rensawamo/Gophercise/urlshort/tree/main
`

	// HTTPハンドラ： SeveHttpというmethodを提供しているinterface
	yamlHandler, err := urlshort.YAMLHandler([]byte(yaml), mapHandler)
	if err != nil {
		panic(err)
	}
	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", yamlHandler) // サーバ起動 (TCPアドレス、http.Handler)
}

// HTTPのリクエストとそれに対応するハンドラを登録
func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}
