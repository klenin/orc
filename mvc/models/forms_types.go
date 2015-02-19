package models

type FormsTypes struct {
    Id           int `name:"id" type:"int" null:"NOT NULL" extra:"PRIMARY"`
    FormId       int `name:"form_id" type:"int" null:"NOT NULL" extra:"REFERENCES" refTable:"forms" refField:"id" refFieldShow:"name"`
    TypeId       int `name:"type_id" type:"int" null:"NOT NULL" extra:"REFERENCES" refTable:"event_types" refField:"id" refFieldShow:"name"`
    SerialNumber int `name:"serial_number" type:"int" null:"NOT NULL" extra:"UNIQUE"`
}

func (c *ModelManager) FormsTypes() *FormsTypesModel {
    model := new(FormsTypesModel)

    model.TableName = "forms_types"
    model.Caption = "Формы - Типы мероприятий"

    model.Columns = []string{"id", "form_id", "type_id", "serial_number"}
    model.ColNames = []string{"ID", "Форма", "Тип", "Порядковый номер"}

    model.Fields = new(FormsTypes)
    model.WherePart = make(map[string]interface{}, 0)

    model.Sub = false
    model.SubTable = nil
    model.SubField = ""

    return model
}

type FormsTypesModel struct {
    Entity
}
