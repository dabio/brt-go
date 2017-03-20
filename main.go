package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"time"

	"github.com/dabio/brt-go/models"
	_ "github.com/lib/pq"
	"github.com/pinub/mux"
	stathat "github.com/stathat/go"
)

type context struct {
	db        models.Datastore
	templates *template.Template
}

func (c *context) redirect(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

func (c *context) index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	data := struct {
		Year int
	}{
		time.Now().Year(),
	}

	c.render(w, "index", data)
}

func (c *context) calendar(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/calendar")

	year, _ := strconv.Atoi(time.Now().Format("2006"))
	events, err := c.db.AllEvents(year)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	c.render(w, "rennen.ics", events)
}

func (c *context) render(w http.ResponseWriter, tmpl string, data interface{}) {
	if err := c.templates.ExecuteTemplate(w, tmpl+".tmpl", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	db, err := models.NewDB("postgres", os.Getenv("DATABASE"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	c := context{
		templates: template.Must(template.ParseGlob("./views/*.tmpl")),
		db:        db,
	}

	h := mux.New()
	h.Get("/", track(c.index))
	h.Get("/rennen.ics", track(c.calendar))
	h.Get("/rennen", track(c.redirect))
	h.Get("/team", track(c.redirect))
	h.Get("/kontakt", track(c.redirect))
	h.Get("/news", track(c.redirect))

	m := http.NewServeMux()
	m.Handle("/css/", http.FileServer(http.Dir("./static/")))
	m.Handle("/img/", http.FileServer(http.Dir("./static/")))
	m.Handle("/", h)

	s := &http.Server{
		Addr:         ":" + os.Getenv("PORT"),
		Handler:      m,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		// Go1.8
		// IdleTimeout: 120 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	log.Fatal(s.ListenAndServe())
}

func track(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func(start time.Time, r *http.Request) {
			elapsed := time.Since(start)
			if os.Getenv("ENV") == "production" {
				sh := os.Getenv("STATHAT")
				stathat.PostEZCountOne("brt visit", sh)
				stathat.PostEZValue("brt duration", sh, float64(elapsed/time.Microsecond))
			} else {
				log.Printf("%s %s %s", r.Method, r.URL, elapsed)
			}
		}(time.Now(), r)

		fn(w, r)
	}
}
