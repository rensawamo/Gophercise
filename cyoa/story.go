// cyoa is a package for building Choose Your Own Adventure
// stories that can be rendered via the resulting http.Handler
package cyoa

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"strings"
)

func init() {
	tpl = template.Must(template.New("").Parse(defaultHandlerTmpl))
}

var tpl *template.Template

var defaultHandlerTmpl = `
<!DOCTYPE html>
<html>
  <head>
    <meta charset="utf-8">
    <title>Choose Your Own Adventure</title>
  </head>
  <body>
    <section class="page">
      <h1>{{.Title}}</h1>
      {{range .Paragraphs}}
        <p>{{.}}</p>
      {{end}}
      {{if .Options}}
        <ul>
        {{range .Options}}
          <li><a href="/{{.Chapter}}">{{.Text}}</a></li>
        {{end}}
        </ul>
      {{else}}
        <h3>The End</h3>
      {{end}}
    </section>
    <style>
      body {
        font-family: helvetica, arial;
      }
      h1 {
        text-align:center;
        position:relative;
      }
      .page {
        width: 80%;
        max-width: 500px;
        margin: auto;
        margin-top: 40px;
        margin-bottom: 40px;
        padding: 80px;
        background: #FFFCF6;
        border: 1px solid #eee;
        box-shadow: 0 10px 6px -6px #777;
      }
      ul {
        border-top: 1px dotted #ccc;
        padding: 10px 0 0 0;
        -webkit-padding-start: 0;
      }
      li {
        padding-top: 10px;
      }
      a,
      a:visited {
        text-decoration: none;
        color: #6295b5;
      }
      a:active,
      a:hover {
        color: #7792a2;
      }
      p {
        text-indent: 1em;
      }
    </style>
  </body>
</html>`

// hander のポインタをとることで handler のめそっとの値を変更したりする関数
type HandlerOption func(h *handler)

// ポインタでこのめそっとをいじる関数を定義したいときとかは、 type func 名 を reterun
func WithTemplate(t *template.Template) HandlerOption {
	return func(h *handler) {
		h.t = t
	}
}

// 引数に関数を受けとって シグネイチャーごと 書いてしまう
func WithPathFunc(fn func(r *http.Request) string) HandlerOption {
	return func(h *handler) {
		h.pathFn = fn
	}
}

// ... で Handelerをいじる関数を複数個   引数にうけとることが可能

func NewHandler(s Story, opts ...HandlerOption) http.Handler {
	h := handler{s, tpl, defaultPathFn}
	for _, opt := range opts {
		//  関数を直接いじる関数は ＊
		// だが呼び出し側は  for分などでoperantで値参照もっているから値をかえられるのか
		opt(&h)
	}
	return h
}

type handler struct {
	s      Story
	t      *template.Template           // templeteは ｛｝とかで html のテンプレートをいじる機能を提供
	pathFn func(r *http.Request) string // (r *http.Request) string このシぐネイチャを持った関数を受け取る
}

func defaultPathFn(r *http.Request) string {
	path := strings.TrimSpace(r.URL.Path)
	if path == "" || path == "/" {
		path = "/intro"
	}
	return path[1:]
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := h.pathFn(r)
	fmt.Println("Paht", path)
	fmt.Println("Pah2t", h.s)
	if chapter, ok := h.s[path]; ok {
		err := h.t.Execute(w, chapter) // 解析されたテンプレートを指定されたデータオブジェクトに適応し wrに書き込む
		if err != nil {
			log.Printf("%v", err)
			http.Error(w, "Something went wrong...", http.StatusInternalServerError)
		}
		return
	}
	http.Error(w, "Chapter not found.", http.StatusNotFound) // 無効なリンクという表示
}

func JsonStory(r io.Reader) (Story, error) {
	d := json.NewDecoder(r)
	var story Story
	// 構造対のにオペラントをつけることでメモリの節約
	if err := d.Decode(&story); err != nil {
		return nil, err
	}
	return story, nil
}

// json map decode要因
// map[チャプターの名前] チャプタの内容の構造体という感じ
type Story map[string]Chapter

// 定義  method 型  `json:"key"`
type Chapter struct {
	Title      string   `json:"title"`
	Paragraphs []string `json:"story"`
	Options    []Option `json:"options"`
}

type Option struct {
	Text    string `json:"text"`
	Chapter string `json:"arc"`
}
