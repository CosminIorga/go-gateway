package main

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/cssivision/reverseproxy"
)

func main() {
	http.ListenAndServe(":8080", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		target := r.URL.Query()["target"][0]

		pathS, err := url.Parse("http://" + target + ".prd.adoreme.com")

		fmt.Println(pathS, err)

		if err != nil {
			panic(err)
			return
		}
		proxy := reverseproxy.NewReverseProxy(pathS)
		proxy.ServeHTTP(w, r)
	}))
}
