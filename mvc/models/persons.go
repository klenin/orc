package models

func (c *ModelManager) Persons() *PersonsModel {
	model := new(PersonsModel)

	model.TableName = "persons"
	model.Caption = "Персоны"

	model.Columns = []string{"id", "fname", "lname", "pname"}
	model.ColNames = []string{"ID", "Фамилия", "Имя", "Отчество"}

	tmp := map[string]*Field{
		"id":    {"id", "ID", "int", false},
		"fname": {"fname", "Фамилия", "text", false},
		"lname": {"lname", "Имя", "text", false},
		"pname": {"pname", "Отчество", "text", false},
	}

	model.Fields = tmp

	model.Ref = false
	model.RefData = nil
	model.RefFields = nil

	model.Sub = true
	model.SubTable = []string{"teams_persons"}
	model.SubField = "person_id"

	return model
}

type PersonsModel struct {
	Entity
}
