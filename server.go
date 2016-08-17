package main

import (
    "flag"
    "log"
    "net/http"
    "database/sql"
    "github.com/klenin/orc/db"
    "github.com/klenin/orc/initial"
    "github.com/klenin/orc/router"
    "github.com/klenin/orc/mvc/controllers"
    "github.com/klenin/orc/config"
)

var err error

func main() {
    db.DB, err = sql.Open("postgres", config.GetValue("DATABASE_URL"))
    defer db.DB.Close()

    if err != nil {
        log.Fatalln("Error DB open:", err.Error())
    }

    if err = db.DB.Ping(); err != nil {
        log.Fatalln("Error DB ping:", err.Error())
    }

    log.Println("Connected to DB")

    testData := flag.Bool("test-data", false, "to load test data")
    resetDB := flag.Bool("reset-db", false, "reset the database")
    flag.Parse()

    initial.Init(*resetDB, *testData)

    // base := new(controllers.BaseController)
    // base.Index().LoadContestsFromCats()

    http.Handle("/", new(router.FastCGIServer))
    http.HandleFunc("/wellcometoprofile/", controllers.WellcomeToProfile)

    fileServer := http.FileServer(http.Dir("./static"))
    http.Handle("/js/", fileServer)
    http.Handle("/css/", fileServer)
    http.Handle("/img/", fileServer)

    addr := config.GetValue("HOSTNAME") + ":" + config.GetValue("PORT")
    log.Println("Server listening on", addr)
    log.Fatalln("Error listening:", http.ListenAndServe(addr, nil))
}
