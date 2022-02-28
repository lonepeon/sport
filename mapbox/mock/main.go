package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc(
		"/styles/v1/mapbox/outdoors-v11/static/",
		func(w http.ResponseWriter, r *http.Request) {
			f, err := os.Open("./map1.png")
			defer f.Close()
			if err != nil {
				http.Error(w, fmt.Sprintf("can't open png: %s", err), http.StatusInternalServerError)
				return
			}
			io.Copy(w, f)
		},
	)

	log.Fatalln(http.ListenAndServe(":8080", nil))
}
