package main

import (
	"net/http"
	"strconv"
)

type apiConfig struct {
	numHits int
}

func (cfg *apiConfig) incrementHits(next http.Handler) http.Handler {
	return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request){
		cfg.numHits++

		next.ServeHTTP(w, r)
	})


}

func (cfg *apiConfig) resetHits() {
	cfg.numHits = 0
}

func (cfg *apiConfig) hitsHandler (w http.ResponseWriter, r *http.Request) {
	hitsString := strconv.Itoa(cfg.numHits)
	
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hits: " + hitsString))
}

func (cfg *apiConfig) resetHandler (w http.ResponseWriter, r *http.Request) {
	cfg.resetHits()
	
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
}

func main(){
	apiCFG := apiConfig{
		numHits: 0,
	}

	mux := http.NewServeMux()
	server := &http.Server{
		Handler: mux,
		Addr: "localhost:8080",
	}

	mux.Handle(
		"/app/*", 
		apiCFG.incrementHits(http.StripPrefix("/app", http.FileServer(http.Dir(".")))),
	)
	mux.HandleFunc("/healthz", HTTPHandler)
	mux.HandleFunc("/metrics", apiCFG.hitsHandler)
	mux.HandleFunc("/reset", apiCFG.resetHandler)

	server.ListenAndServe()
}


func HTTPHandler (w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

