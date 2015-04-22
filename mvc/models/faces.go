package models

type Face struct {
    Id       int `name:"id" type:"int" null:"NOT NULL" extra:"PRIMARY"`
    UserId   int `name:"user_id" type:"int" null:"NOT NULL" extra:"REFERENCES" refTable:"users" refField:"id" refFieldShow:"login"`
}

func (c *ModelManager) Faces() *FaceModel {
    model := new(FaceModel)

    model.TableName = "faces"
    model.Caption = "Лица"

    model.Columns = []string{"id", "user_id"}
    model.ColNames = []string{"ID", "Пользователь"}

    model.Fields = new(Face)
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

type FaceModel struct {
    Entity
}
