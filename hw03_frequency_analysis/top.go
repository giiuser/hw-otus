package hw03frequencyanalysis

import (
	"sort"
	"strings"
)

type WordCount struct {
	word  string
	count int
}

func Top10(text string) []string {
	wordsArr := strings.Fields(text)

	wc := make(map[string]int)
	for _, word := range wordsArr {
		_, matched := wc[word]
		if matched {
			wc[word]++
			continue
		}
		wc[word] = 1
	}

	wordCounts := make([]WordCount, 0, len(wc))
	for key, val := range wc {
		wordCounts = append(wordCounts, WordCount{word: key, count: val})
	}

	sort.Slice(wordCounts, func(i, j int) bool {
		if wordCounts[i].count == wordCounts[j].count {
			return wordCounts[i].word < wordCounts[j].word
		}
		return wordCounts[i].count > wordCounts[j].count
	})

	var finalCount int
	if len(wordCounts) > 10 {
		finalCount = 10
	} else {
		finalCount = len(wordCounts)
	}

	res := make([]string, 0, finalCount)
	for i := 0; i < finalCount; i++ {
		res = append(res, wordCounts[i].word)
	}

	return res
}
