package models

func (c *ModelManager) Forms() *FormsModel {
    model := new(FormsModel)

    model.TableName = "forms"
    model.Caption = "Формы"

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

    model.Sub = true
    model.SubTable = []string{"forms_types"}
    model.SubField = "form_id"

    return model
}

type FormsModel struct {
    Entity
}
