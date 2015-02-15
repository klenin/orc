package models

type Person struct {
    Id        string `name:"id" type:"int" null:"NOT NULL" extra:"PRIMARY"`
    firstName string `name:"fname" type:"text" null:"NOT NULL" extra:""`
    lastName  string `name:"lname" type:"text" null:"NOT NULL" extra:""`
}

func (c *ModelManager) Persons() *PersonsModel {
    model := new(PersonsModel)

    model.TableName = "persons"
    model.Caption = "Персоны"

    model.Columns = []string{"id", "fname", "lname"}
    model.ColNames = []string{"ID", "Фамилия", "Имя"}

    model.Fields = new(Person)

    model.Sub = true
    model.SubTable = []string{"teams_persons"}
    model.SubField = "person_id"

    return model
}

type PersonsModel struct {
    Entity
}
