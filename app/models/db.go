package models

import (
  "github.com/robfig/revel"
  "database/sql"
  "time"
  "github.com/astaxie/beedb"
  _ "github.com/bmizerany/pq"
)

// Models

type Status struct {
  State string `json:"state"`
  Database string `json:"database"`
  App string `json:"app,omitempty"`
  UpdatedAt time.Time `json:"updated_at,omitempty"`
}

type Service struct {
  Id int64 `json:"id"`
  Name string `json:"name"`
  Url string `json:"url"`
  Status string `json:"status"`
  Healthy bool `json:"healthy"`
  Acked bool `json:"acked"`
}

var orm *beedb.Model = nil
func  GetDb() *beedb.Model {
  if (orm != nil) {
    return orm
  }
  rawdb, err := sql.Open("postgres", "user=monitoring dbname=monitoring sslmode=disable")
  if (err != nil) {
    revel.ERROR.Println(err)
  }
  // construct a gorp DbMap
  orm := beedb.New(rawdb, "pg")
  beedb.OnDebug=true
  beedb.PluralizeTableNames=true
  return &orm
}
