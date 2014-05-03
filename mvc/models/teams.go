package models

func (c *ModelManager) Teams() *TeamsModel {
	model := new(TeamsModel)

	model.TableName = "teams"
	model.Caption = "Команды"

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
	model.SubTable = []string{"teams_persons"}
	model.SubField = "team_id"

	return model
}

type TeamsModel struct {
	Entity
}
