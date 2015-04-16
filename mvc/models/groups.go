package models

type Groups struct {
    Id    int    `name:"id" type:"int" null:"NOT NULL" extra:"PRIMARY"`
    Name  string `name:"name" type:"text" null:"NOT NULL" extra:"UNIQUE"`
    Owner int    `name:"face_id" type:"int" null:"NOT NULL" extra:"REFERENCES" refTable:"faces" refField:"id" refFieldShow:"id"`
}

func (c *ModelManager) Groups() *GroupsModel {
    model := new(GroupsModel)

    model.TableName = "groups"
    model.Caption = "Группы"

    model.Columns = []string{"id", "name", "face_id"}
    model.ColNames = []string{"ID", "Название", "Лицо-Владелец"}

    model.Fields = new(Groups)
    model.WherePart = make(map[string]interface{}, 0)
    model.Condition = AND
    model.OrderBy = "id"
    model.Limit = "ALL"
    model.Offset = 0

    model.Sub = true
    model.SubTable = nil
    model.SubField = ""

    return model
}

type GroupsModel struct {
    Entity
}
