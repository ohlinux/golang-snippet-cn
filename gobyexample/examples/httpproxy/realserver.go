package main
 
import (
        "log"
        "net/http"
        "os"
        "strconv"
	"io"
)
 
func main() {
        if len(os.Args) != 2 {
                log.Fatalf("Usage: %s <port>", os.Args[0])
        }
        if _, err := strconv.Atoi(os.Args[1]); err != nil {
                log.Fatalf("Invalid port: %s (%s)\n", os.Args[1], err)
        }
 
        http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		io.WriteString(w, "hello "+os.Args[1])
                println("--->", os.Args[1], req.URL.String())
        })
        http.ListenAndServe(":"+os.Args[1], nil)
}
