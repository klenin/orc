package main

import (
    "flag"
    "database/sql"
    "github.com/orc/db"
    "github.com/orc/resources"
    "github.com/orc/router"
    "github.com/orc/mvc/controllers"
    "log"
    "net/http"
    "os"
)

var err error

func main() {
    log.Println("Server started.")

    db.DB, err = sql.Open("postgres", os.Getenv("DATABASE_URL"))

    if err != nil {
        log.Println("Error db connection: ", err.Error())
        os.Exit(1)
    }

    log.Println("DB CONNECTED")

    testData := flag.Bool("test-data", false, "to load test data")
    flag.Parse()

    // db.Init()
    controllers.CreateRegistrationEvent()
    resources.LoadAdmin()

    if *testData == true {
        resources.Load()
    }

    // base := new(controllers.BaseController)
    // base.Index().LoadContestsFromCats()

    http.Handle("/", new(router.FastCGIServer))
    http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("./static/js"))))
    http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("./static/css"))))
    http.Handle("/img/", http.StripPrefix("/img/", http.FileServer(http.Dir("./static/img"))))

    port := os.Getenv("PORT")
    if port == "" {
        port = "5000"
    }

    if err := http.ListenAndServe(":" + port, nil); err != nil {
        log.Println("Error listening: ", err.Error())
        os.Exit(1)
    }
}
