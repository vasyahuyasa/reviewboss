package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi"
)

type Webserver struct {
	eng *engine
}

func (srv *Webserver) routes() http.Handler {
	mux := chi.NewMux()
	mux.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "<html><body>")
			next.ServeHTTP(w, r)
			fmt.Fprintf(w, "</body></html>")
		})
	})
	mux.Get("/", srv.indexHandler)

	return mux
}

func (srv *Webserver) indexHandler(w http.ResponseWriter, r *http.Request) {
	srv.eng.mu.Lock()
	defer srv.eng.mu.Unlock()

	fmt.Fprintf(w, "<table border=1><tr><td>Date</td><td>Project</td><td>Link</td><td>Assignee</td></tr>")
	for _, m := range srv.eng.mergeRequests {
		assignee := "none"
		if m.Assignee != nil {
			assignee = "@" + m.Assignee.TelegramName
		}

		fmt.Fprintf(w,
			"<tr><td>%s</td><td>%s</td><td><a href=\"%s\">%s</a></td><td>%s</td></tr>",
			m.AddedOn.Format(time.RFC822),
			m.GitlabProject,
			m.Link,
			m.GitlabTitle,
			assignee,
		)
	}
	fmt.Fprintf(w, "</table>")
}

func (srv *Webserver) Run() error {
	log.Println("web server", "http://localhost:8080")
	return http.ListenAndServe("0.0.0.0:8080", srv.routes())
}

func NewWebserver(eng *engine) *Webserver {
	return &Webserver{
		eng: eng,
	}

}
