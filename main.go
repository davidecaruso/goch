package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const coursesDir = "courses"
const courseTitleSelector = ".original-name"

const lessonsSelector = "#lessons-list li"

const lessonTitleSelector = "meta[itemprop=description]"
const lessonTitleAttribute = "content"

const lessonURLSelector = "link[itemprop=url]"
const lessonURLAttribute = "href"

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("Usage : %s url \n ", os.Args[0])
		os.Exit(1)
	}

	// Check if the url is valid
	u, err := url.ParseRequestURI(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	// Parse url to string
	url := u.String()

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
	courseTitle := document.Find(courseTitleSelector).Text()

	// Create the course directory
	courseDir := coursesDir + "/" + courseTitle
	os.Mkdir(courseDir, 0700)

	// Find all video urls
	counter := 1
	document.Find(lessonsSelector).Each(func(index int, element *goquery.Selection) {
		title, err := element.Find(lessonTitleSelector).Attr(lessonTitleAttribute)
		if !err {
			fmt.Println("Cannot find title of element", counter)
			counter++
			return
		}

		url, err := element.Find(lessonURLSelector).Attr(lessonURLAttribute)
		if !err {
			fmt.Println("Cannot find url of element", counter)
			counter++
			return
		}

		filename := fmt.Sprintf("%d%s%s", counter, ". ", title)
		download(url, courseDir, filename)
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
