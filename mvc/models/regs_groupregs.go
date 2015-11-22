package models

import (
    "github.com/klenin/orc/db"
    // "strconv"
)

type RegGroupReg struct {
    id         int `name:"id" type:"int" null:"NOT NULL" extra:"PRIMARY"`
    groupRegId int `name:"groupreg_id" type:"int" null:"NOT NULL" extra:"REFERENCES" refTable:"group_registrations" refField:"id" refFieldShow:"id"`
    regId      int `name:"reg_id" type:"int" null:"NOT NULL" extra:"REFERENCES" refTable:"registrations" refField:"id" refFieldShow:"id"`
}

func (this *RegGroupReg) GetId() int {
    return this.id
}

func (this *RegGroupReg) SetGroupRegId(groupRegId int) {
    this.groupRegId = groupRegId
}

func (this *RegGroupReg) GetGroupRegId() int {
    return this.groupRegId
}

func (this *RegGroupReg) SetRegId(regId int) {
    this.regId = regId
}

func (this *RegGroupReg) GetRegId() int {
    return this.regId
}

type RegsGroupRegsModel struct {
    Entity
}

func (*ModelManager) RegsGroupRegs() *RegsGroupRegsModel {
    model := new(RegsGroupRegsModel)
    model.SetTableName("regs_groupregs").
        SetCaption("Регистрации групп - Регистрации").
        SetColumns([]string{"id", "groupreg_id", "reg_id"}).
        SetColNames([]string{"ID", "Регистрации групп", "Регистрации"}).
        SetFields(new(RegGroupReg)).
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

func (*RegsGroupRegsModel) GetColModel(isAdmin bool, userId int) []map[string]interface{} {
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
