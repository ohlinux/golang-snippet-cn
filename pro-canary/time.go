package main 
import (
    "fmt"
    "time"
)


func main() {
    the_time := time.Date(2014, 1, 7, 5, 50, 4, 0, time.Local)
    unix_time := the_time.Unix()
    fmt.Println(unix_time)


    the_time, err := time.Parse("2006-01-02 15:04:05", "2014-01-08 09:04:41")
    if err == nil {
        unix_time := the_time.Unix()
        fmt.Println(unix_time)
    }
}
