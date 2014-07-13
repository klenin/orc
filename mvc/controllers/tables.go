package controllers

import (
	"fmt"
	"encoding/json"
	"github.com/orc/db"
	"github.com/orc/mvc/models"
	"github.com/orc/sessions"
	"github.com/orc/utils"
	"html/template"
	"strconv"
	"strings"
)

type Model struct {
	Id        string
	TableName string
	Caption   string
	Table     []interface{}
	RefData   map[string]interface{}
	RefFields []string
	Columns   []string
	ColNames  []string
}

type A struct {
	E []interface{} //events
	T []interface{} //event_types
	F []interface{} //forms
	P []interface{} //params
}

func GetModel(tableName string) models.Entity {
	base := new(models.ModelManager)
	var model models.Entity
	switch tableName {
	case "events":
		model = base.Events().Entity
		break
	case "event_types":
		model = base.EventTypes().Entity
		break
	case "events_types":
		model = base.EventsTypes().Entity
		break
	case "teams":
		model = base.Teams().Entity
		break
	case "persons":
		model = base.Persons().Entity
		break
	case "users":
		model = base.Users().Entity
		break
	case "teams_persons":
		model = base.TeamsPersons().Entity
		break
	case "forms":
		model = base.Forms().Entity
		break
	case "params":
		model = base.Params().Entity
		break
	case "forms_types":
		model = base.FormsTypes().Entity
		break
	case "param_values":
		model = base.ParamValues().Entity
		break
	case "persons_events":
		model = base.PersonsEvents().Entity
		break
	}
	return model
}

func (this *Handler) ShowCabinet(tableName string) {
	if flag := sessions.CheackSession(this.Response, this.Request); !flag {
		return
	}

	login := sessions.GetValue("name", this.Request)

	table := GetModel("users")
	data, _ := table.Select([]string{"login", login}, "", []string{"role", "person_id"})

	role := data[0].(map[string]interface{})["role"].(string)
	person_id := data[0].(map[string]interface{})["person_id"].(int64)

	var model Model
	if role == "admin" {
		model = Model{Columns: db.Tables, ColNames: db.TableNames}
	} else if role == "user" {
		m := GetModel("persons")
		data, _ := m.Select([]string{"id", strconv.Itoa(int(person_id))}, "", m.Columns)
		model = Model{Caption: login, Table: data, Columns: m.Columns, ColNames: m.ColNames}
	}

	tmp, err := template.ParseFiles(
		"mvc/views/"+role+".html",
		"mvc/views/header.html",
		"mvc/views/footer.html")
	utils.HandleErr("[Handler.ShowCabinet] ParseFiles: ", err, nil)
	err = tmp.ExecuteTemplate(this.Response, role, model)
	utils.HandleErr("[Handler.ShowCabinet] Execute: ", err, nil)
}

func (this *Handler) Select(tableName string) {
	if flag := sessions.CheackSession(this.Response, this.Request); !flag {
		return
	}
	model := GetModel(tableName)
	answer, refdata := model.Select(nil, "", model.Columns)
	tmp, err := template.ParseFiles(
		"mvc/views/table.html",
		"mvc/views/header.html",
		"mvc/views/footer.html")
	utils.HandleErr("[Handler.Select] template.ParseFiles: ", err, nil)
	err = tmp.ExecuteTemplate(this.Response, "table", Model{
		Table:     answer,
		RefData:   refdata,
		RefFields: model.RefFields,
		TableName: model.TableName,
		ColNames:  model.ColNames,
		Columns:   model.Columns,
		Caption:   model.Caption})
	utils.HandleErr("[Handler.Select] tmp.Execute: ", err, nil)
	fmt.Println(answer)
}

func (this *Handler) Edit(tableName string) {
	if flag := sessions.CheackSession(this.Response, this.Request); !flag {
		return
	}
	var i int
	oper := this.Request.FormValue("oper")
	model := GetModel(tableName)
	params := make([]interface{}, len(model.Columns)-1)

	for i = 0; i < len(model.Columns)-1 && this.Request.FormValue(model.Columns[i+1]) != ""; i++ {
		if model.Columns[i+1] == "date" {
			params[i] = this.Request.FormValue(model.Columns[i+1])[0:10]
		} else {
			params[i] = this.Request.FormValue(model.Columns[i+1])
		}
	}

	switch oper {
	case "edit":
		params = append(params, this.Request.FormValue("id"))
		model.Update(model.Columns[1:], params, "id=$"+strconv.Itoa(i+1))
		break
	case "add":
		model.Insert(model.Columns[1:], params)
		break
	case "del":
		ids := strings.Split(this.Request.FormValue("id"), ",")
		p := make([]interface{}, len(ids))
		for i, v := range ids {
			p[i] = interface{}(v)
		}
		model.Delete("id", p)
		break
	}
}

func MegoJoin(tableName, id string) A {
	var E []interface{}
	var T []interface{}
	var F []interface{}
	var P []interface{}

	E = db.Select("events", []string{"id", id}, "", []string{"id", "name"})

	q := db.InnerJoin(
		[]string{"id", "name"},
		"t",
		"events_types",
		"e_t",
		[]string{"event_id", "type_id"},
		[]string{"events", "event_types"},
		[]string{"e", "t"},
		[]string{"id", "id"},
		"where e.id=$1")

	rows := db.Query(q, []interface{}{id})
	rowsInf := db.Exec(q, []interface{}{id})
	l, _ := rowsInf.RowsAffected()
	c, _ := rows.Columns()
	T = db.ConvertData(c, l, rows)

	for i := 0; i < len(T); i++ {
		item := T[i]
		id := item.(map[string]interface{})["id"]

		q := db.InnerJoin(
			[]string{"id", "name"},
			"f",
			"forms_types",
			"f_t",
			[]string{"form_id", "type_id"},
			[]string{"forms", "event_types"},
			[]string{"f", "t"},
			[]string{"id", "id"},
			"where t.id=$1")

		rows := db.Query(q, []interface{}{id})
		rowsInf := db.Exec(q, []interface{}{id})
		l, _ := rowsInf.RowsAffected()
		c, _ := rows.Columns()
		F = append(F, db.ConvertData(c, l, rows))
	}

	for i := 0; i < len(F); i++ {
		I := F[i]
		var PP []interface{}
		for j := 0; j < len(I.([]interface{})); j++ {
			item := I.([]interface{})[j]
			id := item.(map[string]interface{})["id"]

			q := db.InnerJoin(
				[]string{"id", "name", "type"},
				"p",
				"params",
				"p",
				[]string{"form_id"},
				[]string{"forms"},
				[]string{"f"},
				[]string{"id"},
				"where f.id=$1")

			rows := db.Query(q, []interface{}{id})
			rowsInf := db.Exec(q, []interface{}{id})
			l, _ := rowsInf.RowsAffected()
			c, _ := rows.Columns()
			PP = append(PP, db.ConvertData(c, l, rows))
		}
		P = append(P, PP)
	}
	return A{E: E, T: T, F: F, P: P}
}

func (this *Handler) Show(tableName, id string) {
	if flag := sessions.CheackSession(this.Response, this.Request); !flag {
		return
	}
	tmp, err := template.ParseFiles(
		"mvc/views/item.html",
		"mvc/views/header.html",
		"mvc/views/footer.html")
	utils.HandleErr("[Handler.Show] template.ParseFiles: ", err, nil)

	a, err := json.Marshal(MegoJoin(tableName, id))
	utils.HandleErr("[Handler.Show] template.json.Marshal: ", err, nil)

	err = tmp.ExecuteTemplate(this.Response, "item", template.JS(a))
	utils.HandleErr("[Handler.Show] tmp.Execute: ", err, nil)
}
