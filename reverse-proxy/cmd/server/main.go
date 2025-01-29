package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

func checkErr(err error) {
	if err != nil {
		fmt.Fprint(os.Stderr, err)
	}
}

func main() {
	err := godotenv.Load()
	checkErr(err)

	subDomain := ""

	bUrl := os.Getenv("BUCKET_URL")

	http.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		subDomain = strings.Split(r.Host, ".")[0]

		target := bUrl + subDomain + "/index.html"

		targetUrl, err := url.Parse(target)
		checkErr(err)

		r.URL.Host = targetUrl.Host
		r.URL.Scheme = targetUrl.Scheme
		r.URL.Path = "_output/" + subDomain + "/index.html"
		r.RequestURI = ""

		originServerResponse, err := http.DefaultClient.Do(r)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(rw, err)
			return
		}

		rw.WriteHeader(http.StatusOK)
		io.Copy(rw, originServerResponse.Body)
	})

	fmt.Println("reverse proxy running on PORT:8000")
	err = http.ListenAndServe(":8000", nil)
	checkErr(err)
}
