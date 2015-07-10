package models

type FormsModel struct {
    Entity
}

type Forms struct {
    Id       int    `name:"id" type:"int" null:"NOT NULL" extra:"PRIMARY"`
    Name     string `name:"name" type:"text" null:"NOT NULL" extra:"UNIQUE"`
    Personal bool   `name:"personal" type:"boolean" null:"NOT NULL" extra:""`
}

func (c *ModelManager) Forms() *FormsModel {
    model := new(FormsModel)

    model.TableName = "forms"
    model.Caption = "Формы"

    model.Columns = []string{"id", "name", "personal"}
    model.ColNames = []string{"ID", "Название", "Персональная"}

    model.Fields = new(Forms)
    model.WherePart = make(map[string]interface{}, 0)
    model.Condition = AND
    model.OrderBy = "id"
    model.Limit = "ALL"
    model.Offset = 0

    // model.Sub = true
    // model.SubTable = []string{"events_forms"}
    // model.SubField = "form_id"

    return model
}

func (this *FormsModel) GetColModel(isAdmin bool, userId int) []map[string]interface{} {
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
