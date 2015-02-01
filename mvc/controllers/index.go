package controllers

import (
    "github.com/orc/utils"
    "html/template"
)

func (c *BaseController) Index() *IndexController {
    return new(IndexController)
}

type IndexController struct {
    Controller
}

func (this *IndexController) Index() {
    this.Response.Header().Set("Content-Type", "text/html; charset=utf-8")

    tmp, err := template.ParseFiles("mvc/views/index.html", "mvc/views/header.html", "mvc/views/footer.html")
    utils.HandleErr("[IndexController::Index] ParseFiles: ", err, this.Response)

    err = tmp.ExecuteTemplate(this.Response, "index", nil)
    utils.HandleErr("[IndexController::Index] ExecuteTemplate: ", err, this.Response)
}
