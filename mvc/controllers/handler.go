package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/orc/db"
	"github.com/orc/sessions"
	"github.com/orc/utils"
	"html/template"
	"strconv"
)

func (c *BaseController) Handler() *Handler {
	return new(Handler)
}

type Handler struct {
	Controller
}

func (this *Handler) GetEventList() {
	this.Response.Header().Set("Access-Control-Allow-Origin", "*")
	this.Response.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	this.Response.Header().Set("Content-type", "application/json")

	var request map[string]interface{}
	decoder := json.NewDecoder(this.Request.Body)
	err := decoder.Decode(&request)
	utils.HandleErr("[Handler] Decode :", err, this.Response)

	fields := request["fields"].([]interface{})
	tableName := request["table"].(string)

	model := GetModel(tableName)
	result, _ := model.Select(nil, "", utils.ArrayInterfaceToString(fields))

	response, err := json.Marshal(map[string]interface{}{"data": result})
	utils.HandleErr("[Handle select] json.Marshal: ", err, nil)
	fmt.Fprintf(this.Response, "%s", string(response))
}

func (this *Handler) ResetPassword() {
	if flag := sessions.CheackSession(this.Response, this.Request); !flag {
		return
	}
	this.Response.Header().Set("Access-Control-Allow-Origin", "*")
	this.Response.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	this.Response.Header().Set("Content-type", "application/json")

	var request map[string]string
	decoder := json.NewDecoder(this.Request.Body)
	err := decoder.Decode(&request)
	utils.HandleErr("[Handler] Decode :", err, this.Response)

	id, pass, model := request["id"], request["pass"], GetModel("users")
	result, _ := model.Select([]string{"id", id}, "", []string{"salt"})
	salt := result[0].(map[string]interface{})["salt"].(string)
	hash := GetMD5Hash(pass + salt)
	model.Update([]string{"pass"}, []interface{}{hash, id}, "id=$2")

	response, err := json.Marshal(map[string]interface{}{"result": "ok"})
	utils.HandleErr("[Handle select] json.Marshal: ", err, nil)
	fmt.Fprintf(this.Response, "%s", string(response))
}

func (this *Handler) Index() {
	var data map[string]interface{}
	response := ""

	this.Response.Header().Set("Access-Control-Allow-Origin", "*")
	this.Response.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	this.Response.Header().Set("Content-type", "application/json")

	decoder := json.NewDecoder(this.Request.Body)
	err := decoder.Decode(&data)
	utils.HandleErr("[Handler] Decode :", err, this.Response)

	switch data["action"] {
	case "register":
		login, password := data["login"].(string), data["password"].(string)
		fname, lname, pname := data["fname"].(string), data["lname"].(string), data["pname"].(string)
		response = this.HandleRegister(login, password, "user", fname, lname, pname)
		fmt.Fprintf(this.Response, "%s", response)
		break

	case "login":
		response = this.HandleLogin(data["login"].(string), data["password"].(string))
		fmt.Fprintf(this.Response, "%s", response)
		break

	case "logout":
		response = this.HandleLogout()
		fmt.Fprintf(this.Response, "%s", response)
		break

	case "editProfile":
		id := data["id"].(string)
		tableName := data["table"].(string)
		inf := data["data"].([]interface{})

		var fields []string
		var params []interface{}
		for _, element := range inf {
			fields = append(fields, element.(map[string]interface{})["name"].(string))
			params = append(params, element.(map[string]interface{})["value"])
		}
		params = append(params, id)
		model := GetModel(tableName)
		model.Update(fields, params, "id=$"+strconv.Itoa(len(fields)+1))

		response, err := json.Marshal(map[string]interface{}{"result": "ok"})
		utils.HandleErr("[Handle select] json.Marshal: ", err, nil)
		fmt.Fprintf(this.Response, "%s", string(response))
		break
	}
}

func (this *Handler) ShowCabinet(tableName string) {
	if flag := sessions.CheackSession(this.Response, this.Request); !flag {
		return
	}
	table := GetModel("users")
	login := sessions.GetValue("name", this.Request).(string)
	println(login)
	data, _ := table.Select([]string{"login", login}, "", []string{"role", "person_id"})

	role := data[0].(map[string]interface{})["role"].(string)
	person_id := data[0].(map[string]interface{})["person_id"].(int64)
	println("role", role)

	var model Model
	if role == "admin" {
		model = Model{Columns: db.Tables, ColNames: db.TableNames}
	} else if role == "user" {
		m := GetModel("persons")
		data, _ := m.Select([]string{"id", strconv.Itoa(int(person_id))}, "", m.GetColumns())
		model = Model{Caption: login, Table: data, Columns: m.GetColumns(), ColNames: m.GetColNames()}
	}

	tmp, err := template.ParseFiles(
		"mvc/views/"+role+".html",
		"mvc/views/header.html",
		"mvc/views/footer.html")
	utils.HandleErr("[Handler.ShowCabinet] ParseFiles: ", err, nil)
	err = tmp.ExecuteTemplate(this.Response, role, model)
	utils.HandleErr("[Handler.ShowCabinet] Execute: ", err, nil)
}
