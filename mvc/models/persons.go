package models

type Person struct {
    Id      int    `name:"id" type:"int" null:"NOT NULL" extra:"PRIMARY"`
    FaceId  int    `name:"face_id" type:"int" null:"NULL" extra:"REFERENCES" refTable:"faces" refField:"id" refFieldShow:"id"`
    GroupId int    `name:"group_id" type:"int" null:"NOT NULL" extra:"REFERENCES" refTable:"groups" refField:"id" refFieldShow:"name"`
    Name    string `name:"name" type:"text" null:"NOT NULL" extra:""`
    Token   string `name:"token" type:"text" null:"NOT NULL" extra:""`
    Email   string `name:"email" type:"text" null:"NOT NULL" extra:""`
    Status  bool   `name:"status" type:"boolean" null:"NULL" extra:""`
}

func (c *ModelManager) Persons() *PersonsModel {
    model := new(PersonsModel)

    model.TableName = "persons"
    model.Caption = "Персоны"

    model.Columns = []string{"id", "name", "email", "group_id", "status", "face_id"}
    model.ColNames = []string{"ID", "ФИО", "Почта", "Группа", "Статус", "Лицо"}

    model.Fields = new(Person)
    model.WherePart = make(map[string]interface{}, 0)
    model.Condition = AND
    model.OrderBy = "id"
    model.Limit = "ALL"
    model.Offset = 0

    model.Sub = false
    model.SubTable = nil
    model.SubField = ""

    return model
}

type PersonsModel struct {
    Entity
}
