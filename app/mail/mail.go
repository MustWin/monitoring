package mail


import (
  "github.com/robfig/revel"
  "net/smtp"
  "strconv"
)

/*
2 sample app.conf entries

mail.type = smtp
mail.host = "mail.google.com" // 127.0.0.1
mail.port = 465 // 25
mail.user = "mike@mustw.in" // optional
mail.password = "membersonly" // optional

// OR

mail.type = file

*/

type MailConfig struct {
  Type string
  Host string
  Port int
  AuthType string
  User string
  Password string
}


var (
  // Singleton instance of the underlying job scheduler.
  config MailConfig

  // This limits the number of jobs allowed to run concurrently.
  auth smtp.Auth
)

func Send(to []string, from string, subject string, body string) {
  if config.Type == "file" {
    revel.INFO.Println("\n===========================================\n" + 
        "SENDING EMAIL  -- Subject: " + subject + "\n===========================================\n" +
        body +
        "\n===========================================")
  } else {
    // Only plaintext single part support, for now
    mime := "MIME-version: 1.0;\nContent-Type: text/plain; charset=\"UTF-8\";\n\n";
    subject = "Subject: " + subject + "\n"
    msg := []byte(subject + mime + body)
    err := smtp.SendMail(
      config.Host + ":" + strconv.Itoa(config.Port),
      auth,
      from,
      to,
      msg,
    )
    if err != nil {
      revel.ERROR.Println(err)
    }
  }
}

func init() {
  revel.OnAppStart(func() {
    config.Type= revel.Config.StringDefault("mail.type", "file")
    config.Host = revel.Config.StringDefault("mail.host", "127.0.0.1")
    config.Port = revel.Config.IntDefault("mail.port", 25)
    config.User = revel.Config.StringDefault("mail.user", "")
    config.Password = revel.Config.StringDefault("mail.password", "")
    if config.User != "" {
      auth = smtp.PlainAuth(
        "",
        config.User,
        config.Password,
        config.Host,
      )
    }
  })
}
