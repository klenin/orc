package models

type GroupRegistration struct {
    Id      int  `name:"id" type:"int" null:"NOT NULL" extra:"PRIMARY"`
    EventId int  `name:"event_id" type:"int" null:"NOT NULL" extra:"REFERENCES" refTable:"events" refField:"id" refFieldShow:"name"`
    GroupId int  `name:"group_id" type:"int" null:"NOT NULL" extra:"REFERENCES" refTable:"groups" refField:"id" refFieldShow:"name"`
    Status  bool `name:"status" type:"boolean" null:"NULL" extra:""`
}

func (c *ModelManager) GroupRegistrations() *GroupRegistrationModel {
    model := new(GroupRegistrationModel)

    model.TableName = "group_registrations"
    model.Caption = "Регистрации групп"

    model.Columns = []string{"id", "event_id", "group_id", "status"}
    model.ColNames = []string{"ID", "Мероприятие", "Группа", "Статус"}

    model.Fields = new(GroupRegistration)
    model.WherePart = make(map[string]interface{}, 0)
    model.Condition = AND
    model.OrderBy = "id"
    model.Limit = "ALL"
    model.Offset = 0

    model.Sub = true
    model.SubTable = []string{"persons"}
    model.SubField = "group_id"


    return model
}

type GroupRegistrationModel struct {
    Entity
}
