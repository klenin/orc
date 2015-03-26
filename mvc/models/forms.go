package models

type Forms struct {
    Id   int    `name:"id" type:"int" null:"NOT NULL" extra:"PRIMARY"`
    Name string `name:"name" type:"text" null:"NOT NULL" extra:"UNIQUE"`
}

func (c *ModelManager) Forms() *FormsModel {
    model := new(FormsModel)

    model.TableName = "forms"
    model.Caption = "Формы"

    model.Columns = []string{"id", "name"}
    model.ColNames = []string{"ID", "Название"}

    model.Fields = new(Forms)
    model.WherePart = make(map[string]interface{}, 0)
    model.Condition = AND
    model.OrderBy = "id"
    model.Limit = "ALL"
    model.Offset = 0

    model.Sub = true
    model.SubTable = []string{"events_forms"}
    model.SubField = "form_id"

    return model
}

type FormsModel struct {
    Entity
}
