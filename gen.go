package main

import (
    // "io/ioutil"
    "fmt"
    "github.com/orc/db"
)

const TABLE = "#table"

type Route struct {
    Controller string
    Method     string
    Args       []string
}

func (this Route) GetRoute() {
    var result string

    // if this.Method == "" {
    //     fmt.Println("/"+this.Controller)
    //     return
    // }

    if this.Args == nil {
        fmt.Println("/"+this.Controller+"/"+this.Method)
        return
    }

    for _, v := range this.Args {
        if v == "#table" {
            for _, t := range db.Tables {
                result = "/"+this.Controller+"/"+this.Method+"/"+t
                fmt.Println(result)
                // var a[50]byte
                // copy(a[:], result)
                // err := ioutil.WriteFile("routes.txt", a, 064)
                // if err != nil {
                //     panic(err)
                // }
            }
        } else {
            result = "/"+this.Controller+"/"+this.Method+"/"+v
            fmt.Println(result)
            // var a[50]byte
            // copy(a[:], result)
            // err := ioutil.WriteFile("routes.txt", a, "c")
            // if err != nil {
            //     panic(err)
            // }
        }
    }
}


func generator() {
	/* Routes of User*/

    route := Route{Controller: "handler", Method: "getHistoryRequest", Args: nil}
    route.GetRoute()
    route = Route{Controller: "handler", Method: "saveUserRequest", Args: nil}
    route.GetRoute()
    route = Route{Controller: "handler", Method: "getList", Args: nil}
    route.GetRoute()
    route = Route{Controller: "handler", Method: "index", Args: nil}
    route.GetRoute()

    /* Routes of Admin*/

    route = Route{Controller: "gridhandler", Method: "select", Args: []string{TABLE}}
    route.GetRoute()
    route = Route{Controller: "gridhandler", Method: "load", Args: []string{TABLE}}
    route.GetRoute()
    route = Route{Controller: "gridhandler", Method: "edit", Args: []string{TABLE}}
    route.GetRoute()

    route = Route{Controller: "gridhandler", Method: "resetPassword", Args: nil}
    route.GetRoute()
    route = Route{Controller: "gridhandler", Method: "getSubTable", Args: nil}
    route.GetRoute()
}

func main() {
    generator()    
}