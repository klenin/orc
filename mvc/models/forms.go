package models

func (c *ModelManager) Forms() *FormsModel {
	model := new(FormsModel)

	model.TableName = "forms"
	model.Caption = "Формы"

	model.Columns = []string{"id", "name"}
	model.ColNames = []string{"ID", "Название"}

	tmp := map[string]*Field{
		"id":   {"id", "ID", "int", false},
		"name": {"fname", "Название", "text", false},
	}

	model.Fields = tmp

	model.Ref = false
	model.RefData = nil
	model.RefFields = nil

	model.Sub = true
	model.SubTable = []string{"forms_types"}
	model.SubField = "form_id"

	return model
}

type FormsModel struct {
	Entity
}
