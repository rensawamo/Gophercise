package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"sync"
	"time"
	"github.com/rensawamo/quiet_hn/hn"
)

func main() {
	// parse flags
	var port, numStories int
	flag.IntVar(&port, "port", 3000, "the port to start the web server on")
	flag.IntVar(&numStories, "num_stories", 30, "the number of top stories to display") // 表示するディスプレイの数の指定
	flag.Parse()

	tpl := template.Must(template.ParseFiles("index.gohtml"))

	http.HandleFunc("/", handler(numStories, tpl))

	// Start the server
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}

func handler(numStories int, tpl *template.Template) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		stories, err := getCachedStories(numStories)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		data := templateData{
			Stories: stories,
			Time:   time.Now().Sub(start),  // subは現在との時間の差分で キャッシュやページの読み込み時間の測定などに活用
		}
		err = tpl.Execute(w, data)
		if err != nil {
			http.Error(w, "Failed to process the template", http.StatusInternalServerError)
			return
		}
	})
}

var (
	cache           []item
	cacheExpiration time.Time 
	cacheMutex      sync.Mutex
)


// キャッシュデータ取得
func getCachedStories(numStories int) ([]item, error) {
	cacheMutex.Lock()
	defer cacheMutex.Unlock()
	if time.Now().Sub(cacheExpiration) < 0 { // 一回でもよばれていたらもとのデータをそのまま返す
		return cache, nil
	}
	stories, err := getTopStories(numStories)
	if err != nil {
		return nil, err
	}
	cache = stories
	cacheExpiration = time.Now().Add(5 * time.Minute)

	return cache, nil
}
 // numStories : 30
func getTopStories(numStories int) ([]item, error) {
	var client hn.Client
	ids, err := client.TopItems()
	if err != nil {
		return nil, fmt.Errorf("Failed to load top stories")
	}
	var stories []item
	at := 0

	// 1開始かまわさない
	for len(stories) < numStories { 
		need := (numStories - len(stories)) * 5 / 4
		// 少し余分に37ことる
		stories = append(stories, getStories(ids[at:at+need])...) // idを指定できる
		at += need
		fmt.Println("at",at)
	}
	// 最後にサイズの調節
	return stories[:numStories], nil
}

func getStories(ids []int) []item {
	type result struct {
		idx  int
		item item
		err  error
	}
	resultCh := make(chan result)
	for i := 0; i < len(ids); i++ {
		go func(idx, id int) {
			var client hn.Client
			hnItem, err := client.GetItem(id)
			if err != nil {
				resultCh <- result{idx: idx, err: err}
			}
			resultCh <- result{idx: idx, item: parseHNItem(hnItem)}
		}(i, ids[i])
	}

	var results []result
	for i := 0; i < len(ids); i++ {
		results = append(results, <-resultCh)
	}
	sort.Slice(results, func(i, j int) bool {
		return results[i].idx < results[j].idx
	})

	var stories []item
	for _, res := range results {
		if res.err != nil {
			continue
		}
		if isStoryLink(res.item) {
			stories = append(stories, res.item)
		}
	}
	return stories
}

func isStoryLink(item item) bool {
	return item.Type == "story" && item.URL != ""
}

func parseHNItem(hnItem hn.Item) item {
	ret := item{Item: hnItem}
	url, err := url.Parse(ret.URL)
	if err == nil {
		ret.Host = strings.TrimPrefix(url.Hostname(), "www.")
	}
	return ret
}

// item is the same as the hn.Item, but adds the Host field
type item struct {
	hn.Item
	Host string
}

type templateData struct {
	Stories []item
	Time    time.Duration
}
