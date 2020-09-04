package errhandle

import (
	"log"
	"net/http"
)

// CheckErr for handle erros
func CheckErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// CheckStatusCode for handle bad status code
func CheckStatusCode(res *http.Response) {
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}
}
