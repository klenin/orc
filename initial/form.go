package initial


func createFormParam(name string, formId, typeId, identifier int) {
	params := base.GetModel("params")
	params.LoadModelData(map[string]interface{}{
		"name":          name,
		"form_id":       formId,
		"param_type_id": typeId,
		"identifier":    identifier,
		"required":      true,
		"editable":      true}).
		QueryInsert("").Scan()
}

var paramIdentifier = 0

func loadForms() {
	forms := readJsonFile(getResourcePath("form.json"))
	for _, v := range(forms.([]interface{})) {
		form := v.(map[string]interface{})
		var formId int
		base.Forms().LoadModelData(map[string]interface{}{
			"name": form["name"].(string),
			"personal": form["personal"].(bool),
		}).
			QueryInsert("RETURNING id").Scan(&formId)
		for _, v := range(form["params"].([]interface{})) {
			param := v.(map[string]interface{})
			paramIdentifier++
			createFormParam(param["name"].(string), formId, getOrCreateParamType(param["type"].(string)), paramIdentifier)
		}
	}
}
