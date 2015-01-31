package controllers

import (
	"github.com/orc/mvc/models"
	"net/http"
)

type BaseController struct{}

type Controller struct {
	Request  *http.Request
	Response http.ResponseWriter
}

type Model struct {
	Id        string
	TableName string
	Caption   string
	Table     []interface{}
	RefData   map[string]interface{}
	RefFields []string
	Columns   []string
	ColNames  []string
	Sub       bool
}

type RequestModel struct {
	E []interface{} //events
	T []interface{} //event_types
	F []interface{} //forms
	P []interface{} //params
}

func GetModel(tableName string) models.VirtEntity {
	base := new(models.ModelManager)
	switch tableName {
	case "events":
		return base.Events()
	case "event_types":
		return base.EventTypes()
	case "events_types":
		return base.EventsTypes()
	//case "teams":
	//	return base.Teams()
	case "persons":
		return base.Persons()
	case "users":
		return base.Users()
	//case "teams_persons":
	//	return base.TeamsPersons()
	case "forms":
		return base.Forms()
	case "params":
		return base.Params()
	case "forms_types":
		return base.FormsTypes()
	case "param_values":
		return base.ParamValues()
	case "persons_events":
		return base.PersonsEvents()
	}
	panic(nil)
}
