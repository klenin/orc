package models

func (c *ModelManager) ParamTypes() *ParamTypeModel {
    model := new(ParamTypeModel)

    model.TableName = "param_types"
    model.Caption = "Типы параметров"

    model.Columns = []string{"id", "name"}
    model.ColNames = []string{"ID", "Название"}

    model.Fields = []map[string]string{
        {
            "field": "id",
            "type":  "int",
            "null":  "NOT NULL",
            "extra": "PRIMARY"},
        {
            "field": "name",
            "type":  "text",
            "null":  "NOT NULL",
            "extra": "UNIQUE"},
    }

    model.Ref = false
    model.RefData = nil
    model.RefFields = nil

    model.Sub = false
    model.SubTable = nil
    model.SubField = ""

    return model
}

type ParamTypeModel struct {
    Entity
}
