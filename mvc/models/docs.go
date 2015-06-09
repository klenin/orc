package models

type DocsModel struct {
    Entity
}

type Docs struct {
    Id   int    `name:"id" type:"int" null:"NOT NULL" extra:"PRIMARY"`
    Name string `name:"name" type:"text" null:"NOT NULL" extra:"UNIQUE"`
}

func (c *ModelManager) Docs() *DocsModel {
    model := new(DocsModel)

    model.TableName = "docs"
    model.Caption = "Документы"

    model.Columns = []string{"id", "name"}
    model.ColNames = []string{"ID", "Название"}

    model.Fields = new(Docs)
    model.WherePart = make(map[string]interface{}, 0)
    model.Condition = AND
    model.OrderBy = "id"
    model.Limit = "ALL"
    model.Offset = 0

    return model
}

func (this *DocsModel) GetColModel(bool, int) []map[string]interface{} {
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
            "editrules": map[string]interface{}{"required": true, "edithidden": true},
            "edittype": "file",
            "sortable": false,
            "search": true,
            "editoptions": map[string]interface{}{"enctype": "multipart/form-data"},
        },
    }
}
