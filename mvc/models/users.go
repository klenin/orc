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





func (*ModelManager) Users() *UsersModel {
    model := new(UsersModel)
    model.SetTableName("users").
        SetCaption("Пользователи").
        SetColumns([]string{"id", "login", "role", "enabled"}).
        SetColNames([]string{"ID", "Логин", "Роль", "Состояние"}).
        SetFields(new(User)).
        SetCondition(AND).
        SetOrder("id").
        SetLimit("ALL").
        SetOffset(0).
        SetSorting("ASC").
        SetWherePart(make(map[string]interface{}, 0)).
        SetSub(false).
        SetSubTables(nil).
        SetSubField("")

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
