package models

type Param struct {
    Id          string `name:"id" type:"int" null:"NOT NULL" extra:"PRIMARY"`
    Name        string `name:"name" type:"text" null:"NOT NULL" extra:"UNIQUE"`
    FormId      string `name:"form_id" type:"int" null:"NOT NULL" extra:"REFERENCES" refTable:"forms" refField:"id" refFieldShow:"name"`
    ParamTypeId string `name:"param_type_id" type:"int" null:"NOT NULL" extra:"REFERENCES" refTable:"param_types" refField:"id" refFieldShow:"name"`
    Identifier  string `name:"identifier" type:"int" null:"NOT NULL" extra:"UNIQUE"`
}

func (c *ModelManager) Params() *ParamsModel {
    model := new(ParamsModel)

    model.TableName = "params"
    model.Caption = "Параметры"

    model.Columns = []string{"id", "name", "param_type_id", "form_id", "identifier"}
    model.ColNames = []string{"ID", "Название", "Тип", "Форма", "Идентификатор"}

    model.Fields = new(Param)

    model.Sub = false
    model.SubTable = nil
    model.SubField = ""

    return model
}

type ParamsModel struct {
    Entity
}
