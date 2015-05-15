package models

import "github.com/orc/db"

type RegParamValsModel struct {
    Entity
}

type RegParamVal struct {
    Id         int `name:"id" type:"int" null:"NOT NULL" extra:"PRIMARY"`
    RegId      int `name:"reg_id" type:"int" null:"NOT NULL" extra:"REFERENCES" refTable:"registrations" refField:"id" refFieldShow:"id"`
    ParamValId int `name:"param_val_id" type:"int" null:"NOT NULL" extra:"REFERENCES" refTable:"param_values" refField:"id" refFieldShow:"id"`
}

func (c *ModelManager) RegParamVals() *RegParamValsModel {
    model := new(RegParamValsModel)

    model.TableName = "reg_param_vals"
    model.Caption = "Регистрация - Значение параметра"

    model.Columns = []string{"id", "reg_id", "param_val_id"}
    model.ColNames = []string{"ID", "Регистрация", "Значения параметра"}

    model.Fields = new(RegParamVal)
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

func (this *RegParamValsModel) GetColModel() []map[string]interface{} {
    query := `SELECT array_to_string(
        array(SELECT registrations.id || ':' || registrations.id FROM registrations GROUP BY registrations.id ORDER BY registrations.id), ';') as name;`
    regs := db.Query(query, nil)[0].(map[string]interface{})["name"].(string)

    query = `SELECT array_to_string(
        array(SELECT param_values.id || ':' || param_values.id FROM param_values GROUP BY param_values.id ORDER BY param_values.id), ';') as name;`
    param_vals := db.Query(query, nil)[0].(map[string]interface{})["name"].(string)

    return []map[string]interface{} {
        0: map[string]interface{} {
            "index": "id",
            "name": "id",
            "editable": false,
        },
        1: map[string]interface{} {
            "index": "reg_id",
            "name": "reg_id",
            "editable": true,
            "formatter": "select",
            "edittype": "select",
            "stype": "select",
            "search": true,
            "editrules": map[string]interface{}{"required": true},
            "editoptions": map[string]string{"value": regs},
            "searchoptions": map[string]string{"value": ":Все;"+regs},
        },
        2: map[string]interface{} {
            "index": "param_val_id",
            "name": "param_val_id",
            "editable": true,
            "formatter": "select",
            "edittype": "select",
            "stype": "select",
            "search": true,
            "editrules": map[string]interface{}{"required": true},
            "editoptions": map[string]string{"value": param_vals},
            "searchoptions": map[string]string{"value": ":Все;"+param_vals},
        },
    }
}
