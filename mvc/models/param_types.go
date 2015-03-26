package models

type ParamTypes struct {
    Id   int    `name:"id" type:"int" null:"NOT NULL" extra:"PRIMARY"`
    Name string `name:"name" type:"text" null:"NOT NULL" extra:"UNIQUE"`
}

func (c *ModelManager) ParamTypes() *ParamTypesModel {
    model := new(ParamTypesModel)

    model.TableName = "param_types"
    model.Caption = "Типы параметров"

    model.Columns = []string{"id", "name"}
    model.ColNames = []string{"ID", "Название"}

    model.Fields = new(ParamTypes)
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

type ParamTypesModel struct {
    Entity
}
