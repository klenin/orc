package models

import (
	"github.com/orc/db"
)

func (c *ModelManager) PersonsEvents() *PersonsEventsModel {
	model := new(PersonsEventsModel)

	model.TableName = "persons_events"
	model.Caption = "Персоны - Мероприятия"

	model.Columns = []string{"id", "person_id", "event_id", "reg_date", "last_date"}
	model.ColNames = []string{"ID", "Персона", "Мероприятие", "Дата_регистрации", "Дата_последних_изменений"}

	tmp := map[string]*Field{
		"id":        {"id", "ID", "int", false},
		"person_id": {"person_id", "Персона", "int", true},
		"event_id":  {"event_id", "Мероприятие", "int", true},
		"reg_date":  {"reg_date", "Дата_регистрации", "date", true},
		"last_date": {"last_date", "Дата_последних_изменений", "date", true},
	}

	model.Fields = tmp

	model.Ref = true
	model.RefFields = []string{"name", "fname"}
	model.RefData = make(map[string]interface{}, 2)

	result := db.Select("persons", nil, "", []string{"id", "fname"})
	model.RefData["person_id"] = make([]interface{}, len(result))
	model.RefData["person_id"] = result

	result = db.Select("events", nil, "", []string{"id", "name"})
	model.RefData["event_id"] = make([]interface{}, len(result))
	model.RefData["event_id"] = result

	model.Sub = false
	model.SubTable = nil
	model.SubField = ""

	return model
}

type PersonsEventsModel struct {
	Entity
}
