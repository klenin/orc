package controllers

import (
    "github.com/lib/pq"
    "github.com/klenin/orc/db"
    "github.com/klenin/orc/mailer"
    "github.com/klenin/orc/utils"
    "log"
    "net/http"
    "strconv"
    "strings"
    "time"
)

func (*BaseController) GroupController() *GroupController {
    return new(GroupController)
}

type GroupController struct {
    Controller
}

func (this *GroupController) Register() {
    userId, err := this.CheckSid()
    if err != nil {
        http.Redirect(this.Response, this.Request, "/", http.StatusUnauthorized)
        return
    }

    request, err := utils.ParseJS(this.Request, this.Response)
    if err != nil {
        utils.SendJSReply(map[string]interface{}{"result": err.Error()}, this.Response)
    }

    groupId, err := strconv.Atoi(request["group_id"].(string))
    if err != nil {
        utils.SendJSReply(map[string]interface{}{"result": err.Error()}, this.Response)
    }

    eventId, err := strconv.Atoi(request["event_id"].(string))
    if err != nil {
        utils.SendJSReply(map[string]interface{}{"result": err.Error()}, this.Response)
    }

    var eventName string; var teamEvent bool
    if err = this.GetModel("events").
        LoadWherePart(map[string]interface{}{"id": eventId}).
        SelectRow([]string{"name", "team"}).
        Scan(&eventName, &teamEvent);
        err != nil {
        utils.SendJSReply(map[string]interface{}{"result": err.Error()}, this.Response)
        return
    }

    faceId, groupName := 1, ""
    query := `SELECT groups.face_id, groups.name FROM groups
        INNER JOIN faces ON faces.id = groups.face_id
        INNER JOIN users ON users.id = faces.user_id
        WHERE users.id = $1 AND groups.id = $2;`
    err = db.QueryRow(query, []interface{}{userId, groupId}).Scan(&faceId, &groupName)

    if (err != nil || faceId == 1 || groupName == "") && !this.isAdmin() {
        utils.SendJSReply(map[string]interface{}{"result": "Вы не являетесь владельцем группы"}, this.Response)
        return
    }

    if db.IsExists("group_registrations", []string{"group_id", "event_id"}, []interface{}{groupId, eventId}) {
        utils.SendJSReply(map[string]interface{}{"result": "Группа уже зарегистрированна в этом мероприятии"}, this.Response)
        return
    }

    var groupregId int
    this.GetModel("group_registrations").
        LoadModelData(map[string]interface{}{"event_id": eventId, "group_id": groupId, "status": false}).
        QueryInsert("RETURNING id").
        Scan(&groupregId)

    query = `SELECT persons.status, persons.group_id, f.id as face_id,
        array_to_string(
        array(
            SELECT param_values.value
            FROM param_values
            INNER JOIN registrations ON registrations.id = param_values.reg_id
            INNER JOIN faces ON faces.id = registrations.face_id
            INNER JOIN events ON events.id = registrations.event_id
            INNER JOIN params ON params.id = param_values.param_id
            WHERE param_values.param_id IN (5, 6, 7) AND events.id = 1 AND faces.id = f.id ORDER BY param_values.param_id
        ), ' ') as name,

        (SELECT param_values.value
            FROM param_values
            INNER JOIN registrations ON registrations.id = param_values.reg_id
            INNER JOIN faces ON faces.id = registrations.face_id
            INNER JOIN events ON events.id = registrations.event_id
            INNER JOIN params ON params.id = param_values.param_id
            WHERE param_values.param_id = 4 AND events.id = 1 AND faces.id = f.id
        ) as email

        FROM persons
        INNER JOIN groups ON groups.id = persons.group_id
        INNER JOIN faces as f ON f.id = persons.face_id
        WHERE groups.id = $1;`
    data := db.Query(query, []interface{}{groupId})

    query = `SELECT params.id FROM events_forms
        INNER JOIN events ON events.id = events_forms.event_id
        INNER JOIN forms ON forms.id = events_forms.form_id
        INNER JOIN params ON forms.id = params.form_id
        WHERE events.id = $1 AND forms.personal = true ORDER BY forms.id;`
    params := db.Query(query, []interface{}{eventId})

    date := time.Now().Format("20060102T15:04:05Z00:00")

    for _, v := range data {
        status := v.(map[string]interface{})["status"].(bool)
        personFaceId := v.(map[string]interface{})["face_id"].(int)
        var personUserId int
        this.GetModel("faces").
        LoadWherePart(map[string]interface{}{"id": personFaceId}).
        SelectRow([]string{"user_id"}).
        Scan(&personUserId)

        if !status {
            continue
        }

        regId := this.regExists(personUserId, eventId)
        if regId == -1 {
            this.GetModel("registrations").
                LoadModelData(map[string]interface{}{
                    "face_id": personFaceId,
                    "event_id": eventId,
                    "status": false}).
                QueryInsert("RETURNING id").
                Scan(&regId)

            for _, elem := range params {
                paramId := int(elem.(map[string]interface{})["id"].(int))
                this.GetModel("param_values").
                    LoadModelData(map[string]interface{}{
                        "param_id": paramId,
                        "value": " ",
                        "date": date,
                        "user_id": userId,
                        "reg_id": regId}).
                    QueryInsert("").
                    Scan()
            }
        }

        this.GetModel("regs_groupregs").
            LoadModelData(map[string]interface{}{"groupreg_id": groupregId, "reg_id": regId}).
            QueryInsert("").
            Scan()

        to := v.(map[string]interface{})["name"].(string)
        address := v.(map[string]interface{})["email"].(string)
        if !mailer.AttendAnEvent(to, address, eventName, groupName) {
            utils.SendJSReply(map[string]interface{}{"result": "Ошибка. Письмо с уведомлением не отправлено."}, this.Response)
        }
    }

    if teamEvent == true {
        query = `SELECT params.id FROM events_forms
            INNER JOIN events ON events.id = events_forms.event_id
            INNER JOIN forms ON forms.id = events_forms.form_id
            INNER JOIN params ON forms.id = params.form_id
            WHERE events.id = $1 AND forms.personal = false ORDER BY forms.id;`
        params := db.Query(query, []interface{}{eventId})

        var regId int
        this.GetModel("registrations").
            LoadModelData(map[string]interface{}{"face_id": faceId, "event_id": eventId, "status": false}).
            QueryInsert("RETURNING id").
            Scan(&regId)

        for _, elem := range params {
            this.GetModel("param_values").
                LoadModelData(map[string]interface{}{
                    "param_id": int(elem.(map[string]interface{})["id"].(int)),
                    "value": " ",
                    "date": date,
                    "user_id": userId,
                    "reg_id": regId}).
                QueryInsert("").
                Scan()
        }

        this.GetModel("regs_groupregs").
            LoadModelData(map[string]interface{}{"groupreg_id": groupregId, "reg_id": regId}).
            QueryInsert("").
            Scan()
    }

    utils.SendJSReply(map[string]interface{}{"result": "ok"}, this.Response)
}

func (this *GroupController) ConfirmInvitationToGroup(token string) {
    var faceId int
    if err := this.GetModel("persons").
        LoadWherePart(map[string]interface{}{"token": token}).
        SelectRow([]string{"face_id"}).
        Scan(&faceId);
        err != nil {
        if this.Response != nil {
            this.Render([]string{"mvc/views/msg.html"}, "msg", "Неверный токен.")
        }
        return
    }

    params := map[string]interface{}{"status": true, "token": " "}
    where := map[string]interface{}{"token": token}
    this.GetModel("persons").Update(this.isAdmin(), -1, params, where)

    if this.Response != nil {
        this.Render([]string{"mvc/views/msg.html"}, "msg", "Вы успешно присоединены к группе.")
    }
}

func (this *GroupController) RejectInvitationToGroup(token string) {
    db.Exec("DELETE FROM persons WHERE token = $1;", []interface{}{token})

    if this.Response != nil {
        this.Render([]string{"mvc/views/msg.html"}, "msg", "Запрос о присоединении к группе успешно отклонен.")
    }
}

func (this *GroupController) IsRegGroup() {
    _, err := this.CheckSid()
    if err != nil {
        http.Redirect(this.Response, this.Request, "/", http.StatusUnauthorized)
        return
    }

    request, err := utils.ParseJS(this.Request, this.Response)
    if err != nil {
        utils.SendJSReply(map[string]interface{}{"result": err.Error()}, this.Response)
    }

    groupId, err := strconv.Atoi(request["group_id"].(string))
    if err != nil {
        utils.SendJSReply(map[string]interface{}{"result": err.Error()}, this.Response)
    }

    addDelFlag := !db.IsExists("group_registrations", []string{"group_id"}, []interface{}{groupId})
    utils.SendJSReply(map[string]interface{}{"result": "ok", "addDelFlag": addDelFlag}, this.Response)
}

func (this *GroupController) AddPerson() {
    userId, err := this.CheckSid()
    if err != nil {
        utils.SendJSReply(map[string]interface{}{"result": "Unauthorized"}, this.Response)
        return
    }

    request, err := utils.ParseJS(this.Request, this.Response)
    if err != nil {
        utils.SendJSReply(map[string]interface{}{"result": err.Error()}, this.Response)
        return
    }

    groupId, err := strconv.Atoi(request["group_id"].(string))
    if err != nil {
        utils.SendJSReply(map[string]interface{}{"result": err.Error()}, this.Response)
        return
    }

    var groupName string
    db.QueryRow("SELECT name FROM groups WHERE id = $1;", []interface{}{groupId}).Scan(&groupName)

    date := time.Now().Format("2006-01-02T15:04:05Z00:00")
    token := utils.GetRandSeq(HASH_SIZE)
    to, address, headName := "", "", ""

    query := `SELECT param_values.value
        FROM param_values
        INNER JOIN registrations ON registrations.id = param_values.reg_id
        INNER JOIN params ON params.id = param_values.param_id
        INNER JOIN events ON events.id = registrations.event_id
        INNER JOIN faces ON faces.id = registrations.face_id
        INNER JOIN users ON users.id = faces.user_id
        WHERE params.id in (5, 6, 7) AND users.id = $1 AND events.id = 1 ORDER BY params.id;`
    data := db.Query(query, []interface{}{userId})

    if len(data) < 3 {
        utils.SendJSReply(map[string]interface{}{"result": "Данные о руководителе группы отсутсвуют"}, this.Response)
        return

    } else {
        headName = data[0].(map[string]interface{})["value"].(string)
        headName += " " + data[1].(map[string]interface{})["value"].(string)
        headName += " " + data[2].(map[string]interface{})["value"].(string)
    }

    var faceId int
    this.GetModel("faces").QueryInsert("RETURNING id").Scan(&faceId)

    this.GetModel("persons").
        LoadModelData(map[string]interface{}{"face_id": faceId, "group_id": groupId, "status": false, "token": token}).
        QueryInsert("").
        Scan()

    var regId int
    this.GetModel("registrations").
        LoadModelData(map[string]interface{}{"face_id": faceId, "event_id": 1, "status": false}).
        QueryInsert("RETURNING id").
        Scan(&regId)

    var paramValueIds []string

    for _, element := range request["data"].([]interface{}) {
        paramId, err := strconv.Atoi(element.(map[string]interface{})["id"].(string))
        if err != nil {
            log.Println(err.Error())
            continue
        }

        query := `SELECT params.name FROM params WHERE params.id = $1;`
        res := db.Query(query, []interface{}{paramId})

        name := res[0].(map[string]interface{})["name"].(string)
        value := element.(map[string]interface{})["value"].(string)

        if utils.MatchRegexp("^[ \t\v\r\n\f]{0,}$", value) {
            db.QueryDeleteByIds("param_vals", strings.Join(paramValueIds, ", "))
            db.QueryDeleteByIds("registrations", strconv.Itoa(regId))
            db.QueryDeleteByIds("faces", strconv.Itoa(faceId))
            utils.SendJSReply(map[string]interface{}{"result": "Заполните параметр '"+name+"'."}, this.Response)
            return
        }

        var paramValId int
        paramValues := this.GetModel("param_values")
        err = paramValues.LoadModelData(map[string]interface{}{
            "param_id": paramId,
            "value": value,
            "date": date,
            "user_id": userId,
            "reg_id": regId}).
            QueryInsert("RETURNING id").
            Scan(&paramValId)
        if err, ok := err.(*pq.Error); ok {
            log.Println(err.Code.Name())
        }

        paramValueIds = append(paramValueIds, strconv.Itoa(paramValId))

        if paramId == 4 {
            address = value
        } else if paramId == 5 || paramId == 6 || paramId == 7 {
            to += value + " "
        }
    }

    if !mailer.InviteToGroup(to, address, token, headName, groupName) {
        utils.SendJSReply(
            map[string]interface{}{
                "result": "Вы указали неправильный email, отправить письмо-приглашенине невозможно"},
            this.Response)
        return
    }

    utils.SendJSReply(map[string]interface{}{"result": "ok"}, this.Response)
}
