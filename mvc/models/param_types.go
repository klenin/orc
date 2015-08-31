package models

type ParamType struct {
    id   int    `name:"id" type:"int" null:"NOT NULL" extra:"PRIMARY"`
    name string `name:"name" type:"text" null:"NOT NULL" extra:"UNIQUE"`
}

func (this *ParamType) GetId() int {
    return this.id
}

func (this *ParamType) SetName(name string) {
    this.name = name
}

func (this *ParamType) GetName() string {
    return this.name
}

type ParamTypesModel struct {
    Entity
}

func (*ModelManager) ParamTypes() *ParamTypesModel {
    model := new(ParamTypesModel)
    model.SetTableName("param_types").
        SetCaption("Типы параметров").
        SetColumns([]string{"id", "name"}).
        SetColNames([]string{"ID", "Название"}).
        SetFields(new(ParamType)).
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

func (*ParamTypesModel) GetColModel(isAdmin bool, userId int) []map[string]interface{} {
    return []map[string]interface{} {
        0: map[string]interface{} {
            "index": "id",
            "name": "id",
            "editable": false,
        },
        1: map[string]interface{} {
            "index": "name",
            "name": "name",
            "editable": true,
            "editrules": map[string]interface{}{"required": true},
        },
    }
}
