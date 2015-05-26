package models

import (
    "github.com/orc/db"
    // "strconv"
)

type RegsGroupRegsModel struct {
    Entity
}

type RegsGroupRegs struct {
    Id         int `name:"id" type:"int" null:"NOT NULL" extra:"PRIMARY"`
    GroupRegId int `name:"groupreg_id" type:"int" null:"NOT NULL" extra:"REFERENCES" refTable:"group_registrations" refField:"id" refFieldShow:"id"`
    RegId      int `name:"reg_id" type:"int" null:"NOT NULL" extra:"REFERENCES" refTable:"registrations" refField:"id" refFieldShow:"id"`
}

func (c *ModelManager) RegsGroupRegs() *RegsGroupRegsModel {
    model := new(RegsGroupRegsModel)

    model.TableName = "regs_groupregs"
    model.Caption = "Регистрации групп - Регистрации"

    model.Columns = []string{"id", "groupreg_id", "reg_id"}
    model.ColNames = []string{"ID", "Регистрации групп", "Регистрации"}

    model.Fields = new(RegsGroupRegs)
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

func (this *RegsGroupRegsModel) GetColModel() []map[string]interface{} {
    query := `SELECT array_to_string(
        array(SELECT group_registrations.id || ':' || group_registrations.id FROM group_registrations GROUP BY group_registrations.id ORDER BY group_registrations.id), ';') as name;`
    groupRegs := db.Query(query, nil)[0].(map[string]interface{})["name"].(string)

    query = `SELECT array_to_string(
        array(SELECT registrations.id || ':' || registrations.id FROM registrations GROUP BY registrations.id ORDER BY registrations.id), ';') as name;`
    regs := db.Query(query, nil)[0].(map[string]interface{})["name"].(string)

    return []map[string]interface{} {
        0: map[string]interface{} {
            "index": "id",
            "name": "id",
            "editable": false,
        },
        1: map[string]interface{} {
            "index": "groupreg_id",
            "name": "groupreg_id",
            "editable": true,
            "formatter": "select",
            "edittype": "select",
            "stype": "select",
            "search": true,
            "editrules": map[string]interface{}{"required": true},
            "editoptions": map[string]string{"value": groupRegs},
            "searchoptions": map[string]string{"value": ":Все;"+groupRegs},
        },
        2: map[string]interface{} {
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
    }
}
