package controllers

import (
    "errors"
    "github.com/orc/db"
    "github.com/orc/mvc/models"
    "github.com/orc/sessions"
    "net/http"
    "html/template"
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

func (this *Controller) GetModel(tableName string) models.VirtEntity {
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

func (this *Controller) isAdmin() bool {
    userId, err := this.CheckSid()
    if err != nil {
        return false
    }

    var role string
    user := this.GetModel("users")
    user.LoadWherePart(map[string]interface{}{"id": userId})
    err = db.SelectRow(user, []string{"role"}).Scan(&role)
    if err != nil || role == "user" {
        return false
    }

    return role == "admin"
}

func (this *Controller) CheckSid() (id int, result error)  {
    userSid := sessions.GetValue("sid", this.Request)
    if !sessions.CheckSession(this.Response, this.Request) || userSid == nil {
        return -1, errors.New("Данные в куках отсутствуют.")
    }

    user := this.GetModel("users")
    user.LoadWherePart(map[string]interface{}{"sid": userSid})

    err := db.SelectRow(user, []string{"id"}).Scan(&id)
    if err != nil {
        return -1, errors.New("Данные в куках отсутствуют.")
    }

    return id, nil
}

type VirtController interface {
    GetModel(tableName string) models.VirtEntity
    Render(filename string, data interface{})
    isAdmin() bool
    CheckSid() (id int, result bool)
}
