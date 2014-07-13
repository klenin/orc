package models

import (
	"github.com/orc/db"
)

func (c *ModelManager) ParamValues() *ParamValuesModel {
	model := new(ParamValuesModel)

	model.TableName = "param_values"
	model.Caption = "Знвчение параметров"

	model.Columns = []string{"id", "person_id", "event_id", "param_id", "value"}
	model.ColNames = []string{"ID", "Персона", "Мероприятие", "Параметр", "Значение"}

	tmp := map[string]*Field{
		"id":        {"id", "ID", "int", false},
		"person_id": {"person_id", "Персона", "int", true},
		"event_id":  {"event_id", "Мероприятие", "int", true},
		"param_id":  {"param_id", "Параметр", "int", true},
		"value":     {"value", "Значение", "text", true},
	}

	model.Fields = tmp

	model.Ref = true
	model.RefFields = []string{"fname", "name"}
	model.RefData = make(map[string]interface{}, 3)

	result := db.Select("persons", nil, "", []string{"id", "fname"})
	model.RefData["person_id"] = make([]interface{}, len(result))
	model.RefData["person_id"] = result

	result = db.Select("events", nil, "", []string{"id", "name"})
	model.RefData["event_id"] = make([]interface{}, len(result))
	model.RefData["event_id"] = result

	result = db.Select("params", nil, "", []string{"id", "name"})
	model.RefData["param_id"] = make([]interface{}, len(result))
	model.RefData["param_id"] = result

	model.Sub = false
	model.SubTable = nil
	model.SubField = ""

	return model
}

type ParamValuesModel struct {
	Entity
}
