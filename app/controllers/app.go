package controllers

import (
  "github.com/robfig/revel"
  "monitoring/app/models"
  "monitoring/app/mail"
  "time"
)

type App struct {
	*revel.Controller
}

func (c App) Index() revel.Result {
  services, _ := models.AllServices()
	return c.Render(services)
}

func (c App) Status(app string) revel.Result {
  service := models.FindServiceByName(app)
  return c.RenderJson(service)
}

func (c App) Delete() revel.Result {
  var name string
  c.Params.Bind(&name, "name")
  service := models.FindServiceByName(name)
  models.GetDb().Delete(service)
  return c.RenderJson("{}")
}

func (c App) Create() revel.Result {
  var service models.Service
 // var name, url string
  c.Params.Bind(&service.Name, "name")
  c.Params.Bind(&service.Url, "url")
  service.UpdatedAt = time.Now()
  service.CreatedAt = time.Now()
  models.GetDb().Save(&service)
  return c.Redirect("/")
}

func(c App) Mail(body string) revel.Result {
  mail.Send([]string{"mikejihbe@gmail.com"}, "we@mustw.in", "Service Alert", body)
  return c.RenderText("Success!")
}
