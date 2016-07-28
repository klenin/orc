package initial

import (
	"github.com/klenin/orc/mvc/models"
	"io/ioutil"
	"strings"
	"log"
)

func readStringsFromFile(fileName string) []string {
	content, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatalln("loadData:", err.Error())
	}
	array := strings.Split(string(content), "\n")
	var r []string
	for _, str := range array {
		if str = strings.TrimSpace(str); str != "" {
			r = append(r, str)
		}
	}
	return r
}

func getEntityIdByName(model models.EntityInterface, name string) int {
	if items := model.LoadWherePart(map[string]interface{}{"name": name}).Select_([]string{"id"}); len(items) > 0 {
		return items[0].(map[string]interface{})["id"].(int)
	}
	return -1
}

func getResourcePath(filename string) string {
	return "./initial/resources/" + filename
}
