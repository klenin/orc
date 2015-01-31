package models

import "github.com/orc/db"

func (c *ModelManager) Users() *UsersModel {
    model := new(UsersModel)

    model.TableName = "users"
    model.Caption = "Пользователи"

    model.Columns = []string{"id", "login", "role", "person_id"}
    model.ColNames = []string{"ID", "Логин", "Роль", "Персона"}

    model.Fields = []map[string]string{
        {
            "field": "id",
            "type":  "int",
            "null":  "NOT NULL",
            "extra": "PRIMARY"},
        {
            "field": "login",
            "type":  "text",
            "null":  "NOT NULL",
            "extra": ""},
        {
            "field": "pass",
            "type":  "text",
            "null":  "NOT NULL",
            "extra": ""},
        {
            "field": "salt",
            "type":  "text",
            "null":  "NOT NULL",
            "extra": ""},
        {
            "field": "role",
            "type":  "text",
            "null":  "NOT NULL",
            "extra": ""},
        {
            "field":    "person_id",
            "type":     "int",
            "null":     "NOT NULL",
            "extra":    "REFERENCES",
            "refTable": "persons",
            "refField": "id"},
    }

    model.Ref = true
    model.RefFields = []string{"fname"}
    model.RefData = make(map[string]interface{}, 1)

    result := db.Select("persons", nil, "", []string{"id", "fname"})
    model.RefData["person_id"] = make([]interface{}, len(result))
    model.RefData["person_id"] = result

    model.Sub = false
    model.SubTable = nil
    model.SubField = ""

    return model
}

type UsersModel struct {
    Entity
}
