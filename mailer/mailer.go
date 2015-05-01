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

func SendEmail(address, tmp string, context *SmtpTemplateData) bool {
    var doc bytes.Buffer
    template.Must(template.New("email").Parse(tmp)).Execute(&doc, context)

    err = smtp.SendMail(
        admin.SMTPServer+":"+strconv.Itoa(admin.Port),
        auth,
        admin.EmailAdmin,
        []string{address},
        doc.Bytes())

    return !utils.HandleErr("[SendEmail] Error attempting to send a mail: ", err, nil)
}

func SendConfirmEmail(to, address, token string) bool {
    log.Println("SendConfirmEmail: address: ", address)
    log.Println("SendConfirmEmail: to: ", to)

    context := &SmtpTemplateData{
        From: admin.Name,
        To: to,
        Subject: "Подтверждение регистрации",
        ConfirmationUrl: Server+"/handler/confirmuser/"+token,
        RejectionUrl: Server+"/handler/rejectuser/"+token}

    return SendEmail(address, ComfirmRegistrationEmailTmp, context)
}

func SendEmailToConfirmRejectPersonRequest(to, address, event string, confirm bool) bool {
    log.Println("SendEmailToConfirmRejectPersonRequest: address: ", address)
    log.Println("SendEmailToConfirmRejectPersonRequest: to: ", to)

    var emailTemplate string

    context := &SmtpTemplateData{
        From: admin.Name,
        To: to,
        Subject: `Подтверждение заявки на участие в мероприятии "`+event+`"`,
        EventName: event}

    if !confirm {
        emailTemplate = RejectRequestTmp
        context.Subject = `Отклонение заявки на участие в мероприятии "`+event+`"`
    } else {
        emailTemplate = ConfirmRequestTmp
    }

    return SendEmail(address, emailTemplate, context)
}

func InviteToGroup(to, address, token, headName, groupName string) bool {
    log.Println("InviteToGroup: address: ", address)
    log.Println("InviteToGroup: to: ", to)

    context := &SmtpTemplateData{
        From: admin.Name,
        To: to,
        Subject: `Приглашение в группу "`+groupName+`"`,
        ConfirmationUrl: Server+"/handler/confirminvitationtogroup/"+token,
        RejectionUrl: Server+"/handler/rejectinvitationtogroup/"+token,
        HeadName: headName,
        GroupName: groupName}

    return SendEmail(address, InviteToGroupEmailTmp, context)
}

func AttendAnEvent(to, address, eventName, groupName string) bool {
    log.Println("AttendAnEvent: address: ", address)
    log.Println("AttendAnEvent: to: ", to)

    context := &SmtpTemplateData{
        From: admin.Name,
        To: to,
        Subject: `Уведомление об участии в мероприятии "`+eventName+`"`,
        GroupName: groupName,
        EventName: eventName}

    return SendEmail(address, AttendAnEventEmailTmp, context)
}

func SendEmailWellcomeToProfile(to, address, token string) bool {
    log.Println("SendEmailWellcomeToProfile: address: ", address)
    log.Println("SendEmailWellcomeToProfile: to: ", to)

    context := &SmtpTemplateData{
        From: admin.Name,
        To: to,
        Subject: `Система учета учатников мероприятий`,
        ConfirmationUrl: Server+"/handler/wellcometoprofile/"+token,}

    return SendEmail(address, WellcomeToProfileEmailTmp, context)
}
