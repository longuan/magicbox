package main

import (
	"flag"
	"fmt"
	"net/http"
)

var (
	file = flag.String("file", "", "metric file to be parsed")
	port = flag.Int("port", 8067, "listening port for http server")
)

func main() {
	flag.Parse()

	if *file == "" {
		fmt.Println("need to provide file path")
		return
	}

	err := readFile(*file)
	if err != nil {
		panic(err)
	}

	ftdcServer := http.NewServeMux()
	ftdcServer.HandleFunc("/metrics", metricsHandler)
	ftdcServer.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusOK) })
	ftdcServer.Handle("/query", &queryHandler{})

	err = http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", *port), ftdcServer)
	if err != nil {
		panic(err)
	}

}
