package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/http"
	"os"
	"strings"
)

// Quote represents the structure of a quote
type Quote struct {
	Text   string
	Author string
	Tags   []string
}

// loadConfig reads configuration from config.json
func loadConfig() (string, error) {
	file, err := os.Open("config.json")
	if err != nil {
		return "", err
	}
	defer file.Close()

	config := struct {
		DataSource string `json:"data_source"`
	}{}

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return "", err
	}

	return config.DataSource, nil
}

// GetQuotes fetches quotes from the starting URL and handles pagination
func GetQuotes(startURL string) ([]Quote, error) {
	var allQuotes []Quote
	nextPageURL := startURL

	for nextPageURL != "" {
		quotes, nextPage, err := GetQuotesFromPage(nextPageURL)
		if err != nil {
			return nil, err
		}
		allQuotes = append(allQuotes, quotes...)
		nextPageURL = nextPage
	}

	return allQuotes, nil
}

// GetQuotesFromPage extracts quotes from a single page
func GetQuotesFromPage(url string) ([]Quote, string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, "", err
	}

	var quotes []Quote

	doc.Find(".quote").Each(func(index int, item *goquery.Selection) {
		text := item.Find(".text").Text()
		author := item.Find(".author").Text()

		var tags []string
		item.Find(".tags .tag").Each(func(index int, tagItem *goquery.Selection) {
			tags = append(tags, tagItem.Text())
		})

		quotes = append(quotes, Quote{
			Text:   text,
			Author: author,
			Tags:   tags,
		})
	})

	nextPage := ""
	doc.Find(".pager .next a").Each(func(index int, item *goquery.Selection) {
		href, exists := item.Attr("href")
		if exists {
			nextPageURL := fmt.Sprintf("https://quotes.toscrape.com%s", href)
			nextPage = nextPageURL
		}
	})

	return quotes, nextPage, nil
}

func main() {
	dsn, err := loadConfig()
	if err != nil {
		log.Fatal("Error loading configuration:", err)
	}

	fmt.Println("Connecting to database with DSN:", dsn)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Error opening database connection:", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal("Error pinging database:", err)
	}

	allQuotes, err := GetQuotes("https://quotes.toscrape.com/")
	if err != nil {
		log.Println("Error getting quotes:", err)
		return
	}

	for _, q := range allQuotes {
		_, err := db.Exec("INSERT INTO quotes (text, author, tags) VALUES (?, ?, ?)",
			q.Text, q.Author, strings.Join(q.Tags, ","))
		if err != nil {
			log.Println("Failed to insert quote:", err)
		}
	}

	fmt.Println("Data inserted successfully")
}
