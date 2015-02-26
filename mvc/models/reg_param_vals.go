package models

type RegParamVal struct {
    Id          int `name:"id" type:"int" null:"NOT NULL" extra:"PRIMARY"`
    RegId       int `name:"reg_id" type:"int" null:"NOT NULL" extra:"REFERENCES" refTable:"registrations" refField:"id" refFieldShow:"id"`
    EventId     int `name:"event_id" type:"int" null:"NOT NULL" extra:"REFERENCES" refTable:"events" refField:"id" refFieldShow:"name"`
    EventTypeId int `name:"event_type_id" type:"int" null:"NOT NULL" extra:"REFERENCES" refTable:"event_types" refField:"id" refFieldShow:"name"`
    ParamValId  int `name:"param_val_id" type:"int" null:"NOT NULL" extra:"REFERENCES" refTable:"param_values" refField:"id" refFieldShow:"id"`
}

func (c *ModelManager) RegParamVals() *RegParamValsModel {
    model := new(RegParamValsModel)

    model.TableName = "reg_param_vals"
    model.Caption = "Регистрация-Мероприятие-Тип мероприятия-Значение параметра"

    model.Columns = []string{"id", "reg_id", "event_id", "event_type_id", "param_val_id"}
    model.ColNames = []string{"ID", "Регистрация", "Мероприятие", "Тип мероприятия", "Значения параметра"}

    model.Fields = new(RegParamVal)
    model.WherePart = make(map[string]interface{}, 0)
    model.OrderBy = "id"
    model.Limit = "ALL"
    model.Offset = 0

    model.Sub = false
    model.SubTable = nil
    model.SubField = ""

    return model
}

type RegParamValsModel struct {
    Entity
}
