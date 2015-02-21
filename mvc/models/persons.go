package models

type Person struct {
    Id int `name:"id" type:"int" null:"NOT NULL" extra:"PRIMARY"`
}

func (c *ModelManager) Persons() *PersonsModel {
    model := new(PersonsModel)

    model.TableName = "persons"
    model.Caption = "Персоны"

    model.Columns = []string{"id"}
    model.ColNames = []string{"ID"}

    model.Fields = new(Person)
    model.WherePart = make(map[string]interface{}, 0)
    model.OrderBy = "id"
    model.Limit = "ALL"
    model.Offset = 0

    model.Sub = false
    model.SubTable = nil
    model.SubField = ""

    return model
}

type PersonsModel struct {
    Entity
}
