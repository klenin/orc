package models

type ParamValues struct {
    Id          int    `name:"id" type:"int" null:"NOT NULL" extra:"PRIMARY"`
    ParamId     int    `name:"param_id" type:"int" null:"NOT NULL" extra:"REFERENCES" refTable:"params" refField:"id" refFieldShow:"name"`
    Value       string `name:"value" type:"text" null:"NOT NULL" extra:""`
}

func (c *ModelManager) ParamValues() *ParamValuesModel {
    model := new(ParamValuesModel)

    model.TableName = "param_values"
    model.Caption = "Значение параметров"

    model.Columns = []string{"id", "param_id", "value"}
    model.ColNames = []string{"ID", "Параметр", "Значение"}

    model.Fields = new(ParamValues)
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

type ParamValuesModel struct {
    Entity
}
