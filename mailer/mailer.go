package mailer

import (
    "bytes"
    "github.com/klenin/orc/utils"
    "log"
    "net/smtp"
    "strconv"
    "text/template"
)

const HASH_SIZE = 32
const Server = "https://server/link/"

var err error

type Admin struct {
    Name       string

    Email      string
    Password   string

    SMTPServer string
    Port       int
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

var Admin_ = &Admin{
    Name:       "Name of Admin",
    Email:      "Email of admin",
    Password:   "Password of email",
    SMTPServer: "smtp.gmail.com",
    Port:       587}

var auth = smtp.PlainAuth(
    "",
    Admin_.Email,
    Admin_.Password,
    Admin_.SMTPServer)

func SendEmail(address, tmp string, context *SmtpTemplateData) bool {
    var doc bytes.Buffer
    template.Must(template.New("email").Parse(tmp)).Execute(&doc, context)

    err = smtp.SendMail(
        Admin_.SMTPServer+":"+strconv.Itoa(Admin_.Port),
        auth,
        Admin_.Email,
        []string{address},
        doc.Bytes())

    return !utils.HandleErr("[SendEmail] Error attempting to send a mail: ", err, nil)
}

func SendConfirmEmail(to, address, token string) bool {
    log.Println("SendConfirmEmail: address: ", address)
    log.Println("SendConfirmEmail: to: ", to)

    context := &SmtpTemplateData{
        From: Admin_.Name,
        To: to,
        Subject: "Подтверждение регистрации",
        ConfirmationUrl: Server+"/registrationcontroller/confirmuser/"+token,
        RejectionUrl: Server+"/registrationcontroller/rejectuser/"+token}

    return SendEmail(address, ComfirmRegistrationEmailTmp, context)
}

func SendEmailToConfirmRejectPersonRequest(to, address, event string, confirm bool) bool {
    log.Println("SendEmailToConfirmRejectPersonRequest: address: ", address)
    log.Println("SendEmailToConfirmRejectPersonRequest: to: ", to)

    var emailTemplate string

    context := &SmtpTemplateData{
        From: Admin_.Name,
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
        From: Admin_.Name,
        To: to,
        Subject: `Приглашение в группу "`+groupName+`"`,
        ConfirmationUrl: Server+"/groupcontroller/confirminvitationtogroup/"+token,
        RejectionUrl: Server+"/groupcontroller/rejectinvitationtogroup/"+token,
        HeadName: headName,
        GroupName: groupName}

    return SendEmail(address, InviteToGroupEmailTmp, context)
}

func AttendAnEvent(to, address, eventName, groupName string) bool {
    log.Println("AttendAnEvent: address: ", address)
    log.Println("AttendAnEvent: to: ", to)

    context := &SmtpTemplateData{
        From: Admin_.Name,
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
        From: Admin_.Name,
        To: to,
        Subject: `Система учета учатников мероприятий`,
        ConfirmationUrl: Server+"/wellcometoprofile/"+token,}

    return SendEmail(address, WellcomeToProfileEmailTmp, context)
}
