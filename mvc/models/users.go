package models

import "github.com/orc/db"

func (c *ModelManager) Users() *UsersModel {
	model := new(UsersModel)

	model.TableName = "users"
	model.Caption = "Пользователи"

	model.Columns = []string{"id", "login", "role", "person_id"}
	model.ColNames = []string{"ID", "Логин", "Роль", "Персона"}

	tmp := map[string]*Field{
		"id":        {"id", "ID", "int", false},
		"login":     {"login", "Логин", "text", false},
		"pass":      {"password", "Хеш", "text", false},
		"salt":      {"salt", "Соль", "text", false},
		"role":      {"role", "Роль", "text", false},
		"person_id": {"person_id", "Персона", "text", true},
	}

	model.Fields = tmp

	model.Ref = true
	model.RefFields = []string{"fname"}
	model.RefData = make(map[string]interface{}, 1)

	result := db.Select("persons", nil, "", []string{"id", "fname"})
	model.RefData["person_id"] = make([]interface{}, len(result))
	model.RefData["person_id"] = result

	model.Sub = false
	model.SubTable = nil
	model.SubField = ""

	return model
}

type UsersModel struct {
	Entity
}
