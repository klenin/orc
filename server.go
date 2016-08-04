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
    http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("./static/js"))))
    http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("./static/css"))))
    http.Handle("/img/", http.StripPrefix("/img/", http.FileServer(http.Dir("./static/img"))))
    http.Handle("/vendor/", http.StripPrefix("/vendor/", http.FileServer(http.Dir("./static/vendor"))))

    addr := config.GetValue("HOSTNAME") + ":" + config.GetValue("PORT")
    log.Println("Server listening on", addr)
    log.Fatalln("Error listening:", http.ListenAndServe(addr, nil))
}
