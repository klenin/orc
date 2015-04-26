package mailer

var ComfirmRegistrationEmailTmp = `From: {{ .From }}
To: {{ .To }}
Subject: {{ .Subject }}
Content-Type: text/html; charset=UTF-8

<p>
    Здравствуйте, <b>{{ .To }}</b>.
</p>

<p>
    Спасибо за использование ресурса <a href="`+Server+`">Системы учета участников мероприятий</a>.<br/>
    Для подтверждения вашей учетной записи, пожалуйста, перейдите по ссылке: <a href="{{ .ConfirmationUrl }}">подтвердить</a>.<br/>
    Если это письмо попало к Вам по ошибке, то перейдите по ссылке: <a href="{{ .RejectionUrl }}">отклонить</a>.
</p>`

var RejectRequestTmp = `From: {{ .From }}
To: {{ .To }}
Subject: {{ .Subject }}
Content-Type: text/html; charset=UTF-8

<p>
    Здравствуйте, <b>{{ .To }}</b>.
</p>

<p>
    Вы отправили заявку на участие в мероприятии "{{ .EventName }}", но указанные Вами данные имеют некоторые неточности.
    Пожалуйста, заполните заявку еще раз.
</p>`

var ConfirmRequestTmp = `From: {{ .From }}
To: {{ .To }}
Subject: {{ .Subject }}
Content-Type: text/html; charset=UTF-8

<p>
    Здравствуйте, <b>{{ .To }}</b>.
</p>

<p>
    Ваша заявка на участие в мероприятии "{{ .EventName }}" принята.
</p>`

var InviteToGroupEmailTmp = `From: {{ .From }}
To: {{ .To }}
Subject: {{ .Subject }}
Content-Type: text/html; charset=UTF-8

<p>
    Здравствуйте, <b>{{ .To }}</b>.
</p>

<p>
    <b>{{ .HeadName }}</b> хочет добавить Вас в группу "{{ .GroupName }}".<br/>
    Вам следует залогиниться (зарегистироваться) в <a href="`+Server+`">Системе учета участников мероприятий</a>.<br/>
    Затем для того, чтобы присоединиться к группе "{{ .GroupName }}", пройдите по ссылке: <a href="{{ .ConfirmationUrl }}">присоединиться к группе</a>.<br/>
    Чтобы отклонить приглашение, пройдите по ссылке: <a href="{{ .RejectionUrl }}">отклонить приглашение</a>.
</p>`

var AttendAnEventEmailTmp = `From: {{ .From }}
To: {{ .To }}
Subject: {{ .Subject }}
Content-Type: text/html; charset=UTF-8

<p>
    Здравствуйте, <b>{{ .To }}</b>.
</p>

<p>
    Вы участвуете в мероприятии "{{ .EventName }}".<br/>
    Пожалуйста, заполните анкету в личном кабинете <a href="`+Server+`">Системы учета участников мероприятий</a>.
</p>`
