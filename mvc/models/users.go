package models

type UsersModel struct {
    Entity
}

type User struct {
    Id      int    `name:"id" type:"int" null:"NOT NULL" extra:"PRIMARY"`
    Login   string `name:"login" type:"text" null:"NOT NULL" extra:"UNIQUE"`
    Pass    string `name:"pass" type:"text" null:"NOT NULL" extra:""`
    Salt    string `name:"salt" type:"text" null:"NOT NULL" extra:""`
    Role    string `name:"role" type:"text" null:"NOT NULL" extra:""`
    Sid     string `name:"sid" type:"text" null:"NULL" extra:""`
    Token   string `name:"token" type:"text" null:"NULL" extra:""`
    Enabled bool   `name:"enabled" type:"boolean" null:"NULL" extra:""`
}

func (c *ModelManager) Users() *UsersModel {
    model := new(UsersModel)

    model.TableName = "users"
    model.Caption = "Пользователи"

    model.Columns = []string{"id", "login", "role", "enabled"}
    model.ColNames = []string{"ID", "Логин", "Роль", "Состояние"}

    model.Fields = new(User)
    model.WherePart = make(map[string]interface{}, 0)
    model.Condition = AND
    model.OrderBy = "id"
    model.Limit = "ALL"
    model.Offset = 0

    model.Sub = false
    model.SubTable = nil
    model.SubField = ""

    return model
}

func (this *UsersModel) GetColModel(isAdmin bool, userId int) []map[string]interface{} {
    return []map[string]interface{} {
        0: map[string]interface{} {
            "index": "id",
            "name": "id",
            "editable": false,
        },
        1: map[string]interface{} {
            "index": "login",
            "name": "login",
            "editable": true,
        },
        2: map[string]interface{} {
            "index": "role",
            "name": "role",
            "editable": true,
            "editrules": map[string]interface{}{"required": true},
        },
        3: map[string]interface{} {
            "index": "enabled",
            "name": "enabled",
            "editable": true,
            "editrules": map[string]interface{}{"required": true},
            "formatter": "checkbox",
            "formatoptions": map[string]interface{}{"disabled": true},
            "edittype": "checkbox",
            "editoptions": map[string]interface{}{"value": "true:false"},
        },
    }
}
