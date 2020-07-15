package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func main() {
	file := "./seed.txt"
	fd, err := os.Open(file)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer fd.Close()

	scanner := bufio.NewScanner(fd)
	for scanner.Scan() {
		break
		fmt.Println(string(scanner.Bytes()))
	}

	link := "https://kith.com/collections/vans-collection?sort_by=created-ascending"

	Scrape(link, "div.product-card", ".product-card__title", ".product-card__link@href", ".product-card__price", ".product-card__image-slide@href")
}

func Scrape(url, base string, selector ...string) {
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("user-agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	doc, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	item := doc.Find(base)
	for i := range item.Nodes {
		for _, v := range selector {
			sel := strings.Split(v, "@")
			// Has attribute
			if len(sel) > 1 {
				pr := item.Eq(i).Find(sel[0]).AttrOr(sel[1], "")
				fmt.Println(pr)
				continue
			}
			pr := item.Eq(i).Find(v).Text()
			fmt.Println(strings.TrimSpace(pr))
		}
	}
}
