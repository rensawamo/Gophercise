package urlshort

import (
	"fmt"
	"net/http"

	yaml "gopkg.in/yaml.v2"
)

// MapHandlerはhttp.HandlerFuncを返す。
// // を実装している) を返す。
// // パス (マップのキー) を対応する URL (マップの各キーが指す値、文字列形式) に マップしようとする。
// (マップの各キーが指す値、文字列形式) へのマッピングを試みる。
// マップ内でパスが提供されていない場合、フォールバックとして // http.Handlerが使用。
// // http.Handler が代わりに呼び出される。

// ★ http.HandlerFunc :  Webサーバの一部として機能してHTTPリクエストなどをこなす
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if dest, ok := pathsToUrls[path]; ok {
			fmt.Printf("move")
			http.Redirect(w, r, dest, http.StatusFound) // 存在するパスなら正規のurl に移動
			return
		}
		// パスが 存在しないともとの場所に戻す
		fmt.Println("fallback:", fallback)
		fallback.ServeHTTP(w, r)
	}
}

// YAMLHandler は指定された YAML をパースして
// を返す。
// を返す。
// URL にマップする。パスがYAMLで提供されない場合
// フォールバックのhttp.Handlerが代わりに呼び出される。

// YAMLは次のフォーマットであることが期待される：

//   - パス： パス: /some-path
//     url: https://www.some-url.com/demo
//
// を経由して同様の http.HandlerFunc を作成するには MapHandler を参照。
// を使います。

// yaml
func YAMLHandler(yamlBytes []byte, fallback http.Handler) (http.HandlerFunc, error) {
	pathUrls, err := parseYaml(yamlBytes)
	if err != nil {
		return nil, err
	}
	pathsToUrls := buildMap(pathUrls)
	return MapHandler(pathsToUrls, fallback), nil
}

func buildMap(pathUrls []pathUrl) map[string]string {
	pathsToUrls := make(map[string]string)
	for _, pu := range pathUrls {
		pathsToUrls[pu.Path] = pu.URL
	}
	fmt.Println("rowMaps:", pathsToUrls)
	return pathsToUrls
}

func parseYaml(data []byte) ([]pathUrl, error) {
	var pathUrls []pathUrl
	err := yaml.Unmarshal(data, &pathUrls)
	fmt.Println(pathUrls)
	// [{/part2 github.com/rensawamo/Gophercise/urlshort} {/urlshort-final github.com/rensawamo/Gophercise/urlshort/tree/main}] 写像
	if err != nil {
		return nil, err
	}
	return pathUrls, nil
}

type pathUrl struct {
	Path string `yaml:"path"` // Unmarshal してるときは構造対に型指定がはいる
	URL  string `yaml:"url"`
}
