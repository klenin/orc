package resources

import (
    "github.com/orc/mvc/controllers"
    "github.com/orc/mvc/models"
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
    loadParamTypes()
}

func loadUsers() {
    base := new(controllers.BaseController)
    firstNamesFemaleRussian, _ := ioutil.ReadFile("./resources/first-name-female")
    firstNamesMaleRussian, _ := ioutil.ReadFile("./resources/first-name-male")
    lastNamesFemaleRussian, _ := ioutil.ReadFile("./resources/last-name-female")
    lastNamesMaleRussian, _ := ioutil.ReadFile("./resources/last-name-male")
    genders := []string{"male", "females"}
    for i := 0; i < USER_COUNT; i++ {
        rand.Seed(int64(i))
        gender := genders[rand.Intn(2)]
        //emailProviders := []string{"@gmail.com", "@hotmail.com", "@yandex.ru", "@mail.com"}
        var firstNameSource []byte
        var lastNameSource []byte
        if gender == "male" {
            firstNameSource = firstNamesMaleRussian
            lastNameSource = lastNamesMaleRussian
        } else {
            firstNameSource = firstNamesFemaleRussian
            lastNameSource = lastNamesFemaleRussian
        }
        firstName := strings.TrimSpace(strings.Split(string(firstNameSource), "\n")[rand.Intn(USER_COUNT)])
        lastName := strings.TrimSpace(strings.Split(string(lastNameSource), "\n")[rand.Intn(USER_COUNT)])
        //email := firstName + "_" + lastName + strconv.Itoa(rand.Intn(1024)) + emailProviders[len(emailProviders)-1]
        //println(email)
        base.Handler().HandleRegister("user"+strconv.Itoa(i), "secret"+strconv.Itoa(i), "user", firstName, lastName)
    }
    base.Handler().HandleRegister("admin", "password", "admin", "", "")
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
        params := []interface{}{eventName, dateStart, dateFinish, time, ""}
        entity := base.Events()
        entity.Insert(entity.GetColumnSlice(1), params)
    }
}

func loadEventTypes() {
    eventTypeNames, _ := ioutil.ReadFile("./resources/event-type-name")
    eventTypeNamesSourse := strings.Split(string(eventTypeNames), "\n")
    topicality := []bool{true, false}
    for i := 0; i < len(eventTypeNamesSourse); i++ {
        //rand.Seed(int64(i))
        eventTypeName := strings.TrimSpace(eventTypeNamesSourse[i])
        params := []interface{}{eventTypeName, "", topicality[rand.Intn(2)]}
        entity := base.EventTypes()
        entity.Insert(entity.GetColumnSlice(1), params)
    }
}

func loadForms() {
    formNames, _ := ioutil.ReadFile("./resources/form-name")
    formNamesSourse := strings.Split(string(formNames), "\n")
    for i := 0; i < len(formNamesSourse); i++ {
        formName := strings.TrimSpace(formNamesSourse[i])
        entity := base.Forms()
        entity.Insert(entity.GetColumnSlice(1), []interface{}{formName})
    }
}

func loadParamTypes() {
    paramTypesNames, _ := ioutil.ReadFile("./resources/param-type-name")
    paramTypesSourse := strings.Split(string(paramTypesNames), "\n")
    for i := 0; i < len(paramTypesSourse); i++ {
        paramType := strings.TrimSpace(paramTypesSourse[i])
        entity := base.ParamTypes()
        entity.Insert(entity.GetColumnSlice(1), []interface{}{paramType})
    }
}
