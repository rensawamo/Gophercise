package link

// dom 操作 golnag ver
import (
	"io" // go のos のような感じ
	"strings"

	"golang.org/x/net/html"
)

// パースしたい要素の構造体をつくる
type Link struct {
	Href string
	Text string
}

// Parse parses the HTML file and returns Links
func Parse(r io.Reader) ([]Link, error) {
	doc, err := html.Parse(r)
	if err != nil {
		return nil, err
	}
	return getAllLinks(doc), nil

}

// html の解析には、この nodeが定石
func getAllLinks(n *html.Node) []Link {
	var links []Link

	// ほしい要素をここで指定（型をきめてしまう） if文で <a> タグだけに反応して
	if n.Type == html.ElementNode && n.Data == "a" { // 要素ノードかどうかの判定  プラス  .Data で html のタグを探せる →  <a> tagの要素をさがす
		for _, a := range n.Attr {
			if a.Key == "href" {
				txt := extractText(n)
				links = append(links, Link{a.Val, txt})
			}
		}
		// これで 次のノードに行く 再起的に a タグの精製を繰り返す
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		exLinks := getAllLinks(c)
		println(exLinks)
		links = append(links, exLinks...)
	}
	return links
}

func extractText(n *html.Node) string {
	var text string
	if n.Type != html.ElementNode && n.Data != "a" && n.Type != html.CommentNode {
		text = n.Data
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		text += extractText(c)
	}
	return strings.Trim(text, "\n")
}
