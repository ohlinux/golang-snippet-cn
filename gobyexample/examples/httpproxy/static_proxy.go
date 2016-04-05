package main

import (
    "log"
    "net/http"
    "net/http/httputil"
    "net/url"
)

func main() {
    targetAPI, err := url.Parse("http://10.208.48.11:8080")
    if err != nil {
        log.Fatal(err)
    }

    targetUI, err := url.Parse("http://10.208.48.11:8081")
    if err != nil {
        log.Fatal(err)
    }

    repoFrontend := "front/ops/output"
    http.Handle("/", http.FileServer(http.Dir(repoFrontend)))

    //http.Handle("/users/", http.StripPrefix("/users/", httputil.NewSingleHostReverseProxy(target)))
    http.Handle("/api/", httputil.NewSingleHostReverseProxy(targetAPI))
    http.Handle("/ui/", httputil.NewSingleHostReverseProxy(targetUI))

    log.Fatal(http.ListenAndServe(":9090", nil))
}
