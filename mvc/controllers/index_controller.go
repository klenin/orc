package controllers

import (
    "encoding/json"
    "github.com/klenin/orc/utils"
    "io/ioutil"
    "net/http"
    "strconv"
    "time"
)

func (c *BaseController) IndexController() *IndexController {
    return new(IndexController)
}

type IndexController struct {
    Controller
}

func (this *IndexController) Index() {
    this.Response.Header().Set("Content-Type", "text/html; charset=utf-8")
    model := this.GetModel("events")
    this.Render([]string{"mvc/views/login.html", "mvc/views/index.html"}, "index", map[string]interface{}{"events": Model{
        ColModel:  model.GetColModel(false, -1),
        TableName: model.GetTableName(),
        ColNames:  model.GetColNames(),
        Caption:   model.GetCaption()}})
}

func (this *IndexController) LoadContestsFromCats() {
    url := "http://imcs.dvfu.ru/cats/main.pl?f=contests;filter=unfinished;json=1"
    result, err := http.Get(url)
    if utils.HandleErr("[loadContestsFromCats] http.Get(url): ", err, this.Response) {
        return
    }
    defer result.Body.Close()

    body, err := ioutil.ReadAll(result.Body)
    if utils.HandleErr("[loadContestsFromCats] ioutil.ReadAll(data.Body): ", err, this.Response) {
        return
    }

    var data map[string]interface{}
    err = json.Unmarshal(body, &data)
    if utils.HandleErr("[loadContestsFromCats] json.Unmarshal(body, &data): ", err, this.Response) {
        return
    }

    for _, v := range data["contests"].([]interface{}) {
        contest := v.(map[string]interface{})
        time_, err := time.Parse("20060102T150405", contest["start_time"].(string))
        if utils.HandleErr("[loadContestsFromCats] time.Parse: ", err, this.Response) {
            continue
        }
        startDate, err := time.Parse("02.01.2006 15:04", contest["start_date"].(string))
        if utils.HandleErr("[loadContestsFromCats] time.Parse: ", err, this.Response) {
            continue
        }
        finishDate, err := time.Parse("02.01.2006 15:04", contest["finish_date"].(string))
        if utils.HandleErr("[loadContestsFromCats] time.Parse: ", err, this.Response) {
            continue
        }
        this.GetModel("events").
            LoadModelData(map[string]interface{}{
                "name":        contest["name"],
                "date_start":  startDate.Format("2006-01-02 15:04:05"),
                "date_finish": finishDate.Format("2006-01-02 15:04:05"),
                "time":        time_.Format("15:04:05"),
                "url":         "http://imcs.dvfu.ru/cats/main.pl?f=contests;cid="+strconv.Itoa(int(contest["id"].(float64))),
            }).
            QueryInsert("").
            Scan()
    }
}
