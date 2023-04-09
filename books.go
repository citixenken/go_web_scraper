package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"

	"github.com/gocolly/colly"
)

type Book struct {
	Title string
	Price string
}

func main() {
	// scraping done here

	// save scraped data to CSV file
	file, err := os.Create("books.csv")
	if err != nil {
		log.Fatal(err)
	}

	// delay closing the file until program completes its cycle
	defer file.Close()

	// create a CSV writer
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// write headers
	headers := []string{"Title", "Price"}
	writer.Write(headers)


	// The Collector makes HTTP requests and traverses HTML pages.
	c := colly.NewCollector(
		colly.AllowedDomains("books.toscrape.com"),
	)

	// var books []book
	// execute when matching selectors are found
	c.OnHTML("title", func(h *colly.HTMLElement){
		fmt.Println(h.Text)
	})

	// extract book titles and prices
	c.OnHTML(".product_pod", func(h *colly.HTMLElement){
		// book := book {
		// 	Title: h.ChildAttr(".image_container img", "alt"),
		// 	Price:  h.ChildText(".price_color"),
		// }

		// books = append(books, book)

		// write each book as a single row
		book := Book{}
		book.Title = h.ChildAttr(".image_container img", "alt")
		book.Price = h.ChildText(".price_color")
		row := []string{book.Title, book.Price}
		writer.Write(row)
		// fmt.Println(book.Title, book.Price)
	})

	// examine the response
	c.OnResponse(func(r *colly.Response){
		fmt.Println(r.StatusCode)
	})

	// track which URL is being visited
	c.OnRequest(func(r *colly.Request){
		fmt.Println("Visiting", r.URL.String())
	})

	// handling pagination, then crawl converted URL
	/*
	The existing function that scrapes the book information
	will be called on all of the resulting pages as well.
	No additional code is needed.
	*/
	c.OnHTML(".next > a", func(h *colly.HTMLElement){
		next_page := h.Request.AbsoluteURL(h.Attr("href"))
		c.Visit(next_page)
	})

	// start the scraper
	c.Visit("https://books.toscrape.com/")


	fmt.Println("Done!")
}

