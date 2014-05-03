package models

import (
	"github.com/orc/db"
)

func (c *ModelManager) TeamsPersons() *TeamsPersonsModel {
	model := new(TeamsPersonsModel)

	model.TableName = "teams_persons"
	model.Caption = "Команды - Персоны"

	model.Columns = []string{"id", "team_id", "person_id"}
	model.ColNames = []string{"ID", "Команда", "Персона"}

	tmp := map[string]*Field{
		"id":        {"id", "ID", "int", false},
		"team_id":   {"team_id", "Команда", "int", true},
		"person_id": {"person_id", "Персона", "int", true},
	}

	model.Fields = tmp

	model.Ref = true
	model.RefFields = []string{"name", "fname"}
	model.RefData = make(map[string]interface{}, 2)

	result := db.Select("teams", nil, []string{"id", "name"})
	model.RefData["team_id"] = make([]interface{}, len(result))
	model.RefData["team_id"] = result

	result = db.Select("persons", nil, []string{"id", "fname"})
	model.RefData["person_id"] = make([]interface{}, len(result))
	model.RefData["person_id"] = result

	model.Sub = false
	model.SubTable = nil
	model.SubField = ""

	return model
}

type TeamsPersonsModel struct {
	Entity
}
