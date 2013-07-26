package tasks

import (
  "github.com/robfig/revel"
  "monitoring/app/models"
  "monitoring/app/mail"
  "net"
  "net/http"
  "time"
  "fmt"
  "io/ioutil"
  "errors"
)

// Configure the HTTP client

var timeout = time.Duration(10 * time.Second)

func dialTimeout(network, addr string) (net.Conn, error) {
  return net.DialTimeout(network, addr, timeout)
}

var transport = http.Transport{
  Dial: dialTimeout,
}

var client = &http.Client{
  Transport: &transport,
}

func markAsDown(service *models.Service, status string) {
  service.Healthy = false
  revel.ERROR.Printf("")
  if service.Status != status {
    revel.ERROR.Printf("Error - Marking Service %s as Down", service.Name)
    // ALERT
    mail.Send(
      []string{"we@mustw.in"}, 
      "we@mustw.in", 
      fmt.Sprintf("Service Alert - %s", service.Name), 
      fmt.Sprintf("%s is DOWN.  Here's the status reported from our monitor:\n\n" + status, service.Name),
    )
  }
  service.Status = status
  models.UpdateService(service)
}

func checkService(service models.Service, tries int) (models.Service, error) {
  req, _ := http.NewRequest("GET", service.Url, nil)
  req.Header.Add("User-Agent", "MustWin/health-checker")
  resp, err := client.Do(req)
  if (err != nil) {
    // This handles timeouts, 500s, etc
    if (tries < 1) {
      return checkService(service, tries + 1)
    }
    markAsDown(&service, fmt.Sprintf("{\"error\": \"%s\"}", err))
    return service, err
  }
  // Fetch the body
  defer resp.Body.Close()
  body, err := ioutil.ReadAll(resp.Body)
  if (err != nil) {
    if (tries < 1) {
      return checkService(service, tries + 1)
    }
    markAsDown(&service, fmt.Sprintf("{\"error\": \"%s\"}", err))
    return service, err
  }

  if (resp.StatusCode != 200) {
    if (tries < 1) {
      return checkService(service, tries + 1)
    }
    // TODO: maybe be smarter with json?
    service.Healthy = false
    markAsDown(&service, string(body))
    return service, errors.New(string(body))
  } else {
    // TODO: healthcheck style problem detection. Handle Successful response with error info inside.
    // Reset the service
    service.Healthy = true
    if (service.Status != string(body)) {
      revel.INFO.Printf("Marking Service %s as OK", service.Name)
    }
    service.Status = string(body)
    service.Acked = false
    models.UpdateService(&service)
    return service, nil
  }
}

type CheckServices struct {}

func (cs CheckServices) Run() {
  services, _ := models.AllServices()
  revel.INFO.Println("Performing Service Check")
  for _, service := range services {
    go checkService(service, 0)
  }
}
