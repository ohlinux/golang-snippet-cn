package main
 
import (
        "net/http"
        "net/http/httputil"
        "net/url"
)
 
func main() {
        proxy := httputil.NewSingleHostReverseProxy(&url.URL{
                Scheme: "http",
                Host:   "localhost:8086",
        })
        http.ListenAndServe(":9090", proxy)
}
