package main

import (
	"bytes"
	"fmt"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"time"

	vegeta "github.com/tsenart/vegeta/v12/lib"
)

var count int

const BASEURL = "http://localhost:8080/api/v1"

func generateEmail() string {
	email := fmt.Sprintf("user%d@example.com", count)
	count++

	return email
}

func createForm(email string) (header string, bd []byte, er error) {
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	err := writer.WriteField("email", email)
	if err != nil {
		return "", nil, err
	}
	err = writer.Close()
	if err != nil {
		return "", nil, err
	}
	return writer.FormDataContentType(), body.Bytes(), nil
}

func main() {
	targeter := func(tgt *vegeta.Target) error {
		if time.Now().Unix()%2 == 0 {
			tgt.Method = "GET"
			tgt.URL = fmt.Sprint(BASEURL, "/rate")
			tgt.Body = nil
			tgt.Header = nil
		} else {

			uniqueEmail := generateEmail()
			content, body, err := createForm(uniqueEmail)
			if err != nil {
				return err
			}
			tgt.Method = "POST"
			tgt.URL = fmt.Sprint(BASEURL, "/subscribe")
			tgt.Body = body
			tgt.Header = make(http.Header)
			tgt.Header.Set("Content-Type", content)
		}
		return nil
	}

	duration := 30 * time.Second
	initialRate := 10

	pacer := vegeta.LinearPacer{
		StartAt: vegeta.Rate{Freq: initialRate, Per: time.Second},
		Slope:   1,
	}

	attacker := vegeta.NewAttacker()

	metrics := &vegeta.Metrics{}

	results := attacker.Attack(targeter, pacer, duration, "Load Test")
	for res := range results {
		metrics.Add(res)
	}

	metrics.Close()
	file, err := os.Create("report.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	res := vegeta.NewTextReporter(metrics)
	err = res.Report(file)
	if err != nil {
		log.Println(err)
	}
}
