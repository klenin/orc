package models

import "github.com/orc/db"

type Person struct {
    Id      int    `name:"id" type:"int" null:"NOT NULL" extra:"PRIMARY"`
    FaceId  int    `name:"face_id" type:"int" null:"NULL" extra:"REFERENCES" refTable:"faces" refField:"id" refFieldShow:"id"`
    GroupId int    `name:"group_id" type:"int" null:"NOT NULL" extra:"REFERENCES" refTable:"groups" refField:"id" refFieldShow:"name"`
    Name    string `name:"name" type:"text" null:"NOT NULL" extra:""`
    Token   string `name:"token" type:"text" null:"NOT NULL" extra:""`
    Email   string `name:"email" type:"text" null:"NOT NULL" extra:""`
    Status  bool   `name:"status" type:"boolean" null:"NOT NULL" extra:""`
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

func (this *PersonsModel) GetModelRefDate() (fields []string, result map[string]interface{}) {
    fields = []string{"name", "name"}

    result = make(map[string]interface{})

    result["group_id"] = db.Select(new(ModelManager).GetModel("groups"), []string{"id", "name"})

    query := `SELECT faces.id as id, array_to_string(array_agg(param_values.value), ' ') as name
        FROM reg_param_vals
        INNER JOIN registrations ON registrations.id = reg_param_vals.reg_id
        INNER JOIN faces ON faces.id = registrations.face_id
        INNER JOIN events ON events.id = registrations.event_id
        INNER JOIN param_values ON param_values.id = reg_param_vals.param_val_id
        INNER JOIN params ON params.id = param_values.param_id
        WHERE params.id in (5, 6, 7) AND events.id = 1 GROUP BY faces.id ORDER BY faces.id;`

    result["face_id"] = db.Query(query, nil)

    return fields, result
}
