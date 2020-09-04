package scrapper

import (
	"encoding/csv"
	"fmt"
	"main/errhandle"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type extractedJob struct {
	id       string
	location string
	title    string
	salary   string
	summary  string
}

// Scrape Indeed by Term
func Scrape(term string) {
	var baseURL string = "https://www.indeed.com/jobs?q=" + term + "&limit=50"
	startTime := time.Now()
	ch := make(chan []extractedJob)
	var jobs []extractedJob
	totalPages := getPages(baseURL)

	for i := 0; i < totalPages; i++ {
		go singlePage(i, baseURL, ch)
	}

	for i := 0; i < totalPages; i++ {
		extractedJob := <-ch
		jobs = append(jobs, extractedJob...)
	}

	writeJobs(jobs)
	fmt.Println("Done, exported to CSV", len(jobs))

	endTime := time.Now()
	fmt.Println("Finished in, ", endTime.Sub(startTime))
}

func singlePage(num int, url string, chn chan<- []extractedJob) {
	ch := make(chan extractedJob)
	var jobs []extractedJob
	pageURL := url + "&start=" + strconv.Itoa(num*50)
	fmt.Println("Requestting", pageURL)
	res, err := http.Get(pageURL)
	errhandle.CheckErr(err)
	errhandle.CheckStatusCode(res)

	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	errhandle.CheckErr(err)

	searchCards := doc.Find(".jobsearch-SerpJobCard")

	searchCards.Each(func(i int, card *goquery.Selection) {
		go extractJob(card, ch)
	})

	for i := 0; i < searchCards.Length(); i++ {
		job := <-ch
		jobs = append(jobs, job)
	}

	chn <- jobs
}

func extractJob(card *goquery.Selection, ch chan<- extractedJob) {
	id, _ := card.Attr("data-jk")
	title := CleanString(card.Find(".title>a").Text())
	location := CleanString(card.Find(".location").Text())
	salary := CleanString(card.Find(".salaryText").Text())
	summary := CleanString(card.Find(".summary").Text())

	ch <- extractedJob{
		id,
		location,
		title,
		salary,
		summary,
	}
}

// CleanString is now exported
func CleanString(str string) string {
	return strings.Join(strings.Fields(strings.TrimSpace(str)), " ")
}

func getPages(url string) int {
	res, err := http.Get(url)
	errhandle.CheckErr(err)
	errhandle.CheckStatusCode(res)

	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	errhandle.CheckErr(err)

	pages := 0
	doc.Find(".pagination").Each(func(i int, s *goquery.Selection) {
		pages = s.Find("a").Length()
	})
	return pages
}

func writeJobs(jobs []extractedJob) {
	wch := make(chan []string)
	file, err := os.Create("jobs.csv")
	errhandle.CheckErr(err)

	w := csv.NewWriter(file)
	defer w.Flush()

	headers := []string{"Link", "Title", "Location", "Salary", "Summary"}
	wErr := w.Write(headers)
	errhandle.CheckErr(wErr)

	for _, job := range jobs {
		go writeJobDetails(job, wch)
	}

	for i := 0; i < len(jobs); i++ {
		jobSlice := <-wch
		jwErr := w.Write(jobSlice)
		errhandle.CheckErr(jwErr)
	}
}

func writeJobDetails(job extractedJob, wch chan<- []string) {
	jobSlice := []string{"https://www.indeed.com/viewjob?jk=" + job.id, job.title, job.location, job.salary, job.summary}
	wch <- jobSlice
}
