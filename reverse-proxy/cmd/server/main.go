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

func checkErr(err error, fatal bool) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		if fatal {
			os.Exit(1)
		}
	}
}

func main() {
	err := godotenv.Load()
	checkErr(err, true)

	bucketURL := os.Getenv("BUCKET_URL")
	if bucketURL == "" {
		fmt.Fprintln(os.Stderr, "BUCKET_URL not set in environment")
		os.Exit(1)
	}

	if !strings.HasSuffix(bucketURL, "/") {
		bucketURL += "/"
	}

	http.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		hostParts := strings.Split(r.Host, ".")
		if len(hostParts) < 1 {
			http.Error(rw, "Invalid subdomain", http.StatusBadRequest)
			return
		}
		subDomain := hostParts[0]

		target := bucketURL + subDomain + r.URL.Path
		if r.URL.Path == "/" {
			target += "index.html"
		}

		targetURL, err := url.Parse(target)
		checkErr(err, true)

		resp, err := http.Get(targetURL.String())
		checkErr(err, true)

		for key, value := range resp.Header {
			rw.Header()[key] = value
		}

		rw.WriteHeader(resp.StatusCode)

		_, err = io.Copy(rw, resp.Body)
		checkErr(err, false)

		defer resp.Body.Close()
	})

	fmt.Println("Reverse proxy running on PORT: 8000")
	err = http.ListenAndServe(":8000", nil)
	checkErr(err, true)
}
