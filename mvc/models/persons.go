package models


func (c *ModelManager) Persons() *PersonsModel {
    model := new(PersonsModel)

    model.TableName = "persons"
    model.Caption = "Персоны"

    model.Columns = []string{"id", "fname", "lname",}
    model.ColNames = []string{"ID", "Фамилия", "Имя",}

    model.Fields = []map[string]string{
        {
            "field": "id",
            "type":  "int",
            "null":  "NOT NULL",
            "extra": "PRIMARY"},
        {
            "field": "fname",
            "type":  "text",
            "null":  "NOT NULL",
            "extra": ""},
        {
            "field": "lname",
            "type":  "text",
            "null":  "NOT NULL",
            "extra": ""},
    }

    model.Ref = false
    model.RefData = nil
    model.RefFields = nil

    model.Sub = true
    model.SubTable = []string{"teams_persons"}
    model.SubField = "person_id"

    return model
}

type PersonsModel struct {
    Entity
}
