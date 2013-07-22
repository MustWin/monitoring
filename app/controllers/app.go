package controllers

import (
  "github.com/robfig/revel"
  //"time"
  "monitoring/app/models"
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
