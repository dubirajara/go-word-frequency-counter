package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"sort"
	"strings"

	"stopwords/stopwords"
)

//Word type struct defined.
type Word struct {
	key   string
	value int
}

func wordCount(str string) map[string]int {
	wordList := strings.Fields(str)
	wordCounts := make(map[string]int)
	for _, word := range wordList {
		cleanWord := cleanText(strings.ToLower(word))
		found := checkStopWords(cleanWord)
		if !found {
			wordCounts[cleanWord]++
		}
	}
	return wordCounts
}

func getText(url string) string {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalln(err)
	}
	body, err := ioutil.ReadAll(resp.Body)

	defer resp.Body.Close()
	return string(body)
}

func cleanText(text string) string {
	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		log.Fatal(err)
	}
	proccessedWord := reg.ReplaceAllString(text, "")
	return proccessedWord
}

func checkStopWords(word string) bool {
	for _, item := range stopwords.StopWords() {
		if item == word {
			return true
		}
	}
	return false
}

func sortedWords(words map[string]int) []Word {
	var sorted []Word
	for k, v := range words {
		sorted = append(sorted, Word{k, v})
	}

	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].value > sorted[j].value
	})
	return sorted
}

func main() {
	content := getText(os.Args[1])
	words := wordCount(content)
	for _, word := range sortedWords(words) {
		fmt.Println(word.key, word.value)
	}
}
