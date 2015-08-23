package resources

import (
    "errors"
    "github.com/orc/db"
    "github.com/orc/mvc/controllers"
    "github.com/orc/mvc/models"
    "github.com/orc/utils"
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

func addDate(d, m, y int) string {
    return strconv.Itoa(d) + "-" + strconv.Itoa(m) + "-" + strconv.Itoa(y)
}

func addTime(h, m, s int) string {
    return strconv.Itoa(h) + ":" + strconv.Itoa(m) + ":" + strconv.Itoa(s)
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

    result, regId := base.RegistrationController().Register("admin", "password", "secret.oasis.3805@gmail.com", "admin")
    if result != "ok" {
        utils.HandleErr("[LoadAdmin]: "+result, nil, nil)
        return
    }

    for k := 5; k < 8; k++ {
        query := `INSERT INTO param_values (param_id, value, date, user_id, reg_id)
            VALUES (`+strconv.Itoa(k)+`, 'admin', $1, NULL, $2);`
        db.Exec(query, []interface{}{date, regId})
    }

    query := `SELECT users.token FROM registrations
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

        result, regId := base.RegistrationController().Register(userName, "secret"+strconv.Itoa(i), "", "user")
        if result != "ok" {
            utils.HandleErr("[loadUsers]: "+result, nil, nil)
            continue
        }

        for k := 5; k < 8; k++ {
            query := `INSERT INTO param_values (param_id, value, date, user_id, reg_id)
                VALUES (`+strconv.Itoa(k)+`, '`+userName+`', $1, NULL, $2);`
            db.Exec(query, []interface{}{date, regId})
        }

        query := `SELECT users.token FROM registrations
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
        dateStart := addDate(random(1894, 2014), random(1, 12), random(1, 28))
        dateFinish := addDate(random(1894, 2014), random(1, 12), random(1, 28))
        time := addTime(random(0, 11), random(1, 60), random(1, 60))
        params := map[string]interface{}{"name": eventName, "data_start": dateStart, "date_finish": dateFinish, "time": time, "url": ""}
        entity := base.Events()
        entity.LoadModelData(params)
        db.QueryInsert(entity, "").Scan()
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
        entity := base.EventTypes()
        entity.LoadModelData(params)
        db.QueryInsert(entity, "").Scan()
    }
}

func loadForms() {
    formNames, _ := ioutil.ReadFile("./resources/form-name")
    formNamesSourse := strings.Split(string(formNames), "\n")
    for i := 0; i < len(formNamesSourse); i++ {
        formName := strings.TrimSpace(formNamesSourse[i])
        entity := base.Forms()
        entity.LoadModelData(map[string]interface{}{"name": formName, "personal": true})
        db.QueryInsert(entity, "").Scan()
    }
}

func LoadParamTypes() {
    paramTypesNames, _ := ioutil.ReadFile("./resources/param-type-name")
    paramTypesSourse := strings.Split(string(paramTypesNames), "\n")
    for i := 0; i < len(paramTypesSourse); i++ {
        paramType := strings.TrimSpace(paramTypesSourse[i])
        entity := base.ParamTypes()
        entity.LoadModelData(map[string]interface{}{"name": paramType})
        db.QueryInsert(entity, "").Scan()
    }
}
