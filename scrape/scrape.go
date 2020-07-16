package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	json "github.com/json-iterator/go"
)

func main() {
	file := "scrape/seed.txt"
	fd, err := os.Open(file)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer fd.Close()

	scanner := bufio.NewScanner(fd)
	for line := 0; scanner.Scan(); line++ {
		row := scanner.Text()
		col := strings.Split(row, ",")
		// Skip incomplete row/Not enough parameters
		if len(col) < 3 {
			continue
		}

		products, err := Scrape(col[0], col[1], col[2], col[3], col[4], col[5])
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}
		if len(products) < 1 {
			fmt.Fprintln(os.Stdout, "Got no results from", col[0])
			continue
		}

		out, err := os.Create("./data/" + strconv.Itoa(line) + ".json")
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}
		defer out.Close()

		data, err := json.Marshal(map[string]interface{}{
			"from":  baseURL(col[0]),
			"count": len(products),
			"items": products,
		})
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}
		_, err = out.Write(data)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}
	}

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
				result = item.Eq(i).Find(sel[0]).First().AttrOr(sel[1], "")
			} else {
				result = strings.TrimSpace(item.Eq(i).Find(v).First().Text())
			}
			switch k { // Find image first
			case 0:
				product.Name = result
			case 1:
				product.Link = result
			case 2:
				product.Price = result
			case 3:
				// Ditch the product if has no image
				if result == "" {
					break
				}
				product.Image = result
			default:
			}
		}

		if product.Image == "" {
			continue
		}
		collection = append(collection, product)
	}
	return collection, err
}

func baseURL(url string) string {
	base := strings.Split(url, "/")[2]
	return base
}
