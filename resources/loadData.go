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
    "fmt"
)

const USER_COUNT = 20
const EVENTS_COUNT = 20

var base = new(models.ModelManager)

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
    eventNames := readStringsFromFile("./resources/event-type.txt")
    subjectNames := readStringsFromFile("./resources/event-subject.txt")
    for i := 0; i < EVENTS_COUNT; i++ {
        eventName := fmt.Sprintf("%s по дисциплине \"%s\"",
            eventNames[rand.Intn(len(eventNames))], subjectNames[rand.Intn(len(subjectNames))])

        var secInYear int64 = 365 * 24 * 60 * 60
        timeRangeFrom := time.Now().Unix() - secInYear * 5
        timeRangeTo := time.Now().Unix() + secInYear
        timeStart := time.Unix(timeRangeFrom + rand.Int63n(timeRangeTo - timeRangeFrom), 0)
        timeFinish := time.Unix(timeStart.Unix() + rand.Int63n(7 * 24 * 60 * 60), 0)
        params := map[string]interface{}{
            "name": eventName,
            "date_start": timeStart.Format("2006-01-02"),
            "date_finish": timeFinish.Format("2006-01-02"),
            "time": timeStart.Format("15:04:05"),
            "team": rand.Int() % 3 == 2,
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
    paramTypes := readStringsFromFile("./resources/param-type-name.txt")
    for _, paramType := range(paramTypes) {
        base.ParamTypes().
            LoadModelData(map[string]interface{}{"name": paramType}).
            QueryInsert("").
            Scan()
    }
}
