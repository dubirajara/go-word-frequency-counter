package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"sort"
	"strings"
)

type word struct {
	Key   string
	Value int
}

func wordCount(str string) map[string]int {
	stpwords, _ := openStopWords("./StopWords.txt")
	wordList := strings.Fields(str)
	wordCounts := make(map[string]int)
	for _, word := range wordList {
		cleanWord := cleanText(strings.ToLower(word))
		found := checkStopWords(stpwords, cleanWord)
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

func openStopWords(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func checkStopWords(stpwords []string, word string) bool {
	for _, item := range stpwords {
		if item == word {
			return true
		}
	}
	return false
}

func sortedWords(words map[string]int) []word {
	var sorted []word
	for k, v := range words {
		sorted = append(sorted, word{k, v})
	}

	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Value > sorted[j].Value
	})
	return sorted
}

func main() {
	content := getText("https://storage.googleapis.com/apache-beam-samples/shakespeare/romeoandjuliet.txt")
	words := wordCount(content)
	for _, word := range sortedWords(words) {
		fmt.Println(word.Key, word.Value)
	}
}
