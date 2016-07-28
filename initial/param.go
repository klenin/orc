package initial


func getOrCreateParamType(name string) (id int) {
	model := base.GetModel("param_types")
	if id = getEntityIdByName(model, name); id != -1 {
		return id
	}
	model.LoadModelData(map[string]interface{}{"name": name}).
		QueryInsert("RETURNING id").Scan(&id)
	return id
}

func loadParamTypes() {
	paramTypes := readStringsFromFile(getResourcePath("param-type-name.txt"))
	for _, paramType := range(paramTypes) {
		getOrCreateParamType(paramType);
	}
}
