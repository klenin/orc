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
    EventUrl        string
    HeadName        string
    GroupName       string
    Login           string
    Password        string
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

var comfirmRegistrationEmailTmp = `From: {{ .From }}
To: {{ .To }}
Subject: {{ .Subject }}

Здравствуйте!

Спасибо за использование нашего ресурса secret-oasis-3805.com!
Для подтверждения вашей учетной записи, пожалуйста, перейдите по ссылке: {{ .ConfirmationUrl }}

Если это письмо попало к Вам по ошибке, то, чтобы больше не получать писем от ` + admin.Name + `, перейдите по этой ссылке: {{ .RejectionUrl }}`

var rejectRequestTmp = `From: {{ .From }}
To: {{ .To }}
Subject: {{ .Subject }}

Здравствуйте!

Вы отправили заявку на участие в мероприятии "{{ .EventName }}", но указанные Вами данные имеют некоторые неточности.
Пожалуйста, заполните заявку еще раз.`

var confirmRequestTmp = `From: {{ .From }}
To: {{ .To }}
Subject: {{ .Subject }}

Здравствуйте!

Ваша заявка на участие в мероприятии "{{ .EventName }}" принята.`

var inviteToGroupEmailTmp = `From: {{ .From }}
To: {{ .To }}
Subject: {{ .Subject }}

Здравствуйте, {{ .To }}!

{{ .HeadName }} хочет добавить Вас в группу "{{ .GroupName }}".

Вы ДОЛЖНЫ залогиниться (зарегистироваться) в системе `+Server+`.

Затем для того, чтобы присоединиться к группе "{{ .GroupName }}", пройдите по ссылке: {{ .ConfirmationUrl }}

Чтобы отклонить приглашение, пройдите по ссылке: {{ .RejectionUrl }}`

var attendAnEventEmailTmp = `From: {{ .From }}
To: {{ .To }}
Subject: {{ .Subject }}

Здравствуйте, {{ .To }}!

Вы участвуете в мероприятии "{{ .EventName }}".
Пожалуйста, заполните анкету в личном кабинете `+Server

func SendConfirmEmail(to, address, token string) bool {

    log.Println("SendConfirmEmail: address: ", address)
    log.Println("SendConfirmEmail: to: ", to)

    context := &SmtpTemplateData{
        From: admin.Name,
        To: to,
        Subject: "Подтверждение регистрации",
        ConfirmationUrl: Server+"/handler/confirmuser/"+token,
        RejectionUrl: Server+"/handler/rejectuser/"+token}

    t, err := template.New("confirmationmail").Parse(comfirmRegistrationEmailTmp)
    if utils.HandleErr("[SendConfirmEmail] Error trying to parse mail template: ", err, nil) {
        return false
    }

    var doc bytes.Buffer
    err = t.Execute(&doc, context)
    if utils.HandleErr("[SendConfirmEmail] Error trying to execute mail template: ", err, nil) {
        return false
    }

    err = smtp.SendMail(
        admin.SMTPServer+":"+strconv.Itoa(admin.Port),
        auth,
        admin.EmailAdmin,
        []string{address},
        doc.Bytes())

    return !utils.HandleErr("[SendConfirmEmail] Error attempting to send a mail: ", err, nil)
}

func SendEmailToConfirmRejectPersonRequest(to, address, event string, confirm bool) bool {

    var emailTemplate string

    context := &SmtpTemplateData{
        From: admin.Name,
        To: to,
        Subject: `Подтверждение заявки на участие в мероприятии "`+event+`"`,
        EventName: event}

    if !confirm {
        emailTemplate = rejectRequestTmp
        context.Subject = `Отклонение заявки на участие в мероприятии "`+event+`"`
    } else {
        emailTemplate = confirmRequestTmp
    }

    t, err := template.New("confirmationmail").Parse(emailTemplate)
    if utils.HandleErr("[SendEmailToConfirmRejectPersonRequest] Error trying to parse mail template: ", err, nil) {
        return false
    }

    var doc bytes.Buffer
    err = t.Execute(&doc, context)
    if utils.HandleErr("[SendEmailToConfirmRejectPersonRequest] Error trying to execute mail template: ", err, nil) {
        return false
    }

    err = smtp.SendMail(
        admin.SMTPServer+":"+strconv.Itoa(admin.Port),
        auth,
        admin.EmailAdmin,
        []string{address},
        doc.Bytes())

    return !utils.HandleErr("[SendEmailToConfirmRejectPersonRequest] Error attempting to send a mail: ", err, nil)
}

func InviteToGroup(to, address, token, headName, groupName string) bool {

    log.Println("SendConfirmEmail: address: ", address)
    log.Println("SendConfirmEmail: to: ", to)

    context := &SmtpTemplateData{
        From: admin.Name,
        To: to,
        Subject: `Приглашение в группу "`+groupName+`"`,
        ConfirmationUrl: Server+"/handler/confirminvitationtogroup/"+token,
        RejectionUrl: Server+"/handler/rejectinvitationtogroup/"+token,
        HeadName: headName,
        GroupName: groupName}

    t, err := template.New("mail").Parse(inviteToGroupEmailTmp)
    if utils.HandleErr("[InviteToGroup] Error trying to parse mail template: ", err, nil) {
        return false
    }

    var doc bytes.Buffer
    err = t.Execute(&doc, context)
    if utils.HandleErr("[InviteToGroup] Error trying to execute mail template: ", err, nil) {
        return false
    }

    err = smtp.SendMail(
        admin.SMTPServer+":"+strconv.Itoa(admin.Port),
        auth,
        admin.EmailAdmin,
        []string{address},
        doc.Bytes())

    return !utils.HandleErr("[InviteToGroup] Error attempting to send a mail: ", err, nil)
}

func AttendAnEvent(to, address, eventName, groupName string) bool {

    log.Println("SendConfirmEmail: address: ", address)
    log.Println("SendConfirmEmail: to: ", to)

    context := &SmtpTemplateData{
        From: admin.Name,
        To: to,
        Subject: `Уведомление об участии в мероприятии "`+eventName+`"`,
        GroupName: groupName,
        EventName: eventName}

    t, err := template.New("mail").Parse(attendAnEventEmailTmp)
    if utils.HandleErr("[InviteToGroup] Error trying to parse mail template: ", err, nil) {
        return false
    }

    var doc bytes.Buffer
    err = t.Execute(&doc, context)
    if utils.HandleErr("[InviteToGroup] Error trying to execute mail template: ", err, nil) {
        return false
    }

    err = smtp.SendMail(
        admin.SMTPServer+":"+strconv.Itoa(admin.Port),
        auth,
        admin.EmailAdmin,
        []string{address},
        doc.Bytes())

    return !utils.HandleErr("[InviteToGroup] Error attempting to send a mail: ", err, nil)
}
