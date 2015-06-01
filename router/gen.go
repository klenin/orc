// package main
package router

import (
    "fmt"
    "github.com/orc/db"
    "os"
)

const TABLE = "#table"
const prefix = "ProxyPass "
var host = os.Getenv("HOST")
var link = os.Getenv("LINK")

type Route struct {
    Controller string
    Method     string
    Args       []string
}

func (this Route) GetRoute() {
    var result string

    if this.Method == "" {
        fmt.Println("/"+this.Controller)
        return
    }

    if this.Args == nil {
        result = "/"+this.Controller+"/"+this.Method
        fmt.Println(prefix+result+" "+host+result)
        return
    }

    for _, v := range this.Args {
        if v == "#table" {
            for _, t := range db.Tables {
                result = "/"+this.Controller+"/"+this.Method+"/"+t
                fmt.Println(prefix+result+" "+host+result)
            }
        } else {
            result = "/"+this.Controller+"/"+this.Method+"/"+v
            fmt.Println(prefix+result+" "+host+result)
        }
    }
}


func Generate() {

    /* js css img */

    fmt.Println(prefix+"/js"+" "+host+"/js")
    fmt.Println(prefix+"/css"+" "+host+"/css")
    fmt.Println(prefix+"/img"+" "+host+"/img")

    /* Routes of User*/

    route := Route{Controller: "handler", Method: "gethistoryrequest", Args: nil}
    route.GetRoute()
    route = Route{Controller: "handler", Method: "saveuserrequest", Args: nil}
    route.GetRoute()
    route = Route{Controller: "handler", Method: "getlist", Args: nil}
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

    route = Route{Controller: "gridhandler", Method: "resetpassword", Args: nil}
    route.GetRoute()
    route = Route{Controller: "gridhandler", Method: "getsubtable", Args: nil}
    route.GetRoute()

    fmt.Println("ProxyPassReverse /"+link+" "+host)
}

// func main() {
//     Generate()
// }
