package models

func (c *ModelManager) Events() *EventModel {
	model := new(EventModel)

	model.TableName = "events"
	model.Caption = "Мероприятия"

	model.Columns = []string{"id", "name", "date_start", "date_end", "time", "url"}
	model.ColNames = []string{"ID", "Название", "Дата_start", "Дата_end", "Время", "Сайт"}

	tmp := map[string]*Field{
		"id":         {"id", "ID", "int", false},
		"name":       {"name", "Название", "text", false},
		"date_start": {"date", "Дата_start", "date", false},
		"date_end":   {"date", "Дата_end", "date", false},
		"time":       {"time", "Время", "time", false},
		"url":        {"url", "Сайт", "text", false},
	}

	model.Fields = tmp

	model.Ref = false
	model.RefData = nil
	model.RefFields = nil

	model.Sub = true
	model.SubTable = []string{"events_types"}
	model.SubField = "event_id"
	return model
}

type EventModel struct {
	Entity
}
