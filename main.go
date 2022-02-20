package main

import (
	"encoding/csv"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"sort"
	"strconv"
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
	checkError("Cannot get url content", err)

	body, err := ioutil.ReadAll(resp.Body)

	defer resp.Body.Close()
	return string(body)
}

func cleanText(text string) string {
	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	checkError("Cannot clean text", err)

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

func saveCsvResults(words map[string]int) {
	file, err := os.Create("./word_frequencies_report.csv")
	checkError("Cannot create file", err)
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for _, word := range sortedWords(words) {
		row := []string{word.key, strconv.Itoa(word.value)}

		err := writer.Write(row)
		checkError("Cannot write to file", err)
	}

}

func checkError(message string, err error) {
	if err != nil {
		log.Fatal(message, err)
	}
}

func main() {
	content := getText(os.Args[1])
	words := wordCount(content)
	saveCsvResults(words)
}
