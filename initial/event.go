package initial

import (
	"math/rand"
	"time"
	"fmt"
	"log"
)

func generateEvents(count int) {
	eventNames := readStringsFromFile(getResourcePath("event-type.txt"))
	subjectNames := readStringsFromFile(getResourcePath("event-subject.txt"))
	for i := 0; i < count; i++ {
		eventName := fmt.Sprintf("%s по дисциплине \"%s\"",
			eventNames[rand.Intn(len(eventNames))], subjectNames[rand.Intn(len(subjectNames))])

		var secInYear int64 = 365 * 24 * 60 * 60
		timeRangeFrom := time.Now().Unix() - secInYear * 5
		timeRangeTo := time.Now().Unix() + secInYear
		timeStart := time.Unix(timeRangeFrom + rand.Int63n(timeRangeTo - timeRangeFrom), 0)
		timeFinish := time.Unix(timeStart.Unix() + rand.Int63n(7 * 24 * 60 * 60), 0)
		params := map[string]interface{}{
			"name": eventName,
			"date_start": timeStart.Format("2006-01-02"),
			"date_finish": timeFinish.Format("2006-01-02"),
			"time": timeStart.Format("15:04:05"),
			"team": rand.Int() % 3 == 2,
			"url": ""}
		var eventId int
		base.Events().LoadModelData(params).QueryInsert("RETURNING id").Scan(&eventId)

		addFormToEvent(getEntityIdByName(base.Forms(), "All field types"), eventId);
		formIds := base.Forms().SetLimit(10).Select_([]string{"id"})
		addedFormIds := map[int]bool{}
		count := rand.Intn(len(formIds))
		for i := 0; i < count; i++ {
			r := rand.Intn(len(formIds))
			id := formIds[r].(map[string]interface{})["id"].(int)
			if addedFormIds[id] {
				continue
			}
			addFormToEvent(id, eventId)
			addedFormIds[id] = true
		}
	}
}

func loadEventTypes() {
	eventTypes := readStringsFromFile(getResourcePath("event-type.txt"))
	topicality := []bool{true, false}
	for _, eventType := range eventTypes {
		params := map[string]interface{}{"name": eventType, "description": "", "topicality": topicality[rand.Intn(2)]}
		base.EventTypes().LoadModelData(params).QueryInsert("").Scan()
	}
}

func addFormToEvent(formId, eventId int) {
	if formId == -1 {
		log.Fatalln("Invalid formId")
	}
	base.GetModel("events_forms").
		LoadModelData(map[string]interface{}{"form_id": formId, "event_id": eventId}).
		QueryInsert("").Scan()
}

func createRegistrationEvent() {
	var eventId int
	base.GetModel("events").
		LoadModelData(map[string]interface{}{
		"name": "Регистрация для входа в систему",
		"date_start": "2006-01-02",
		"date_finish": "2006-01-02",
		"time": "00:00:00"}).
		QueryInsert("RETURNING id").Scan(&eventId)
	addFormToEvent(getEntityIdByName(base.Forms(), "Регистрационные данные"), eventId);
	addFormToEvent(getEntityIdByName(base.Forms(), "Общие сведения"), eventId);
}
