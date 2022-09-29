package main

import (
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"

	"github.com/andrdru/sitemap"
)

func main() {
	router := httprouter.New()

	s := sitemap.NewSitemap()

	urlset1 := sitemap.NewURLSet()
	urlset1.URL = append(urlset1.URL, sitemap.URL{
		Loc: "http://localhost:8080/test1",
	})

	urlset2 := sitemap.NewURLSet()
	urlset2.URL = append(urlset2.URL, sitemap.URL{
		Loc: "http://localhost:8080/test2",
	})

	_ = s.AddURLSet(urlset1)
	_ = s.AddURLSet(urlset2)

	router.GET("/sitemap:ignoredParam", func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
		s.Handler()(writer, request)
	})

	err := http.ListenAndServe(":8080", router)
	if err != nil {
		log.Fatal(err)
	}
}
