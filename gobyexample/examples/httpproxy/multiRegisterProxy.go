package main
 
import (
        "log"
        "math/rand"
        "net"
        "net/http"
        "net/http/httputil"
        "net/url"
        "time"
	"fmt"
	"strings"
)
 
type Registry map[string][]string
 
func extractNameVersion(target *url.URL) (name, version string, err error) {
        path := target.Path
        // Trim the leading `/`
        if len(path) > 1 && path[0] == '/' {
                path = path[1:]
        }
        // Explode on `/` and make sure we have at least
        // 2 elements (service name and version)
        tmp := strings.Split(path, "/")
        if len(tmp) < 2 {
                return "", "", fmt.Errorf("Invalid path")
        }
        name, version = tmp[0], tmp[1]
        // Rewrite the request's path without the prefix.
        target.Path = "/" + strings.Join(tmp[2:], "/")
        return name, version, nil
}
 

func NewMultipleHostReverseProxy(reg Registry) *httputil.ReverseProxy {
    director := func(req *http.Request) {
        name, version, err := extractNameVersion(req.URL)
        if err != nil {
            log.Print(err)
            return
        }
        req.URL.Scheme = "http"
        req.URL.Host = name + "/" + version
    }
    return &httputil.ReverseProxy{
        Director: director,
        Transport: &http.Transport{
            Proxy: http.ProxyFromEnvironment,
            Dial: func(network, addr string) (net.Conn, error) {
                // Trim the `:80` added by Scheme http.
                addr = strings.Split(addr, ":")[0]
                endpoints := reg[addr]
                if len(endpoints) == 0 {
                    return nil, fmt.Errorf("Service/Version not found")
                }
                return net.Dial(network, endpoints[rand.Int()%len(endpoints)])
            },
            TLSHandshakeTimeout: 10 * time.Second,
        },
    }
}

func main() {
	repoFrontend := "/Users/ajian/Code/myself/icode.baidu.com/ops/front/ops/output" 
	http.Handle("/", http.FileServer(http.Dir(repoFrontend)))

        proxy := NewMultipleHostReverseProxy(Registry{
                        "serviceone/v1": {"localhost:9091"},
                        "serviceone/v2": {"localhost:9092"},
        })
        log.Fatal(http.ListenAndServe(":9090", proxy))
}
