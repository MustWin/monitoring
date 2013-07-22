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
  services := make([]models.Service, 0)
  err := models.GetDb().FindAll(&services)
  revel.INFO.Print(services)
  if (err != nil) {
    revel.ERROR.Println(err)
  }
  /*apps := []models.Service {
              models.Service{"app1", "blah", models.Status{"healthy", "healthy", "healthy", time.Now()}, true, false}, 
              models.Service{"app2", "blah", models.Status{"error", "down", "error", time.Now()}, false, false},
              models.Service{"app3", "blah", models.Status{"down", "unknown", "down", time.Now()}, false, true},
          }*/
	return c.Render(services)
}

func (c App) Status() revel.Result {
  app := models.Service{1, "app2", "blah", "{}"/*models.Status{"error", "down", "error", time.Now()}*/, false, false}
  return c.RenderJson(app)
}
