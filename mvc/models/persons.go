package models

type Person struct {
    Id        int    `name:"id" type:"int" null:"NOT NULL" extra:"PRIMARY"`
    FirstName string `name:"fname" type:"text" null:"NOT NULL" extra:""`
    LastName  string `name:"lname" type:"text" null:"NOT NULL" extra:""`
}

func (c *ModelManager) Persons() *PersonsModel {
    model := new(PersonsModel)

    model.TableName = "persons"
    model.Caption = "Персоны"

    model.Columns = []string{"id", "fname", "lname"}
    model.ColNames = []string{"ID", "Имя", "Фамилия"}

    model.Fields = new(Person)
    model.WherePart = make(map[string]interface{}, 0)
    model.OrderBy = "id"
    model.Limit = "ALL"
    model.Offset = 0

    model.Sub = true
    model.SubTable = []string{"teams_persons"}
    model.SubField = "person_id"

    return model
}

type PersonsModel struct {
    Entity
}
