package mailer

import (
    "bytes"
    "github.com/orc/utils"
    "net/smtp"
    "strconv"
    "text/template"
)

const HASH_SIZE = 32
const Server = "http://secret-oasis-3805.herokuapp.com"

var err error

type Admin struct {
    Name   string

    EmailAdmin  string
    Password    string

    SMTPServer string
    Port        int
}

type SmtpTemplateData struct {
    From            string
    To              string
    Subject         string
    ConfirmationUrl string
    RejactionUrl    string
}

func SendConfirmEmail(to, address, token string) {
    admin := &Admin{
        Name:       "Secret Oasis",
        EmailAdmin: "secret.oasis.3805@gmail.com",
        Password:   "mysterious-reef-6215",
        SMTPServer: "smtp.gmail.com",
        Port:       587}

    auth := smtp.PlainAuth(
        "",
        admin.EmailAdmin,
        admin.Password,
        admin.SMTPServer)

    var emailTemplate = `From: {{.From}}
To: {{.To}}
Subject: {{.Subject}}

Здравствуйте!

Спасибо за использование нашего ресурса secret-oasis-3805.com!
Для подтверждения вашей учетной записи, пожалуйста, перейдите по ссылке: {{ .ConfirmationUrl }}

Если это письмо попало к Вам по ошибке, то, чтобы больше не получать писем от ` + admin.Name + `, перейдите по этой ссылке: {{ .RejactionUrl }}`

    context := &SmtpTemplateData{
        admin.Name,
        to,
        "Подтверждение регистрации",
        Server+"/handler/confirmuser/"+token,
        Server+"/handler/rejectuser/"+token}

    t, err := template.New("confirmationmail").Parse(emailTemplate)
    if utils.HandleErr("[SendEmail] Error trying to parse mail template: ", err, nil) {
        return
    }

    var doc bytes.Buffer
    err = t.Execute(&doc, context)
    if utils.HandleErr("[SendEmail] Error trying to execute mail template: ", err, nil) {
        return
    }

    err = smtp.SendMail(
        admin.SMTPServer+":"+strconv.Itoa(admin.Port),
        auth,
        admin.EmailAdmin,
        []string{address},
        doc.Bytes())

    if utils.HandleErr("[SendEmail] Error attempting to send a mail: ", err, nil) {
        return
    }
}
