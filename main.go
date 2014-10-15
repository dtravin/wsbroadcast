package main

import (
	"flag"
	"fmt"
	"github.com/dtravin/wsbroadcast/wsb"
	"html/template"
	"log"
	"net/http"
	"os"
)

var ListenPort int

type Page struct {
	Title      string
	ListenPort int
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	var p = &Page{Title: "WS broadcast DEMO", ListenPort: ListenPort}
	t, _ := template.ParseFiles("html_templates/index.html")
	t.Execute(w, p)
}

func main() {

	listenPortPtr := flag.Int("l", 8080, "A port to listen for client connections")
	var inputStreamURI string
	flag.StringVar(&inputStreamURI, "i", "", "Input WS stream URL")
	flag.Parse()

	if inputStreamURI == "" {
		fmt.Fprintf(os.Stderr, "Specify input stream URL as -i=ws://IP:PORT\n")
		os.Exit(1)
	}
	ListenPort = *listenPortPtr

	log.Println(fmt.Sprintf("Listening server on :%d", ListenPort))

	http.HandleFunc("/", viewHandler)

	server := wsb.NewServer("/stream", inputStreamURI)
	go server.Listen()
	go server.ReadInput()
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", ListenPort), nil))
}
