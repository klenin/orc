package mailer

import (
    "bytes"
    "github.com/orc/utils"
    "log"
    "net/smtp"
    "strconv"
    "text/template"
)

const HASH_SIZE = 32
const Server = "secret-oasis-3805.herokuapp.com"

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
    RejectionUrl    string
    EventName       string
}

var admin = &Admin{
    Name:       "Secret Oasis",
    EmailAdmin: "secret.oasis.3805@gmail.com",
    Password:   "mysterious-reef-6215",
    SMTPServer: "smtp.gmail.com",
    Port:       587}

var auth = smtp.PlainAuth(
    "",
    admin.EmailAdmin,
    admin.Password,
    admin.SMTPServer)

var comfirmRegistrationEmailTmp = `From: {{.From}}
To: {{.To}}
Subject: {{.Subject}}

Здравствуйте!

Спасибо за использование нашего ресурса secret-oasis-3805.com!
Для подтверждения вашей учетной записи, пожалуйста, перейдите по ссылке: {{ .ConfirmationUrl }}

Если это письмо попало к Вам по ошибке, то, чтобы больше не получать писем от ` + admin.Name + `, перейдите по этой ссылке: {{ .RejectionUrl }}`

var rejectRequestTmp = `From: {{.From}}
To: {{.To}}
Subject: {{.Subject}}

Здравствуйте!

Спасибо за использование нашего ресурса secret-oasis-3805.com!
Вы отправили заявку на участие в мероприятии "{{ .EventName }}", но указанные Вами данные имеют некоторые неточности.
Пожалуйста, заполните заявку еще раз.`

var confirmRequestTmp = `From: {{.From}}
To: {{.To}}
Subject: {{.Subject}}

Здравствуйте!

Спасибо за использование нашего ресурса secret-oasis-3805.com!
Ваша заявка на участие в мероприятии "{{ .EventName }}" принята.`

func SendConfirmEmail(to, address, token string) {

    log.Println("SendConfirmEmail: address: ", address)
    log.Println("SendConfirmEmail: to: ", to)

    context := &SmtpTemplateData{
        From: admin.Name,
        To: to,
        Subject: "Подтверждение регистрации",
        ConfirmationUrl: Server+"/handler/confirmuser/"+token,
        RejectionUrl: Server+"/handler/rejectuser/"+token}

    t, err := template.New("confirmationmail").Parse(comfirmRegistrationEmailTmp)
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

    utils.HandleErr("[SendEmail] Error attempting to send a mail: ", err, nil)
}

func SendEmailToConfirmRejectPersonRequest(to, address, event string, confirm bool) {

    var emailTemplate string
    if !confirm {
        emailTemplate = rejectRequestTmp
    } else {
        emailTemplate = confirmRequestTmp
    }

    context := &SmtpTemplateData{
        From: admin.Name,
        To: to,
        Subject: `Подтверждение заявки на участие в мероприятии "`+event+`"`,
        EventName: event}

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
