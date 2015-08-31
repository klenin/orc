package models

type User struct {
    id      int    `name:"id" type:"int" null:"NOT NULL" extra:"PRIMARY"`
    login   string `name:"login" type:"text" null:"NOT NULL" extra:"UNIQUE"`
    pass    string `name:"pass" type:"text" null:"NOT NULL" extra:""`
    salt    string `name:"salt" type:"text" null:"NOT NULL" extra:""`
    role    string `name:"role" type:"text" null:"NOT NULL" extra:""`
    sid     string `name:"sid" type:"text" null:"NULL" extra:""`
    token   string `name:"token" type:"text" null:"NULL" extra:""`
    enabled bool   `name:"enabled" type:"boolean" null:"NULL" extra:""`
}

func (this *User) GetId() int {
    return this.id
}

func (this *User) SetEnabled(enabled bool) {
    this.enabled = enabled
}

func (this *User) GetEnabled() bool {
    return this.enabled
}

func (this *User) SetSid(sid string) {
    this.sid = sid
}

func (this *User) GetSid() string {
    return this.sid
}

func (this *User) SetRole(role string) {
    this.role = role
}

func (this *User) GetRole() string {
    return this.role
}

func (this *User) SetSalt(salt string) {
    this.salt = salt
}

func (this *User) GetSalt() string {
    return this.salt
}

func (this *User) SetPass(pass string) {
    this.pass = pass
}

func (this *User) GetPass() string {
    return this.pass
}

func (this *User) SetLogin(login string) {
    this.login = login
}

func (this *User) GetLogin() string {
    return this.login
}

func (this *User) SetToken(token string) {
    this.token = token
}

func (this *User) GetToken() string {
    return this.token
}

type UsersModel struct {
    Entity
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

func (*UsersModel) GetColModel(isAdmin bool, userId int) []map[string]interface{} {
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
