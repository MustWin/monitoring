package controllers

import (
  "github.com/robfig/revel"
  "monitoring/app/models"
  "monitoring/app/mail"
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

func(c App) Mail(body string) revel.Result {
  mail.Send([]string{"mikejihbe@gmail.com"}, "we@mustw.in", "Service Alert", body)
  return c.RenderText("Success!")
}
