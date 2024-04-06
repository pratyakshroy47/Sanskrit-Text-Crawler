package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
)

func extractSanskritTexts(url string, wg *sync.WaitGroup, ch chan<- []string) {
	defer wg.Done()

	// Send HTTP request to the URL
	response, err := http.Get(url)
	if err != nil {
		fmt.Println("Error fetching URL:", err)
		return
	}
	defer response.Body.Close()

	// Read response body
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	// Parse HTML content
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(body)))
	if err != nil {
		fmt.Println("Error parsing HTML:", err)
		return
	}

	// Extract Sanskrit texts from HTML elements
	var sanskritTexts []string
	doc.Find(".vedic").Each(func(i int, s *goquery.Selection) {
		sanskritTexts = append(sanskritTexts, s.Text())
	})

	ch <- sanskritTexts
}

func main() {
	// List of URLs to scrape concurrently
	urls := []string{
		"https://sanskritdocuments.org/doc_veda/goShThasUkta.html",
		"https://sanskritdocuments.org/doc_upanishhat/atharvashikha.html",
		// Add more URLs here...
	}

	// Channel to receive scraped data
	ch := make(chan []string)

	// WaitGroup to wait for all goroutines to finish
	var wg sync.WaitGroup

	// Start a goroutine for each URL
	for _, url := range urls {
		wg.Add(1)
		go extractSanskritTexts(url, &wg, ch)
	}

	// Close the channel once all goroutines are done
	go func() {
		wg.Wait()
		close(ch)
	}()

	// Create or open the output file
	file, err := os.Create("sanskrit_texts.txt")
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	// Collect scraped data from the channel and write to file
	for sanskritTexts := range ch {
		// Write scraped data to file
		for _, text := range sanskritTexts {
			_, err := file.WriteString(text + "\n")
			if err != nil {
				fmt.Println("Error writing to file:", err)
				return
			}
		}
	}
}
