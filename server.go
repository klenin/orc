package main

import (
    "flag"
    "github.com/klenin/orc/db"
    "database/sql"
    "github.com/klenin/orc/resources"
    "github.com/klenin/orc/router"
    "github.com/klenin/orc/mvc/controllers"
    "log"
    "net/http"
    "os"
)

var err error

func main() {
    log.Println("Server started.")

    db.DB, err = sql.Open("postgres", os.Getenv("DATABASE_URL"))
    defer db.DB.Close()

    if err != nil {
        log.Println("Error db connection: ", err.Error())
        os.Exit(1)
    }

    log.Println("DB CONNECTED")

    testData := flag.Bool("test-data", false, "to load test data")
    resetDB := flag.Bool("reset-db", false, "reset the database")
    flag.Parse()

    if *resetDB == true {
        new(controllers.BaseController).IndexController().Init(*testData)
        resources.LoadAdmin()
        resources.LoadParamTypes()
    }

    if *testData == true {
        resources.Load()
    }

    // base := new(controllers.BaseController)
    // base.Index().LoadContestsFromCats()

    http.Handle("/", new(router.FastCGIServer))
    http.HandleFunc("/wellcometoprofile/", controllers.WellcomeToProfile)
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
