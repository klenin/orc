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
)

const USER_COUNT = 20

var base = new(models.ModelManager)

func random(min, max int) int {
    rand.Seed(int64(time.Now().Second()))

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
    loadUsers()
    loadEvents()
    loadEventTypes()
    loadForms()
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

    for i := 0; i < USER_COUNT; i++ {
        rand.Seed(int64(i))
        userName := "user"+strconv.Itoa(i)
        userEmail := userName+"@mail.ru"

        result, regId := base.RegistrationController().Register(userName, "secret"+strconv.Itoa(i), userEmail, "user")
        if result != "ok" {
            utils.HandleErr("[loadUsers]: "+result, nil, nil)
            continue
        }

        query := `INSERT INTO param_values (param_id, value, date, user_id, reg_id)
            VALUES (4, $1, $2, NULL, $3);`
        db.Exec(query, []interface{}{userEmail, date, regId})

        for k := 5; k < 8; k++ {
            query := `INSERT INTO param_values (param_id, value, date, user_id, reg_id)
                VALUES (`+strconv.Itoa(k)+`, '`+userName+`', $1, NULL, $2);`
            db.Exec(query, []interface{}{date, regId})
        }

        query = `SELECT users.token FROM registrations
            INNER JOIN events ON registrations.event_id = events.id
            INNER JOIN faces ON faces.id = registrations.face_id
            INNER JOIN users ON users.id = faces.user_id
            WHERE events.id = $1 AND registrations.id = $2;`
        res := db.Query(query, []interface{}{1, regId})

        if len(res) == 0 {
            utils.HandleErr("[loadUsers]: ", errors.New("Data are not faund."), nil)
            continue
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
        rand.Seed(int64(i))
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
    eventTypeNames, _ := ioutil.ReadFile("./resources/event-type-name")
    eventTypeNamesSourse := strings.Split(string(eventTypeNames), "\n")
    topicality := []bool{true, false}
    for i := 0; i < len(eventTypeNamesSourse); i++ {
        //rand.Seed(int64(i))
        eventTypeName := strings.TrimSpace(eventTypeNamesSourse[i])
        params := map[string]interface{}{"name": eventTypeName, "description": "", "topicality": topicality[rand.Intn(2)]}
        base.EventTypes().LoadModelData(params).QueryInsert("").Scan()
    }
}

func loadForms() {
    formNames, _ := ioutil.ReadFile("./resources/form-name")
    formNamesSourse := strings.Split(string(formNames), "\n")
    for i := 0; i < len(formNamesSourse); i++ {
        formName := strings.TrimSpace(formNamesSourse[i])
        base.Forms().
            LoadModelData(map[string]interface{}{"name": formName, "personal": true}).
            QueryInsert("").
            Scan()
    }
}

func LoadParamTypes() {
    paramTypesNames, _ := ioutil.ReadFile("./resources/param-type-name")
    paramTypesSourse := strings.Split(string(paramTypesNames), "\n")
    for i := 0; i < len(paramTypesSourse); i++ {
        paramType := strings.TrimSpace(paramTypesSourse[i])
        base.ParamTypes().
            LoadModelData(map[string]interface{}{"name": paramType}).
            QueryInsert("").
            Scan()
    }
}
