package resources

import (
    "errors"
    "github.com/klenin/orc/mailer"
    "github.com/klenin/orc/db"
    "github.com/klenin/orc/mvc/controllers"
    "github.com/klenin/orc/mvc/models"
    "github.com/klenin/orc/utils"
    "io/ioutil"
    "math/rand"
    "strconv"
    "strings"
    "time"
    "log"
)

const USER_COUNT = 20

var base = new(models.ModelManager)

func random(min, max int) int {
    return rand.Intn(max-min) + min
}

func prepare(v1, v2, v3 string) (v1_, v2_, v3_ string) {
    if len(v1) <= 1 {
        v1 = "0" + v1
    }
    if len(v2) <= 1 {
        v2 = "0" + v2
    }
    if len(v3) <= 1 {
        v3 = "0" + v3
    }

    return v1, v2, v3
}

func addDate(d, m, y string) string {
    d, m, y = prepare(d, m, y)

    return d + "-" + m + "-" + y
}

func addTime(h, m, s string) string {
    h, m, s = prepare(h, m, s)

    return h + ":" + m + ":" + s
}

func Load() {
    rand.Seed(time.Now().UnixNano())

    loadUsers()
    loadEvents()
    loadEventTypes()
    loadForms()
}

func readStringsFromFile(fileName string) []string {
    content, err := ioutil.ReadFile(fileName)
    if err != nil {
        log.Fatalln("loadData:", err.Error())
    }
    array := strings.Split(string(content), "\n")
    var r []string
    for _, str := range array {
        if str = strings.TrimSpace(str); str != "" {
            r = append(r, str)
        }
    }
    return r
}

func LoadAdmin() {
    base := new(controllers.BaseController)
    date := time.Now().Format("2006-01-02T15:04:05Z00:00")

    result, regId := base.RegistrationController().Register("admin", "password", mailer.Admin_.Email, "admin")
    if result != "ok" {
        utils.HandleErr("[LoadAdmin]: "+result, nil, nil)

        return
    }

    query := `INSERT INTO param_values (param_id, value, date, user_id, reg_id)
        VALUES (4, $1, $2, NULL, $3);`
    db.Exec(query, []interface{}{mailer.Admin_.Email, date, regId})

    for k := 5; k < 8; k++ {
        query := `INSERT INTO param_values (param_id, value, date, user_id, reg_id)
            VALUES (`+strconv.Itoa(k)+`, 'admin', $1, NULL, $2);`
        db.Exec(query, []interface{}{date, regId})
    }

    query = `SELECT users.token FROM registrations
        INNER JOIN events ON registrations.event_id = events.id
        INNER JOIN faces ON faces.id = registrations.face_id
        INNER JOIN users ON users.id = faces.user_id
        WHERE events.id = $1 AND registrations.id = $2;`
    res := db.Query(query, []interface{}{1, regId})

    if len(res) == 0 {
        utils.HandleErr("[LoadAdmin]: ", errors.New("Data are not faund."), nil)

        return
    }

    token := res[0].(map[string]interface{})["token"].(string)
    base.RegistrationController().ConfirmUser(token)
}

func loadUsers() {
    base := new(controllers.BaseController)
    date := time.Now().Format("2006-01-02T15:04:05Z00:00")

    type FullNames struct {
        firstNames, lastNames, patronymics []string
    }

    male := FullNames{
        firstNames: readStringsFromFile("./resources/first-name-male.txt"),
        lastNames: readStringsFromFile("./resources/last-name-male.txt"),
        patronymics: readStringsFromFile("./resources/patronymic-male.txt"),
    }
    female := FullNames{
        firstNames: readStringsFromFile("./resources/first-name-female.txt"),
        lastNames: readStringsFromFile("./resources/last-name-female.txt"),
        patronymics: readStringsFromFile("./resources/patronymic-female.txt"),
    }

    for i := 0; i < USER_COUNT; i++ {
        userName := "user" + strconv.Itoa(i + 1)
        userEmail := userName + "@example.com"

        result, regId := base.RegistrationController().Register(userName, "password", userEmail, "user")
        if result != "ok" {
            log.Fatalln("[loadUsers]:", result)
        }

        query := `INSERT INTO param_values (param_id, value, date, reg_id)
            VALUES ($1, $2, $3, $4);`

        db.Exec(query, []interface{}{4, userEmail, date, regId})
        var fullNames FullNames
        if rand.Int() % 2 == 1 {
            fullNames = male
        } else {
            fullNames = female
        }
        db.Exec(query, []interface{}{6, fullNames.firstNames[rand.Intn(len(fullNames.firstNames))], date, regId})
        db.Exec(query, []interface{}{5, fullNames.lastNames[rand.Intn(len(fullNames.lastNames))], date, regId})
        db.Exec(query, []interface{}{7, fullNames.patronymics[rand.Intn(len(fullNames.patronymics))], date, regId})

        query = `SELECT users.token FROM registrations
            INNER JOIN events ON registrations.event_id = events.id
            INNER JOIN faces ON faces.id = registrations.face_id
            INNER JOIN users ON users.id = faces.user_id
            WHERE events.id = $1 AND registrations.id = $2;`
        res := db.Query(query, []interface{}{1, regId})

        if len(res) == 0 {
            log.Fatalln("[loadUsers]:", "Data are not found")
        }

        token := res[0].(map[string]interface{})["token"].(string)
        base.RegistrationController().ConfirmUser(token)
    }
}

func loadEvents() {
    eventNames, _ := ioutil.ReadFile("./resources/event-name")
    subjectNames, _ := ioutil.ReadFile("./resources/subject-name")
    eventNameSource := strings.Split(string(eventNames), "\n")
    subjectNameSource := strings.Split(string(subjectNames), "\n")
    for i := 0; i < len(eventNameSource); i++ {
        eventName := strings.TrimSpace(eventNameSource[rand.Intn(len(eventNameSource))])
        eventName += " по дисциплине "
        eventName += "\"" + strings.TrimSpace(subjectNameSource[rand.Intn(len(subjectNameSource))]) + "\""
        dateStart := addDate(strconv.Itoa(random(1894, 2014)), strconv.Itoa(random(1, 12)), strconv.Itoa(random(1, 28)))
        dateFinish := addDate(strconv.Itoa(random(1894, 2014)), strconv.Itoa(random(1, 12)), strconv.Itoa(random(1, 28)))
        time := addTime(strconv.Itoa(random(0, 11)), strconv.Itoa(random(1, 60)), strconv.Itoa(random(1, 60)))
        params := map[string]interface{}{
            "name": eventName,
            "date_start": dateStart,
            "date_finish": dateFinish,
            "time": time,
            "team": false,
            "url": ""}
        base.Events().LoadModelData(params).QueryInsert("").Scan()
    }
}

func loadEventTypes() {
    eventTypes := readStringsFromFile("./resources/event-type-name")
    topicality := []bool{true, false}
    for _, eventType := range eventTypes {
        params := map[string]interface{}{"name": eventType, "description": "", "topicality": topicality[rand.Intn(2)]}
        base.EventTypes().LoadModelData(params).QueryInsert("").Scan()
    }
}

func loadForms() {
    formNames := readStringsFromFile("./resources/form-name")
    for _, formName := range(formNames) {
        base.Forms().
            LoadModelData(map[string]interface{}{"name": formName, "personal": true}).
            QueryInsert("").
            Scan()
    }
}

func LoadParamTypes() {
    paramTypes := readStringsFromFile("./resources/param-type-name")
    for _, paramType := range(paramTypes) {
        base.ParamTypes().
            LoadModelData(map[string]interface{}{"name": paramType}).
            QueryInsert("").
            Scan()
    }
}
