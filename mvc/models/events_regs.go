package models

type EventReg struct {
    Id      int `name:"id" type:"int" null:"NOT NULL" extra:"PRIMARY"`
    EventId int `name:"event_id" type:"int" null:"NOT NULL" extra:"REFERENCES" refTable:"events" refField:"id" refFieldShow:"name"`
    RegId   int `name:"reg_id" type:"int" null:"NULL" extra:"REFERENCES" refTable:"registrations" refField:"id" refFieldShow:"id"`
}

func (c *ModelManager) EventsRegs() *EventRegModel {
    model := new(EventRegModel)

    model.TableName = "events_regs"
    model.Caption = "Мероприятия-Регистрации"

    model.Columns = []string{"id", "event_id", "reg_id"/*, "reg_date"*/}
    model.ColNames = []string{"ID", "Мероприятие", "Регистрация"/*, "Дата регистрации"*/}

    model.Fields = new(EventReg)
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

type EventRegModel struct {
    Entity
}
