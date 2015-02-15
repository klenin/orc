package models

type FormsTypes struct {
    Id           string `name:"id" type:"int" null:"NOT NULL" extra:"PRIMARY"`
    FormId       string `name:"form_id" type:"int" null:"NOT NULL" extra:"REFERENCES" refTable:"forms" refField:"id" refFieldShow:"name"`
    TypeId       string `name:"type_id" type:"int" null:"NOT NULL" extra:"REFERENCES" refTable:"event_types" refField:"id" refFieldShow:"name"`
    SerialNumber string `name:"serial_number" type:"int" null:"NOT NULL" extra:"UNIQUE"`
}

func (c *ModelManager) FormsTypes() *FormsTypesModel {
    model := new(FormsTypesModel)

    model.TableName = "forms_types"
    model.Caption = "Формы - Типы мероприятий"

    model.Columns = []string{"id", "form_id", "type_id", "serial_number"}
    model.ColNames = []string{"ID", "Форма", "Тип", "Порядковый номер"}

    model.Fields = new(FormsTypes)

    model.Sub = false
    model.SubTable = nil
    model.SubField = ""

    return model
}

type FormsTypesModel struct {
    Entity
}
