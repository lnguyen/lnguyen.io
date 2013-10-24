package main

import (
  "github.com/gorilla/mux"
  "github.com/gorilla/handlers"
  "github.com/daaku/go.grace/gracehttp"
  "net/http"
  "os"
)
func HomeHandler(w http.ResponseWriter, r *http.Request) {
  http.ServeFile(w, r, "index.html")
  return
}

func main() {
  r := mux.NewRouter()
  r.HandleFunc("/", HomeHandler)


  r.HandleFunc("/assets/", func(w http.ResponseWriter, r *http.Request) {
    http.ServeFile(w, r, r.URL.Path[1:])
  })
  loggingHandler := handlers.LoggingHandler(os.Stdout, r)
  http.Handle("/", loggingHandler)

  gracehttp.Serve(
    &http.Server{Addr: ":8080", Handler: loggingHandler },
  )

}
