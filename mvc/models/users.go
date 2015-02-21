package models

type User struct {
    Id       int    `name:"id" type:"int" null:"NOT NULL" extra:"PRIMARY"`
    Login    string `name:"login" type:"text" null:"NOT NULL" extra:""`
    Pass     string `name:"pass" type:"text" null:"NOT NULL" extra:""`
    Salt     string `name:"salt" type:"text" null:"NOT NULL" extra:""`
    Role     string `name:"role" type:"text" null:"NOT NULL" extra:""`
    Hash     string `name:"hash" type:"text" null:"NULL" extra:""`
    PersonId int    `name:"person_id" type:"int" null:"NOT NULL" extra:"REFERENCES" refTable:"persons" refField:"id" refFieldShow:"fname"`
}

func (c *ModelManager) Users() *UsersModel {
    model := new(UsersModel)

    model.TableName = "users"
    model.Caption = "Пользователи"

    model.Columns = []string{"id", "login", "role", "person_id"}
    model.ColNames = []string{"ID", "Логин", "Роль", "Персона"}

    model.Fields = new(User)
    model.WherePart = make(map[string]interface{}, 0)
    model.OrderBy = "id"
    model.Limit = "ALL"
    model.Offset = 0

    model.Sub = false
    model.SubTable = nil
    model.SubField = ""

    return model
}

type UsersModel struct {
    Entity
}
