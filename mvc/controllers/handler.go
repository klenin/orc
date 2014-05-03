package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/orc/utils"
	"strconv"
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
	utils.HandleErr("[Handler] Decode :", err)
	data := request.(map[string]interface{})

	switch data["action"] {
	case "register":
		login, password := data["login"].(string), data["password"].(string)
		response = this.HandleRegister(login, password)
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

	case "change-pass":
		id := data["id"].(string)
		pass := data["pass"].(string)
		model := GetModel("users")
		result, _ := model.Select([]string{"id", id}, []string{"salt"})
		salt := result[0].(map[string]interface{})["salt"].(string)
		hash := GetMD5Hash(pass + salt)
		model.Update([]string{"pass"}, []interface{}{hash, id}, "id=$2")
		break

	case "select":
		tableName := data["table"].(string)
		count, err := strconv.Atoi(data["count"].(string))
		utils.HandleErr("[Handle select] strconv.Atoi: ", err)
		fields := data["fields"].([]interface{})
		model := GetModel(tableName)
		result, _ := model.Select(nil, utils.ArrayInterfaceToString(fields, count))
		answer, err := json.Marshal(map[string]interface{}{
			"data": result})
		utils.HandleErr("[Handle select] json.Marshal: ", err)
		fmt.Fprintf(this.Response, "%s", string(answer))
		break

	case "getsubgrid":
		tableName := data["table"].(string)
		id := data["id"].(string)
		index, _ := strconv.Atoi(data["index"].(string))
		model := GetModel(tableName)
		subTableName := model.SubTable[index]
		subModel := GetModel(subTableName)
		result, refdata := subModel.Select([]string{model.SubField, id}, subModel.Columns)
		answer, err := json.Marshal(map[string]interface{}{
			"data":      result,
			"name":      subModel.TableName,
			"caption":   subModel.Caption,
			"colnames":  subModel.ColNames,
			"columns":   subModel.Columns,
			"refdata":   refdata,
			"reffields": subModel.RefFields})
		utils.HandleErr("[HandleRegister] json.Marshal: ", err)
		fmt.Fprintf(this.Response, "%s", string(answer))
		break
	}
}
