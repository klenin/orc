package models

import  "github.com/klenin/orc/db"

type Blank struct {
    personal bool
    faceId   int
    regId    int
}

func (this *Blank) SetPersonal(personal bool) *Blank {
    this.personal = personal

    return this
}

func (this *Blank) SetFaceId(faceId int) *Blank {
    this.faceId = faceId

    return this
}

func (this *Blank) SetRegId(regId int) *Blank {
    this.regId = regId

    return this
}

var queryGetBlankByRegId = `SELECT forms.id as form_id,
        forms.name as form_name, params.id as param_id,
        params.name as param_name, params.required,
        params.editable, events.name as event_name, events.id as event_id,
        param_types.name as type, param_values.id as param_val_id,
        param_values.value
    FROM events_forms
    INNER JOIN events ON events.id = events_forms.event_id
    INNER JOIN forms ON forms.id = events_forms.form_id
    INNER JOIN params ON forms.id = params.form_id
    INNER JOIN param_types ON param_types.id = params.param_type_id
    INNER JOIN param_values ON params.id = param_values.param_id
    INNER JOIN registrations ON registrations.id = param_values.reg_id
        AND events.id = registrations.event_id
    WHERE registrations.id = $1 AND forms.personal = $2
    ORDER BY forms.id, params.id;`

func (this *Blank) GetBlank() []interface{} {
    return db.Query(queryGetBlankByRegId, []interface{}{this.regId, this.personal})
}

func (this *Blank) GetEmptyBlank(eventId int) []interface {} {
    query := `SELECT forms.id as form_id, forms.name as form_name,
        params.id as param_id, params.name as param_name, params.required,
        params.editable, param_types.name as type, events.name as event_name,
        events.id as event_id
    FROM events_forms
    INNER JOIN events ON events.id = events_forms.event_id
    INNER JOIN forms ON forms.id = events_forms.form_id
    INNER JOIN params ON forms.id = params.form_id
    INNER JOIN param_types ON param_types.id = params.param_type_id
    WHERE events.id = $1 AND forms.personal = $2 ORDER BY forms.id, params.id;`

    return db.Query(query, []interface{}{eventId, this.personal})
}

//-----------------------------------------------------------------------------
type GroupBlank struct {
    groupRegId int
    Blank
}

func (this *GroupBlank) SetGroupRegId(groupRegId int) *GroupBlank {
    this.groupRegId = groupRegId

    return this
}

var queryGetBlankByGroupRegId = `SELECT forms.id as form_id,
        forms.name as form_name, params.id as param_id,
        params.name as param_name, params.required, params.editable,
        events.name as event_name, events.id as event_id,
        param_types.name as type, param_values.id as param_val_id,
        param_values.value
    FROM events_forms
    INNER JOIN events ON events.id = events_forms.event_id
    INNER JOIN forms ON forms.id = events_forms.form_id
    INNER JOIN params ON forms.id = params.form_id
    INNER JOIN param_types ON param_types.id = params.param_type_id
    INNER JOIN param_values ON params.id = param_values.param_id
    INNER JOIN registrations ON registrations.id = param_values.reg_id
    INNER JOIN faces ON faces.id = registrations.face_id
    INNER JOIN group_registrations ON group_registrations.event_id = events.id
    INNER JOIN groups ON group_registrations.group_id = groups.id
    INNER JOIN regs_groupregs ON regs_groupregs.reg_id = registrations.id
        AND regs_groupregs.groupreg_id = group_registrations.id
    WHERE group_registrations.id = $1 AND faces.id = $2 AND forms.personal = $3
    ORDER BY forms.id, params.id;`

func (this *GroupBlank) GetBlank() []interface{} {
    return db.Query(queryGetBlankByGroupRegId, []interface{}{this.groupRegId, this.faceId, this.personal})
}

var queryGetTeamBlank = `SELECT DISTINCT forms.id as form_id,
        forms.name as form_name, params.id as param_id,
        params.name as param_name, params.required, params.editable,
        events.name as event_name, events.id as event_id,
        param_types.name as type, param_values.id as param_val_id,
        param_values.value
    FROM events_forms
    INNER JOIN events ON events.id = events_forms.event_id
    INNER JOIN forms ON forms.id = events_forms.form_id
    INNER JOIN params ON forms.id = params.form_id
    INNER JOIN param_types ON param_types.id = params.param_type_id
    INNER JOIN param_values ON params.id = param_values.param_id
    INNER JOIN group_registrations ON group_registrations.event_id = events.id
    INNER JOIN groups ON group_registrations.group_id = groups.id
    INNER JOIN faces ON faces.id = groups.face_id
    INNER JOIN regs_groupregs ON
        regs_groupregs.groupreg_id = group_registrations.id
    INNER JOIN registrations ON regs_groupregs.reg_id = registrations.id
        AND registrations.event_id = events.id
        AND registrations.id = param_values.reg_id
    WHERE group_registrations.id = $1 AND forms.personal = FALSE
    ORDER BY forms.id, params.id;`

func (this *GroupBlank) GetTeamBlank() []interface{} {
    return db.Query(queryGetTeamBlank, []interface{}{this.groupRegId})
}

//-----------------------------------------------------------------------------
type BlankManager struct{}

func (*BlankManager) NewPersonalBlank(personal bool) *Blank {
    return new(Blank).SetPersonal(personal)
}

func (*BlankManager) NewGroupBlank(personal bool) *GroupBlank {
    blank := new(GroupBlank)
    blank.SetPersonal(personal)
    return blank
}

type BlankInterface interface {
    SetPersonal(bool) *Blank
    SetFaceId(int) *Blank
    GetBlank() []interface{}
}
