package models

import "github.com/orc/db"

type Registration struct {
    Id      int  `name:"id" type:"int" null:"NOT NULL" extra:"PRIMARY"`
    FaceId  int  `name:"face_id" type:"int" null:"NOT NULL" extra:"REFERENCES" refTable:"faces" refField:"id" refFieldShow:"id"`
    EventId int  `name:"event_id" type:"int" null:"NOT NULL" extra:"REFERENCES" refTable:"events" refField:"id" refFieldShow:"name"`
    Status  bool `name:"status" type:"boolean" null:"NOT NULL" extra:""`
}

func (c *ModelManager) Registrations() *RegistrationModel {
    model := new(RegistrationModel)

    model.TableName = "registrations"
    model.Caption = "Регистрации"

    model.Columns = []string{"id", "face_id", "event_id", "status"}
    model.ColNames = []string{"ID", "Лицо", "Мероприятие", "Статус"}

    model.Fields = new(Registration)
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

type RegistrationModel struct {
    Entity
}

func (this *RegistrationModel) GetModelRefDate() (fields []string, result map[string]interface{}) {
    fields = []string{"name", "name"}

    result = make(map[string]interface{})

    result["event_id"] = db.Select(new(ModelManager).GetModel("events"), []string{"id", "name"})

    query := `SELECT faces.id as id, array_to_string(array_agg(param_values.value), ' ') as name
        FROM reg_param_vals
        INNER JOIN registrations ON registrations.id = reg_param_vals.reg_id
        INNER JOIN faces ON faces.id = registrations.face_id
        INNER JOIN events ON events.id = registrations.event_id
        INNER JOIN param_values ON param_values.id = reg_param_vals.param_val_id
        INNER JOIN params ON params.id = param_values.param_id
        WHERE params.id in (5, 6, 7) GROUP BY faces.id ORDER BY faces.id;`

    result["face_id"] = db.Query(query, nil)

    return fields, result
}
