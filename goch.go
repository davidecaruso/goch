package main

import (
    "io"
    "os"
    "fmt"
    "log"
    "strings"
    "net/http"
	"path/filepath"
    "github.com/PuerkitoBio/goquery"
)

const coursesDir = "courses"
const titleSelector = ".original-name"
const lessonsSelector = "#lessons-list li"

func main() {
	// Get the url
	url := os.Args[1]
	
    // Make HTTP GET request
    response, err := http.Get(url)
    if err != nil {
        log.Fatal(err)
    }
    defer response.Body.Close()

	// Get the DOM
	document, err := goquery.NewDocumentFromReader(response.Body)
    if err != nil {
        log.Fatal("Error loading HTTP response body. ", err)
	}

	// Get the course title
	courseTitle := document.Find(titleSelector).Text()
	
	// Course destination directory
	dir := coursesDir + "/" + courseTitle

	// Create course directory
	os.Mkdir(dir, 0700)

	// Find all video urls
	counter := 1
    document.Find(lessonsSelector).Each(func(index int, element *goquery.Selection) {
		title, err := element.Find("meta[itemprop=description]").Attr("content")
		if !err {
			fmt.Println("Cannot find title of element", counter)
			return
		}
		
		url, err := element.Find("link[itemprop=url]").Attr("href")
		if !err {
			fmt.Println("Cannot find url of element", counter)
			return
		}
		
		filename := fmt.Sprintf("%d%s%s", counter, ". ", title)
		download(url, dir, filename)
		counter++
    })
}

func download(url string, dir string, title string) {
	tokens := strings.Split(url, "/")
	ext := filepath.Ext(tokens[len(tokens)-1])
	filename := title + ext

	fmt.Println("Downloading", title, "...")

	output, err := os.Create(dir + "/" + filename)
	if err != nil {
		fmt.Println("Error while creating", filename, "-", err)
		return
	}
	defer output.Close()

	response, err := http.Get(url)
	if err != nil {
		fmt.Println("Error while downloading", url, "-", err)
		return
	}
	defer response.Body.Close()

	n, err := io.Copy(output, response.Body)
	if err != nil {
		fmt.Println("Error while downloading", url, "-", err)
		return
	}

	fmt.Println(n, "bytes downloaded.")
}
