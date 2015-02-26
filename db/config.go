package db

const user string = "admin"
const dbname string = "orc"
const password string = "admin"

var Tables = []string{
    "events",
    "event_types",
    "events_types",
    "persons",
    "users",
    "forms",
    "param_types",
    "params",
    "events_forms",
    "param_values",
    "faces",
    "registrations",
    "reg_param_vals",
    "events_regs",
}

var TableNames = []string{
    "Мероприятия",
    "Типы мероприятий",
    "Мероприятия-Типы",
    "Персоны",
    "Пользователи",
    "Формы",
    "Типы параметров",
    "Параметры",
    "Мероприятия-Формы",
    "Значения параметров",
    "Лица",
    "Регистрации",
    "Регистрация-Мероприятие-Значение параметра",
    "Мероприятия-Регистрации",
}
