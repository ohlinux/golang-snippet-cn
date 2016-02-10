package main

import (

    "fmt"

)

type USB interface {

    Name() string 
    Connecter
}

type Connecter interface {
    Connect()
}

type PhoneConnecter struct { 
        name string 
}


func (pc PhoneConnecter) Name() string {
    return pc.name
}

func (pc PhoneConnecter) Connect()  {
  fmt.Println("Connect:",pc.name)
}

func main() {
    var a USB
    a = PhoneConnecter{"PhoneConnecter"}
    a.Connect()
    Disconnect(a)
}

func Disconnect(usb USB) {
    if pc,ok:=usb.(PhoneConnecter);ok{
        fmt.Println("Disconnect:",pc.name)
        return
    }
    fmt.Println("Unknow decive.")
}
