package models

func (c *ModelManager) Events() *EventModel {
	model := new(EventModel)

	model.TableName = "events"
	model.Caption = "Мероприятия"

	model.Columns = []string{"id", "name", "date_start", "date_end", "time", "url"}
	model.ColNames = []string{"ID", "Название", "Дата начала", "Дата окончания", "Время", "Сайт"}

	tmp := map[string]*Field{
		"id":         {"id", "ID", "int", false},
		"name":       {"name", "Название", "text", false},
		"date_start": {"date", "Дата начала", "date", false},
		"date_end":   {"date", "Дата окончания", "date", false},
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

//func (this EventModel) GetColModel() map[string]interface{} {
//	result := make([string]interface{}, len(this.Columns))
//	result["id"] = make()
//	result["name"] = make()
//	result["date_start"] = make()
//	result["date_end"] = make()
//	result["time"] = make()
//	result["url"] = make()
//}
