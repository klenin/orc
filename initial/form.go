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

func loadForms() {
	formNames := readStringsFromFile(getResourcePath("form-name.txt"))
	for _, formName := range(formNames) {
		base.Forms().
			LoadModelData(map[string]interface{}{"name": formName, "personal": true}).
			QueryInsert("").
			Scan()
	}
}

func getOrCreateRegForm() (id int) {
	name := "Регистрационные данные"
	model := base.GetModel("forms")
	if id = getEntityIdByName(model, name); id != -1 {
		return id
	}
	model.LoadModelData(map[string]interface{}{"name": name, "personal": true}).
		QueryInsert("RETURNING id").Scan(&id)
	paramTextTypeId := getOrCreateParamType("text")
	paramPassTypeId := getOrCreateParamType("password")
	createFormParam("Логин", id, paramTextTypeId, 2)
	createFormParam("Пароль", id, paramPassTypeId, 3)
	createFormParam("Подтвердите пароль", id, paramPassTypeId, 4)
	createFormParam("E-mail", id, paramTextTypeId, 5)
	return id
}

func getOrCreateNamesForm() (id int) {
	name := "Общие сведения"
	model := base.GetModel("forms")
	if id = getEntityIdByName(model, name); id != -1 {
		return id
	}
	model.LoadModelData(map[string]interface{}{"name": name, "personal": true}).
		QueryInsert("RETURNING id").Scan(&id)
	paramTextTypeId := getOrCreateParamType("text")
	createFormParam("Фамилия", id, paramTextTypeId, 6)
	createFormParam("Имя", id, paramTextTypeId, 7)
	createFormParam("Отчество", id, paramTextTypeId, 8)
	return id
}
