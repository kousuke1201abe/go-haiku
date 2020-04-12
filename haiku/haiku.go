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
	args := make(map[string]string)
	args = map[string]string{"dicdir": "/usr/lib/x86_64-linux-gnu/mecab/dic/mecab-ipadic-neologd"}
	mecab, err := mecab.New(args)
	if err != nil {
		panic(err)
	}
	defer mecab.Destroy()

	node, err := mecab.ParseToNode(text)

	var words Words
	for ; !node.IsZero(); node = node.Next() {
		feature_arr := strings.Split(node.Feature(), ",")
		yomi := feature_arr[len(feature_arr)-1]
		replacer := strings.NewReplacer("ゃ", "", "ゅ", "", "ょ", "", "ぁ", "", "ぃ", "", "ぅ", "", "ぇ", "", "ぉ", "", "ゎ", "", "ャ", "", "ュ", "", "ョ", "", "ァ", "", "ィ", "", "ゥ", "", "ェ", "", "ォ", "", "ヮ", "")
		if feature_arr[0] != "記号" && feature_arr[len(feature_arr)-1] != "*" {
			feature_arr[len(feature_arr)-1] = replacer.Replace(yomi)
			var word Word
			word.count = len([]rune(yomi))
			word.pos = feature_arr[0]
			words = append(words, word)
		}
	}

	return words
}

func (words Words) CheckHaiku() bool {
	result := false
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
	count := 0

	for _, word := range words {
		count += word.count
	}

	return count
}

func (words Words) checkKu(startIndex int, expectedCount int) (bool, int) {
	count := 0
	endIndex := 0
	result := false

	if words[startIndex].pos != "名詞" && words[startIndex].pos != "形容詞" && words[startIndex].pos != "形容動詞" && words[startIndex].pos != "副詞" && words[startIndex].pos != "連体詞" && words[startIndex].pos != "感動詞" || words[startIndex].pos == "接頭詞" {
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
