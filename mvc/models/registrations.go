package models

type Registration struct {
    Id     int `name:"id" type:"int" null:"NOT NULL" extra:"PRIMARY"`
    FaceId int `name:"face_id" type:"int" null:"NOT NULL" extra:"REFERENCES" refTable:"faces" refField:"id" refFieldShow:"id"`
}

func (c *ModelManager) Registrations() *RegistrationModel {
    model := new(RegistrationModel)

    model.TableName = "registrations"
    model.Caption = "Регистрации"

    model.Columns = []string{"id", "face_id", "param_values_id"}
    model.ColNames = []string{"ID", "Лицо", "Значение"}

    model.Fields = new(Registration)
    model.WherePart = make(map[string]interface{}, 0)
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
