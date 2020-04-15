package haiku

import (
	"strings"

	mecab "github.com/shogo82148/go-mecab"
)

type Word struct {
	count int
	pos   string
}

type Words []Word

func NewWords(text string) Words {
	var words Words

	args := make(map[string]string)
	args = map[string]string{"dicdir": "/usr/lib/x86_64-linux-gnu/mecab/dic/mecab-ipadic-neologd"}
	mecab, err := mecab.New(args)
	if err != nil {
		panic(err)
	}
	defer mecab.Destroy()

	node, err := mecab.ParseToNode(text)

	for ; !node.IsZero(); node = node.Next() {
		feature_arr := strings.Split(node.Feature(), ",")
		yomi := feature_arr[len(feature_arr)-1]
		replacer := strings.NewReplacer("ゃ", "", "ゅ", "", "ょ", "", "ぁ", "", "ぃ", "", "ぅ", "", "ぇ", "", "ぉ", "", "ゎ", "", "ャ", "", "ュ", "", "ョ", "", "ァ", "", "ィ", "", "ゥ", "", "ェ", "", "ォ", "", "ヮ", "")
		if feature_arr[0] != "記号" && feature_arr[len(feature_arr)-1] != "*" {
			var word Word
			word.count = len([]rune(replacer.Replace(yomi)))
			word.pos = feature_arr[0]
			words = append(words, word)
		}
	}

	return words
}

func (words Words) CheckHaiku() bool {
	var result bool
	var endIndex int

	if words.count() == 17 {
		result, endIndex = words.checkKu(0, 5)
	}
	if result == true {
		result, endIndex = words.checkKu(endIndex+1, 7)
	}
	if result == true {
		result, endIndex = words.checkKu(endIndex+1, 5)
	}

	return result
}

func (words Words) count() int {
	var count int

	for _, word := range words {
		count += word.count
	}

	return count
}

func (words Words) checkKu(startIndex int, expectedCount int) (bool, int) {
	var count, endIndex int
	var result bool

	pos := words[startIndex].pos

	if pos != "動詞" && pos != "名詞" && pos != "形容詞" && pos != "形容動詞" && pos != "副詞" && pos != "連体詞" && pos != "感動詞" || pos == "接頭詞" {
		return result, startIndex
	}

	for i := startIndex; i < len(words)+1; i++ {
		if count == expectedCount {
			result = true
			break
		}
		count += words[i].count
		endIndex = i
	}

	return result, endIndex
}
