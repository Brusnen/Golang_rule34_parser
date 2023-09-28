package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func get_pages(character string, page_num int) io.ReadCloser {
	client := &http.Client{}

	page := fmt.Sprintf("https://rule34.xxx/index.php?page=post&s=list&tags=%s+&pid=%d", character, page_num)

	req, err := http.NewRequest("GET", page, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Linux; Android 13; "+
		"22101316UG Build/TP1A.220624.014; wv) AppleWebKit/537.36 (KHTML, like Gecko) "+
		"Version/4.0 Chrome/117.0.0.0 Mobile Safari/537.36 [FB_IAB/FB4A;FBAV/432.0.0.29.102;]")

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	return resp.Body
}

func get_page(formatted_url string) io.ReadCloser {
	client := &http.Client{}

	page := formatted_url

	req, err := http.NewRequest("GET", page, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Linux; Android 13; "+
		"22101316UG Build/TP1A.220624.014; wv) AppleWebKit/537.36 (KHTML, like Gecko) "+
		"Version/4.0 Chrome/117.0.0.0 Mobile Safari/537.36 [FB_IAB/FB4A;FBAV/432.0.0.29.102;]")

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	return resp.Body
}

func get_single_pages(doc *goquery.Document) []string {
	pages := []string{}

	doc.Find("a").Each(func(index int, element *goquery.Selection) {
		href, _ := element.Attr("href")
		if strings.Contains(href, "/index.php?page=post&s=view&id=") {
			pages = append(pages, href)
		}
	})
	return pages
}

func get_image_pages(page string) string {
	formatted_url := fmt.Sprintf("https://rule34.xxx%s", page)
	doc_body := get_page(formatted_url)
	doc, _ := goquery.NewDocumentFromReader(doc_body)
	temp := doc.Find("[id=\"image\"]")
	res_link, _ := temp.Attr("data-cfsrc")
	return res_link
}

func download_image(url string, index int) error {
	response, err := http.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("Статус код ответа: %d", response.StatusCode)
	}

	file, err := os.Create(strconv.Itoa(index) + ".jpg")
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, response.Body)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	overall_pages := []string{}
	image_link := []string{}
	var author_name string
	var pic_num int
	fmt.Println("Enter author name")
	fmt.Scanln(&author_name)
	fmt.Println("Enter num of pictures")
	fmt.Scanln(&pic_num)

	for i := 0; i < pic_num; i = i + 42 {
		doc_body := get_pages(author_name, i)
		doc, _ := goquery.NewDocumentFromReader(doc_body)
		overall_pages = append(overall_pages, get_single_pages(doc)...)
	}
	for i := 0; i < len(overall_pages); i++ {
		res := get_image_pages(overall_pages[i])
		image_link = append(image_link, res)
		fmt.Println("scrapping:", i)
	}

	for i := 0; i < len(image_link); i++ {
		download_image(image_link[i], i)
		fmt.Println("saving: ", i)
	}
}
