package models

func (c *ModelManager) EventTypes() *EventTypesModel {
	model := new(EventTypesModel)

	model.TableName = "event_types"
	model.Caption = "Типы мероприятий"

	model.Columns = []string{"id", "name", "description", "topicality"}
	model.ColNames = []string{"ID", "Тип", "Описание", "Актуальность"}

	tmp := map[string]*Field{
		"id":          {"id", "ID", "int", false},
		"name":        {"name", "Тип", "text", false},
		"description": {"description", "Описание", "text", false},
		"topicality":  {"topicality", "Актуальность", "boolean", false},
	}

	model.Fields = tmp

	model.Ref = false
	model.RefData = nil
	model.RefFields = nil

	model.Sub = true
	model.SubTable = []string{"forms_types", "events_types"}
	model.SubField = "type_id"

	return model
}

type EventTypesModel struct {
	Entity
}
