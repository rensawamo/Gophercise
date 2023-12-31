package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"runtime/debug"
	"strconv"
	"strings"

	"github.com/alecthomas/chroma/formatters/html"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/debug/", sourceCodeHandler)
	mux.HandleFunc("/panic/", panicDemo)
	mux.HandleFunc("/panic-after/", panicAfterDemo)
	mux.HandleFunc("/", hello)
	log.Fatal(http.ListenAndServe(":3000", devMw(mux)))
}


// パスを受けってコードをブラウザで表示させるその際に chromaで デザインや プログラミング言語を指定して 出力
func sourceCodeHandler(w http.ResponseWriter, r *http.Request) {
	// http://localhost:3000/debug/?line=2009&path=%2Fusr%2Flib%2Fgo-1.21%2Fsrc%2Fnet%2Fhttp%2Fserver.go
	// 上記の   line のとこをとる

	// chomeはデザイン
	path := r.FormValue("path")
	lineStr := r.FormValue("line")
	fmt.Println("linestr",lineStr)
	line, err := strconv.Atoi(lineStr)
	if err != nil {
		line = -1
	}
	file, err := os.Open(path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	b := bytes.NewBuffer(nil)
	// io.Copy は  failをよみこんで b に書き込みをおこなう
	_, err = io.Copy(b, file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	//  ２つの要素を持つ配列の宣言
	var lines [][2]int
	if line > 0 {
		lines = append(lines, [2]int{line, line})
	}
	// lines [[78 78]]
	fmt.Println("lines",lines)
	lexer := lexers.Get("Go") //Go のソースコード 指定 
	iterator, _ := lexer.Tokenise(nil, b.String())
	fmt.Println("iterator",iterator)
	style := styles.Get("github") //github風のテキストで出力
	if style == nil {
		style = styles.Fallback
	}
	// html の contentの設定からいじれる
	formatter := html.New(html.TabWidth(2), html.WithLineNumbers(true), html.LineNumbersInTable(true), html.HighlightLines(lines))
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, "<style>pre { font-size: 4.2em; }</style>")
	formatter.Format(w, style, iterator)
	// _ = quick.Highlight(w, b.String(), "go", "html", "github")
}

func devMw(app http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Println(err)
				stack := debug.Stack()
				log.Println(string(stack))
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintf(w, "<h1>panic: %v</h1><pre>%s</pre>", err, makeLinks(string(stack)))
			}
		}()
		app.ServeHTTP(w, r)
	}
}

func panicDemo(w http.ResponseWriter, r *http.Request) {
	funcThatPanics()
}

func panicAfterDemo(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "<h1>Hello!</h1>")
	funcThatPanics()
}

func funcThatPanics() {
	panic("Oh no!")
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "<h1>Hello!</h1>")
}

func makeLinks(stack string) string {
	lines := strings.Split(stack, "\n")
	for li, line := range lines {
		if len(line) == 0 || line[0] != '\t' {
			continue
		}
		file := ""
		for i, ch := range line {
			if ch == ':' {
				file = line[1:i]
				break
			}
		}
		var lineStr strings.Builder
		for i := len(file) + 2; i < len(line); i++ {
			if line[i] < '0' || line[i] > '9' {
				break
			}
			lineStr.WriteByte(line[i])
		}
		v := url.Values{}
		v.Set("path", file)
		v.Set("line", lineStr.String())
		lines[li] = "\t<a href=\"/debug/?" + v.Encode() + "\">" + file + ":" + lineStr.String() + "</a>" + line[len(file)+2+len(lineStr.String()):]
	}
	return strings.Join(lines, "\n")
}