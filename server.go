package main

import (
  "github.com/gorilla/mux"
  "github.com/rcrowley/goagain"
  "net/http"
  "net"
  "time"
  "log"
  "os"
)
func HomeHandler(w http.ResponseWriter, r *http.Request) {
  http.ServeFile(w, r, "index.html")
  return
}

func main() {
  r := mux.NewRouter()
  r.HandleFunc("/", HomeHandler)

  http.HandleFunc("/assets/", func(w http.ResponseWriter, r *http.Request) {
    http.ServeFile(w, r, r.URL.Path[1:])
  })
  http.Handle("/", r)

  var (
    err error
    l net.Listener
    ppid int
  )

  l, ppid, err = goagain.GetEnvs()

  if err != nil {
    // Listen on a TCP or a UNIX domain socket (the latter is commented).
    laddr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:48879")
    if nil != err {
      log.Fatalln(err)
    }
    log.Printf("listening on %v", laddr)
    l, err = net.ListenTCP("tcp", laddr)
    if nil != err {
      log.Fatalln(err)
    }

    // Accept connections in a new goroutine.
    go http.Serve(l,nil);
  } else {
    // Resume listening and accepting connections in a new goroutine.
    log.Printf("resuming listening on %v", l.Addr())
    go http.Serve(l,nil);

    // Kill the parent, now that the child has started successfully.
    if err := goagain.KillParent(ppid); nil != err {
      log.Fatalln(err)
    }
  }

  // Block the main goroutine awaiting signals.
  if err := goagain.AwaitSignals(l); nil != err {
    log.Fatalln(err)
  }
  log.Println(os.Args[0])

  // Do whatever's necessary to ensure a graceful exit like waiting for
  // goroutines to terminate or a channel to become closed.

  // In this case, we'll simply stop listening and wait one second.
  if err := l.Close(); nil != err {
    log.Fatalln(err)
  }
  time.Sleep(1e9)



}
