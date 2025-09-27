package main

import (
    "fmt"
    "github.com/go-chi/chi/v5"
    "io"
    "log"
    "net/http"
)

func main() {
    mux := chi.NewRouter()
    mux.Get("/", index)
    mux.Get("/hello", hello)
    mux.Get("/headers", headers)
    mux.Post("/body", body)
    http.ListenAndServe(":8000", mux)
}

func index(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("Hello World"))
}

func hello(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Method : %s, Host : %s\n", r.Method, r.Host)
    fmt.Fprintf(w, "Path : %s, Query : %s\n", r.URL.Path, r.URL.Query())
}

func headers(w http.ResponseWriter, r *http.Request) {
    for k, v := range r.Header {
        fmt.Fprintf(w, "%s: %s\n", k, v)
    }
}

func body(w http.ResponseWriter, r *http.Request) {
    body, err := io.ReadAll(r.Body)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Fprintf(w, "%s", body)
}
