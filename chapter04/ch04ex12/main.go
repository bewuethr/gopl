// ch04ex12 builds an index of xkcd transcripts that can be queried.
package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"
)

type comic struct {
	Month      string
	Num        int
	Link       string
	Year       string
	News       string
	SafeTitle  string `json:"safe_title"`
	Transcript string
	Alt        string
	Img        string
	Title      string
	Day        string
}

var comics []comic
var index = make(map[string][]int)

// Build index: read files into array of JSON objects; while reading, build a
// map[string][]int for the words in title and transcript, mapping to the array
// indices for comics containing the words
func init() {
	os.Chdir("infofiles")
	files, err := ioutil.ReadDir(".")
	if err != nil {
		log.Fatal(err)
	}

	comics = make([]comic, len(files), len(files))

	for idx, file := range files {
		jsonData, err := ioutil.ReadFile(file.Name())
		if err != nil {
			log.Fatal(err)
		}
		json.Unmarshal(jsonData, &comics[idx])

		// Split title and transcript into lowercase words after removing all
		// interpunction, deduplicate, then build index
		re := regexp.MustCompile("[^[:alnum:][:blank:]]")
		noPunctStr := re.ReplaceAllString(comics[idx].Title+" "+comics[idx].Transcript, "")
		words := dedupe(strings.Fields(strings.ToLower(noPunctStr)))
		for _, word := range words {
			index[word] = append(index[word], idx)
		}
	}
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("xkcd>> ")
		scanned := scanner.Scan()
		if !scanned {
			os.Exit(0)
		}
		input := scanner.Text()
		for i, val := range index[input] {
			printComic(val)
			if i < len(index[input])-1 {
				fmt.Printf("\n====\n\n")
			}
		}
	}
}

// dedupe removes duplicates from a string slice.
func dedupe(a []string) []string {
	seen := make(map[string]bool)
	d := []string{}
	for _, val := range a {
		if !seen[val] {
			d = append(d, val)
			seen[val] = true
		}
	}
	return d
}

func printComic(idx int) {
	fmt.Printf("https://xkcd.com/%d\n----\n%s\n\n%s\n",
		comics[idx].Num, comics[idx].Title, comics[idx].Transcript)
}
