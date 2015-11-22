package controllers

import (
    "database/sql"
    "errors"
    "github.com/klenin/orc/db"
    "github.com/klenin/orc/mvc/models"
    "github.com/klenin/orc/sessions"
    "github.com/klenin/orc/utils"
    "html/template"
    "net/http"
    "strings"
)

const HASH_SIZE = 32

var err error

type BaseController struct{}

type Controller struct {
    Request  *http.Request
    Response http.ResponseWriter
}

type Model struct {
    Id           string
    TableName    string
    Caption      string
    Table        []interface{}
    Columns      []string
    ColNames     []string
    ColModel     []map[string]interface{}
    Sub          bool
    SubTableName string
    SubCaption   string
    SubColNames  []string
    SubColModel  []map[string]interface{}
}

func (*Controller) GetModel(tableName string) models.EntityInterface {
    return new(models.ModelManager).GetModel(tableName)
}

func (this *Controller) Render(filenames []string, tmpname string, data interface{}) {
    filenames = append(filenames, "mvc/views/header.html")
    filenames = append(filenames, "mvc/views/footer.html")
    tmpl, err := template.ParseFiles(filenames...)
    if err != nil {
        http.Error(this.Response, err.Error(), http.StatusInternalServerError)
    }
    if err := tmpl.ExecuteTemplate(this.Response, tmpname, data); err != nil {
        http.Error(this.Response, err.Error(), http.StatusInternalServerError)
    }
}

func (this *Controller) CheckSid() (id int, result error)  {
    userSid := sessions.GetValue("sid", this.Request)
    if !sessions.CheckSession(this.Response, this.Request) || userSid == nil {
        return -1, errors.New("Данные в куках отсутствуют.")
    }

    if err := this.GetModel("users").
        LoadWherePart(map[string]interface{}{"sid": userSid}).
        SelectRow([]string{"id"}).
        Scan(&id);
        err != nil {
        return -1, errors.New("Данные в куках отсутствуют.")
    }

    return id, nil
}

func (this *Controller) isAdmin() bool {
    userId, err := this.CheckSid()
    if err != nil {
        return false
    }

    var role string
    err = this.GetModel("users").
        LoadWherePart(map[string]interface{}{"id": userId}).
        SelectRow([]string{"role"}).
        Scan(&role)
    if err != nil || role == "user" {
        return false
    }

    return role == "admin"
}

func (*Controller) regExists(userId, eventId int) int {
    var regId int
    query := `SELECT registrations.id
        FROM registrations
        INNER JOIN events ON events.id = registrations.event_id
        INNER JOIN faces ON faces.id = registrations.face_id
        INNER JOIN users ON users.id = faces.user_id
        WHERE users.id = $1 AND events.id = $2;`
    err := db.QueryRow(query, []interface{}{userId, eventId}).Scan(&regId)
    if err != sql.ErrNoRows {
        return regId
    } else {
        return -1
    }
}

func WellcomeToProfile(w http.ResponseWriter, r *http.Request) {
    newContreoller := new(BaseController).Handler()

    parts := strings.Split(r.URL.Path, "/")
    token := parts[len(parts)-1]

    var id int
    err := newContreoller.GetModel("users").
        LoadWherePart(map[string]interface{}{"token": token}).
        SelectRow([]string{"id"}).
        Scan(&id)
    if utils.HandleErr("[WellcomeToProfile]: ", err, w) || id == 0 {
        return
    }

    sid := utils.GetRandSeq(HASH_SIZE)
    params := map[string]interface{}{"sid": sid, "enabled": true}
    where := map[string]interface{}{"id": id}
    newContreoller.GetModel("users").Update(false, -1, params, where)
    sessions.SetSession(w, map[string]interface{}{"sid": sid})
    http.Redirect(w, r, "/usercontroller/showcabinet", 200)
}

type VirtController interface {
    GetModel(tableName string) models.EntityInterface
    Render(filename string, data interface{})
    CheckSid() (id int, result bool)
    isAdmin() bool
    regExists(userId, eventId int) bool
}
