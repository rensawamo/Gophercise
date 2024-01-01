//go:generate stringer -type=Suit,Rank

package deck

import (
	"fmt"
	"math/rand"
	"sort"
	"time"
)

type Suit uint8 // unit8は メモリ節約

const (
	Spade Suit = iota // iotaでも 初期値の0を回避しないと 0が入る
	Diamond
	Club
	Heart
	Joker // this is a special case
)

var suits = [...]Suit{Spade, Diamond, Club, Heart} // [...]で初期値によって配列のサイズが決定

type Rank uint8

const (
	_ Rank = iota // 定数内自動インクリメント 0を無視
	// この定義の仕方で const Two Rank = 2 の値が割り振られる  = for 分も回せる
	Ace
	Two
	Three
	Four
	Five
	Six
	Seven
	Eight
	Nine
	Ten
	Jack
	Queen
	King
)

const (
	minRank = Ace // 1
	maxRank = King
)

type Card struct {
	Suit // 柄
	Rank // ランク
}

func (c Card) String() string {
	if c.Suit == Joker {
		return c.Suit.String()
	}
	return fmt.Sprintf("%s of %ss", c.Rank.String(), c.Suit.String())
}

// 可変長引数は0も値も許容する
// 引数をとらないときは 下記の opts のfor が回らない
func New(opts ...func([]Card) []Card) []Card {
	var cards []Card
	for _, suit := range suits {
		for rank := minRank; rank <= maxRank; rank++ {
			fmt.Println("rank", rank)
			cards = append(cards, Card{Suit: suit, Rank: rank})
		}
	}
	fmt.Println("cards", cards)
	// 可変長が 0のときは動かない

	// fiter などの特別の設定をくわえたいときなどに こちらの optでが動くｙほうにしたいから funcを可変長引数でうけとっている
	for _, opt := range opts {
		cards = opt(cards)
	}
	return cards
}

func DefaultSort(cards []Card) []Card {
	fmt.Println("引数何", cards)
	// sort.Slice(ソート対象、how)
	sort.Slice(cards, Less(cards)) // less func(i int, j int) bool が必要
	return cards
}

func Sort(less func(cards []Card) func(i, j int) bool) func([]Card) []Card {
	return func(cards []Card) []Card {
		sort.Slice(cards, less(cards))
		return cards
	}
}

// 柄番号とランクの掛け算で判定
func Less(cards []Card) func(i, j int) bool {
	return func(i, j int) bool { // シグネイチャを強制的につくらされたけど
		return absRank(cards[i]) < absRank(cards[j]) // 別に返り値に依存できる
	}
}

func absRank(c Card) int {
	return int(c.Suit)*int(maxRank) + int(c.Rank)
}

// 毎回違う数列をつくるために 時間という概念をつかって作成
var shuffleRand = rand.New(rand.NewSource(time.Now().Unix()))

func Shuffle(cards []Card) []Card {
	ret := make([]Card, len(cards))

	perm := shuffleRand.Perm(len(cards))
	// 毎回len がおなじの内容が異なる数列randam で作成するイメージ
	fmt.Println("perm", perm)
	for i, j := range perm {
		ret[i] = cards[j]
	}
	return ret
}

func Jokers(n int) func([]Card) []Card {
	return func(cards []Card) []Card {
		for i := 0; i < n; i++ {
			cards = append(cards, Card{
				Rank: Rank(i),
				Suit: Joker,
			})
		}
		fmt.Println("クロージャの引数は var になりうるか", cards)
		return cards
	}
}

func Filter(f func(card Card) bool) func([]Card) []Card {
	return func(cards []Card) []Card {
		var ret []Card
		for _, c := range cards {
			if !f(c) {
				ret = append(ret, c)
			}
		}
		return ret
	}
}

func Deck(n int) func([]Card) []Card {
	return func(cards []Card) []Card {
		var ret []Card
		for i := 0; i < n; i++ {
			ret = append(ret, cards...)
		}
		return ret
	}
}
