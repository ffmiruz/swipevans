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

	products, err := Scrape(link, "div.product-card", ".product-card__title", ".product-card__link@href", ".product-card__price", ".product-card__image-slide@href")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	fmt.Println(len(products))
}

type Item struct {
	Name  string
	Link  string
	Price string
	Image string
}

func Scrape(url, base string, selector ...string) ([]Item, error) {
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("user-agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	doc, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		return nil, err
	}
	item := doc.Find(base)
	collection := make([]Item, 0, len(item.Nodes))
	for i := range item.Nodes {
		var product Item
		for k, v := range selector {
			sel := strings.Split(v, "@")
			result := ""
			// Has attribute
			if len(sel) > 1 {
				result = item.Eq(i).Find(sel[0]).AttrOr(sel[1], "")
			} else {
				result = strings.TrimSpace(item.Eq(i).Find(v).Text())
			}
			switch k {
			case 0:
				product.Name = result
			case 1:
				product.Link = result
			case 2:
				product.Price = result
			case 3:
				product.Image = result
			default:
			}
		}
		collection = append(collection, product)
	}
	return collection, err
}
