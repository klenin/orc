package models

import (
	"github.com/orc/db"
)

func (c *ModelManager) EventsTypes() *EventsTypesModel {
	model := new(EventsTypesModel)

	model.TableName = "events_types"
	model.Caption = "Мероприятия - Типы"

	model.Columns = []string{"id", "event_id", "type_id"}
	model.ColNames = []string{"ID", "Мероприятие", "Тип"}

	tmp := map[string]*Field{
		"id":       {"id", "ID", "int", false},
		"event_id": {"event_id", "Мероприятие", "int", true},
		"type_id":  {"type_id", "Тип", "int", true},
	}

	model.Fields = tmp

	model.Ref = true
	model.RefFields = []string{"name"}
	model.RefData = make(map[string]interface{}, 2)

	result := db.Select("events", nil, "", []string{"id", "name"})
	model.RefData["event_id"] = make([]interface{}, len(result))
	model.RefData["event_id"] = result

	result = db.Select("event_types", nil, "", []string{"id", "name"})
	model.RefData["type_id"] = make([]interface{}, len(result))
	model.RefData["type_id"] = result

	model.Sub = false
	model.SubTable = nil
	model.SubField = ""

	return model
}

type EventsTypesModel struct {
	Entity
}
