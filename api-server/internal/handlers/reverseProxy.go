package handlers

import (
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/0xMishra/relay/api-server/internal/utils"
	"github.com/joho/godotenv"
)

var envErr = godotenv.Load()

var bucketURL = os.Getenv("BUCKET_URL")

func ReverseProxy(rw http.ResponseWriter, r *http.Request) {
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
	utils.CheckErr(err, true)

	resp, err := http.Get(targetURL.String())
	utils.CheckErr(err, true)

	for key, value := range resp.Header {
		rw.Header()[key] = value
	}

	rw.WriteHeader(resp.StatusCode)

	_, err = io.Copy(rw, resp.Body)
	utils.CheckErr(err, false)

	defer resp.Body.Close()
}
