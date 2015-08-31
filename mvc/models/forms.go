package models

type Form struct {
    id       int    `name:"id" type:"int" null:"NOT NULL" extra:"PRIMARY"`
    name     string `name:"name" type:"text" null:"NOT NULL" extra:"UNIQUE"`
    personal bool   `name:"personal" type:"boolean" null:"NOT NULL" extra:""`
}

func (this *Form) GetId() int {
    return this.id
}

func (this *Form) SetName(name string) {
    this.name = name
}

func (this *Form) GetName() string {
    return this.name
}

func (this *Form) SetPersonal(personal bool) {
    this.personal = personal
}

func (this *Form) GetPersonal() bool {
    return this.personal
}

type FormsModel struct {
    Entity
}

func (*ModelManager) Forms() *FormsModel {
    model := new(FormsModel)
    model.SetTableName("forms").
        SetCaption("Формы").
        SetColumns([]string{"id", "name", "personal"}).
        SetColNames([]string{"ID", "Название", "Персональная"}).
        SetFields(new(Form)).
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

func (*FormsModel) GetColModel(isAdmin bool, userId int) []map[string]interface{} {
    return []map[string]interface{} {
        0: map[string]interface{} {
            "index": "id",
            "name": "id",
            "editable": false,
            "width": "20",
        },
        1: map[string]interface{} {
            "index": "name",
            "name": "name",
            "editable": true,
            "editrules": map[string]interface{}{"required": true},
        },
        2: map[string]interface{} {
            "index": "personal",
            "name": "personal",
            "editable": true,
            "editrules": map[string]interface{}{"required": true},
            "formatter": "checkbox",
            "formatoptions": map[string]interface{}{"disabled": true},
            "edittype": "checkbox",
            "editoptions": map[string]interface{}{"value": "true:false"},
        },
    }
}
