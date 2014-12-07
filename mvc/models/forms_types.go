package models

import (
	"github.com/orc/db"
)

func (c *ModelManager) FormsTypes() *FormsTypesModel {
	model := new(FormsTypesModel)

	model.TableName = "forms_types"
	model.Caption = "Формы - Типы мероприятий"

	model.Columns = []string{"id", "form_id", "type_id", "serial_number"}
	model.ColNames = []string{"ID", "Форма", "Тип", "Порядковый номер"}

	tmp := map[string]*Field{
		"id":            {"id", "ID", "int", false},
		"form_id":       {"form_id", "Форма", "int", true},
		"type_id":       {"type_id", "Тип", "int", true},
		"serial_number": {"serial_number", "Порядковый номер", "int", false},
	}

	model.Fields = tmp

	model.Ref = true
	model.RefFields = []string{"name"}
	model.RefData = make(map[string]interface{}, 2)

	result := db.Select("forms", nil, "", []string{"id", "name"})
	model.RefData["form_id"] = make([]interface{}, len(result))
	model.RefData["form_id"] = result

	result = db.Select("event_types", nil, "", []string{"id", "name"})
	model.RefData["type_id"] = make([]interface{}, len(result))
	model.RefData["type_id"] = result

	model.Sub = false
	model.SubTable = nil
	model.SubField = ""

	return model
}

type FormsTypesModel struct {
	Entity
}
