package tasks

import (
  "github.com/robfig/revel"
  "monitoring/app/models"
  "net"
  "net/http"
  "time"
  "fmt"
  "io/ioutil"
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


func checkService(service *models.Service) {
  req, _ := http.NewRequest("GET", service.Url, nil)
  req.Header.Add("User-Agent", "MustWin/health-checker")
  resp, err := client.Do(req)
  if (err != nil) {
  // PICK UP HERE, FOR SOME REASON THE SERVICE THAT IS SAVED LOSES ALL THE REST OF THE DATAS
    // This handles timeouts, 500s, etc
    revel.ERROR.Println(err)
    service.Healthy = false
    service.Status = fmt.Sprintf("{\"error\": \"%s\"}", err)
    models.UpdateService(service)
    return
  }
  // Fetch the body
  defer resp.Body.Close()
  body, err := ioutil.ReadAll(resp.Body)
  if (err != nil) {
    revel.ERROR.Println(err)
    service.Healthy = false
    service.Status = fmt.Sprintf("{\"error\": \"%s\"}", err)
    models.UpdateService(service)
    return
  }

  // Update the state
  // TODO: maybe be smarter with json?
  service.Status = string(body)

  if (resp.StatusCode != 200) {
    service.Healthy = false
    if (service.Status != string(body)) {
      // TODO: add time condition to re-send alerts
      // TODO: SEND ALERTS, RAISE HELL
      revel.ERROR.Println(err)
      revel.ERROR.Printf("Error in %s: %s", service.Name, body)
    }
  } else {
    // Reset
    service.Healthy = true
    service.Acked = false
  }
  models.UpdateService(service)
}

type CheckServices struct {

}

func (cs CheckServices) Run() {
  services, _ := models.AllServices()
  for _, service := range services {
    go checkService(&service)
  }
}
