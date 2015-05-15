package controllers

import (
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
    Id        string
    TableName string
    Caption   string
    Table     []interface{}
    RefData   map[string]interface{}
    RefFields []string
    Columns   []string
    ColNames  []string
    ColModel  []map[string]interface{}
    Sub       bool
    SubTableName string
    SubCaption   string
    SubRefData   map[string]interface{}
    SubRefFields []string
    SubColumns   []string
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
    var role string

    user_id := sessions.GetValue("id", this.Request)
    if user_id == nil {
        http.Redirect(this.Response, this.Request, "/", http.StatusUnauthorized)
        return false
    }

    user := this.GetModel("users")
    user.LoadWherePart(map[string]interface{}{"id": user_id})
    err := db.SelectRow(user, []string{"role"}).Scan(&role)
    if err != nil || role == "user" {
        return false
    }

    return role == "admin"
}

type VirtController interface {
    GetModel(tableName string) models.VirtEntity
    Render(filename string, data interface{})
    isAdmin() bool
}
