package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/orc/db"
	"github.com/orc/sessions"
	"github.com/orc/utils"
	"reflect"
	"strconv"
	"time"
)

func (c *BaseController) Handler() *Handler {
	return new(Handler)
}

type Handler struct {
	Controller
}

func (this *Handler) Index() {
	var (
		request  interface{}
		response = ""
	)
	this.Response.Header().Set("Access-Control-Allow-Origin", "*")
	this.Response.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	this.Response.Header().Set("Content-type", "application/json")

	decoder := json.NewDecoder(this.Request.Body)
	err := decoder.Decode(&request)
	utils.HandleErr("[Handler] Decode :", err, this.Response)
	if err != nil {
		return
	}
	data := request.(map[string]interface{})

	switch data["action"] {
	case "register":
		login, password := data["login"].(string), data["password"].(string)
		fname, lname, pname := data["fname"].(string), data["lname"].(string), data["pname"].(string)
		response = this.HandleRegister(login, password, "user", fname, lname, pname)
		fmt.Fprintf(this.Response, "%s", response)
		break

	case "login":
		login, password := data["login"].(string), data["password"].(string)
		response = this.HandleLogin(login, password)
		fmt.Fprintf(this.Response, "%s", response)
		break

	case "logout":
		response = this.HandleLogout()
		fmt.Fprintf(this.Response, "%s", response)
		break

	//admin
	case "change-pass":
		id := data["id"].(string)
		pass := data["pass"].(string)
		model := GetModel("users")
		result, _ := model.Select([]string{"id", id}, "", []string{"salt"})
		salt := result[0].(map[string]interface{})["salt"].(string)
		hash := GetMD5Hash(pass + salt)
		model.Update([]string{"pass"}, []interface{}{hash, id}, "id=$2")
		r := map[string]interface{}{"result": "ok"}
		answer, err := json.Marshal(r)
		utils.HandleErr("[Handle select] json.Marshal: ", err, nil)
		fmt.Fprintf(this.Response, "%s", string(answer))
		break

	//select events
	case "select":
		tableName := data["table"].(string)
		fields := data["fields"].([]interface{})
		model := GetModel(tableName)
		result, _ := model.Select(nil, "", utils.ArrayInterfaceToString(fields))
		fmt.Println("result: ", result)
		answer, err := json.Marshal(map[string]interface{}{"data": result})
		utils.HandleErr("[Handle select] json.Marshal: ", err, nil)
		fmt.Fprintf(this.Response, "%s", string(answer))
		break

	case "are-alive":
		answer1, err := json.Marshal(sessions.CheackSession(this.Response, this.Request))
		utils.HandleErr("[Handle are_alive] json.Marshal: ", err, nil)
		fmt.Fprintf(this.Response, "%s", string(answer1))
		break

	case "getsubgrid":
		tableName := data["table"].(string)
		id := data["id"].(string)
		index, _ := strconv.Atoi(data["index"].(string))
		model := GetModel(tableName)
		subTableName := model.GetSubTable(index)
		subModel := GetModel(subTableName)
		result, refdata := subModel.Select([]string{model.GetSubField(), id}, "", subModel.GetColumns())
		answer, err := json.Marshal(map[string]interface{}{
			"data":      result,
			"name":      subModel.GetTableName(),
			"caption":   subModel.GetCaption(),
			"colnames":  subModel.GetColNames(),
			"columns":   subModel.GetColumns(),
			"refdata":   refdata,
			"reffields": subModel.GetRefFields()})
		utils.HandleErr("[HandleRegister] json.Marshal: ", err, nil)
		fmt.Fprintf(this.Response, "%s", string(answer))
		break

	//update table "persons"
	case "update":
		tableName := data["table"].(string)
		id := data["id"].(string)
		inf := data["data"].([]interface{})

		var fields []string
		var params []interface{}
		for _, element := range inf {
			e := element.(map[string]interface{})
			fields = append(fields, e["name"].(string))
			params = append(params, e["value"])
		}
		params = append(params, id)
		model := GetModel(tableName)
		model.Update(fields, params, "id=$"+strconv.Itoa(len(fields)+1))
		r := map[string]interface{}{"result": "ok"}
		answer, err := json.Marshal(r)
		utils.HandleErr("[Handle select] json.Marshal: ", err, nil)
		fmt.Fprintf(this.Response, "%s", string(answer))
		break

	//insert into param_values
	case "save-update":
		event_id := int(data["event_id"].(float64))

		oper := data["oper"].(string)

		id := sessions.GetValue("id", this.Request)
		users := GetModel("users")
		d, _ := users.Select([]string{"id", id}, "", []string{"person_id"})
		person_id := int(d[0].(map[string]interface{})["person_id"].(int64))

		persons_events := GetModel("persons_events")
		d, _ = persons_events.Select([]string{"person_id", strconv.Itoa(person_id), "event_id", strconv.Itoa(event_id)}, "AND", []string{"person_id"})

		var r interface{}
		inf := data["data"].([]interface{})
		param_values := GetModel("param_values")

		t := time.Now()

		if len(d) == 0 && oper == "save" {

			persons_events.Insert([]string{"person_id", "event_id", "reg_date", "last_date"}, []interface{}{person_id, event_id, t.Format("2006-01-02"), t.Format("2006-01-02")})
			for _, element := range inf {
				e := element.(map[string]interface{})
				param_id := e["name"]
				value := e["value"]
				param_values.Insert([]string{"person_id", "event_id", "param_id", "value"}, []interface{}{person_id, event_id, param_id, value})
			}
			r = map[string]interface{}{"result": "ok"}

		} else if len(d) != 0 && oper == "update" {
			for _, element := range inf {
				e := element.(map[string]interface{})
				param_id := e["name"]
				value := e["value"]
				param_values.Update([]string{"value"}, []interface{}{value, person_id, event_id, param_id}, "person_id=$"+strconv.Itoa(2)+" AND event_id=$"+strconv.Itoa(3)+" AND param_id=$"+strconv.Itoa(4))
			}
			persons_events.Update([]string{"last_date"}, []interface{}{t.Format("2006-01-02"), person_id, event_id}, "person_id=$"+strconv.Itoa(2)+" AND event_id=$"+strconv.Itoa(3))
			r = map[string]interface{}{"result": "ok"}
		} else {
			r = map[string]interface{}{"result": "exists"}
		}

		answer, err := json.Marshal(r)
		utils.HandleErr("[Handle select] json.Marshal: ", err, nil)
		fmt.Fprintf(this.Response, "%s", string(answer))
		break

	case "select-h-events":
		ids := utils.ArrayInterfaceToString(data["form_ids"].([]interface{}))

		id := sessions.GetValue("id", this.Request)
		users := GetModel("users")
		d, _ := users.Select([]string{"id", id}, "", []string{"person_id"})
		person_id := int(d[0].(map[string]interface{})["person_id"].(int64))

		model := GetModel("forms_types")
		result, _ := model.Select(ids, "OR", []string{"type_id"})
		fmt.Println("result: ", result)

		q := `SELECT DISTINCT event_id, name FROM param_values 
		inner join events on events.id = param_values.event_id
		WHERE event_id IN (SELECT DISTINCT event_id FROM events_types WHERE `
		s := reflect.ValueOf(result)
		var i int
		var params []interface{}
		for i = 1; i < s.Len(); i++ {
			q += "type_id=$" + strconv.Itoa(i) + " OR "
			params = append(params, result[i-1].(map[string]interface{})["type_id"])
		}
		q += "type_id=$" + strconv.Itoa(i) + ") AND person_id=$" + strconv.Itoa(i+1)
		params = append(params, result[i-1].(map[string]interface{})["type_id"])
		params = append(params, person_id)
		fmt.Println("params: ", params)
		rows := db.Query(q, params)
		rowsInf := db.Exec(q, params)
		l, _ := rowsInf.RowsAffected()
		c, _ := rows.Columns()
		T := db.ConvertData(c, l, rows)

		answer, err := json.Marshal(T)
		utils.HandleErr("[Handle select] json.Marshal: ", err, nil)
		fmt.Fprintf(this.Response, "%s", string(answer))
		break

	case "get-h-forms":
		event_id := data["event_id"].(string)

		id := sessions.GetValue("id", this.Request)
		users := GetModel("users")
		d, _ := users.Select([]string{"id", id}, "", []string{"person_id"})
		person_id := int(d[0].(map[string]interface{})["person_id"].(int64))

		q := `select param_id, p.name param_name, p.type, value, form_id, forms.name form_name from param_values 
			inner join params p on param_values.param_id = p.id
			inner join forms on forms.id = p.form_id
			where person_id = $1 and event_id = $2;`
		rows := db.Query(q, []interface{}{person_id, event_id})
		rowsInf := db.Exec(q, []interface{}{person_id, event_id})
		l, _ := rowsInf.RowsAffected()
		c, _ := rows.Columns()
		R := db.ConvertData(c, l, rows)

		answer, err := json.Marshal(R)
		utils.HandleErr("[Handle select] json.Marshal: ", err, nil)
		fmt.Fprintf(this.Response, "%s", string(answer))
		break
	}
}
