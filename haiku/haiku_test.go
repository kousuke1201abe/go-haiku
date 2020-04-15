package haiku

import (
	"testing"
)

func TestCheckHaiku(t *testing.T) {
	valid_haiku := []string{
		"古池や蛙飛びこむ水の音",
		"松たか子藤原竜也志村けん",
		"図書館は今日はお休み休館日",
		"今週も、一所懸命頑張ろう",
	}

	for _, haiku := range valid_haiku {
		words := NewWords(haiku)
		result := words.CheckHaiku()
		if result != true {
			t.Fatal("failed test")
		}
	}

	invalid_haiku := []string{
		"これは俳句じゃないよ",
		"あああああああああああああああああ",
		"ProjectManagementProject",
	}

	for _, haiku := range invalid_haiku {
		words := NewWords(haiku)
		result := words.CheckHaiku()
		if result == true {
			t.Fatal("failed test")
		}
	}
}
