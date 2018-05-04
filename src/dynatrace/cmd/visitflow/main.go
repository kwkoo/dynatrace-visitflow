package main

import (
	"dynatrace"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
)

var application string
var fileserver http.Handler

func route(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		get(w, r)
	} else if r.Method == http.MethodPost || r.Method == http.MethodPut {
		post(w, r)
	} else {
		fmt.Fprintf(w, "Unknown Method\n")
	}
}

func get(w http.ResponseWriter, r *http.Request) {
	log.Println("Serving URI", r.URL.Path)
	if r.URL.Path == "/api" {
		maxchildren := 9999
		linkvaluethreshold := 0
		m, _ := url.ParseQuery(r.URL.RawQuery)
		if mc := m["maxchildren"]; mc != nil && len(mc) > 0 {
			mci, err := strconv.Atoi(mc[0])
			if err == nil {
				maxchildren = mci
			}
		}
		if lvt := m["linkvaluethreshold"]; lvt != nil && len(lvt) > 0 {
			lvti, err := strconv.Atoi(lvt[0])
			if err == nil {
				linkvaluethreshold = lvti
			}
		}
		log.Println("Processing /api request maxchildren:", maxchildren, "linkvaluethreshold:", linkvaluethreshold)
		dynatrace.GenerateFlow(w, maxchildren, linkvaluethreshold)
	} else if r.URL.Path == "/dump" {
		dynatrace.PrintGraph(w)
	} else {
		fileserver.ServeHTTP(w, r)
	}
}

func post(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	visits := dynatrace.Parse(r.Body)
	log.Println("Parsed", len(visits), "visit(s)")
	for _, v := range visits {
		dynatrace.AddVisit(v, application)
	}

	fmt.Fprintf(w, "OK\n")
}

func main() {
	var dir string
	var port int

	cwd, err := os.Getwd()
	if err != nil {
		panic("Could not get current directory.")
	}
	flag.StringVar(&dir, "dir", cwd, "directory containing static HTML content")

	flag.IntVar(&port, "port", 8080, "HTTP listener port")
	flag.StringVar(&application, "application", "", "Application filter")
	flag.Parse()

	dynatrace.InitGraph()

	if len(application) == 0 {
		log.Println("No application filter")
	} else {
		log.Println("Application filter set to", application)
	}
	log.Println("Serving files from directory", dir)
	log.Println("Listening on port", port)

	fileserver = http.FileServer(http.Dir(dir))
	http.HandleFunc("/", route)
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
