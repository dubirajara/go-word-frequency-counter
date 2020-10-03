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

func wordCount(str string) map[string]int {
	stpwords, _ := openStopWordsFile("./StopWords.txt")
	wordList := strings.Fields(str)
	wordcounts := make(map[string]int)
	for _, word := range wordList {
		cleanWord := cleanText(strings.ToLower(word))
		found := checkStopWords(stpwords, cleanWord)
		if !found {
			wordcounts[cleanWord]++
		}
	}
	return wordcounts
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

func openStopWordsFile(path string) ([]string, error) {
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

func main() {
	content := getText("https://storage.googleapis.com/apache-beam-samples/shakespeare/romeoandjuliet.txt")
	words := wordCount(content)
	type kv struct {
		Key   string
		Value int
	}
	var sorted []kv
	for k, v := range words {
		sorted = append(sorted, kv{k, v})
	}

	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Value > sorted[j].Value
	})

	for _, kv := range sorted {
		fmt.Println(kv.Key, kv.Value)
	}

}
