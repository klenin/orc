package db

const user string = "admin"
const dbname string = "orc"
const password string = "admin"

var Tables = []string{
    "events",
    "event_types",
    "events_types",
    //"teams",
    "persons",
    "persons_events",
    "users",
    //"teams_persons",
    "forms",
    "param_types",
    "params",
    "forms_types",
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
    //"Команды",
    "Персоны",
    "Персоны-Мероприятия",
    "Пользователи",
    //"Команды-Персоны",
    "Формы",
    "Типы параметров",
    "Параметры",
    "Формы-Типы мероприятий",
    "Значения параметров",
    "Лица",
    "Регистрации",
    "Регистрация-Мероприятие-Тип мероприятия-Значение параметра",
    "Мероприятия-Регистрации",
}

//var Tables = []map[string]string{
//    {"name": "events", "ru-name": "Мероприятия"},
//    {"name": "event_types", "ru-name": "Типы мероприятий"},
//    {"name": "events_types", "ru-name": "Мероприятия-Типы"},
//    //"teams",
//    {"name": "persons", "ru-name": "Персоны"},
//    {"name": "persons_events", "ru-name": "Персоны-Мероприятия"},
//    {"name": "users", "ru-name": "Пользователи"},
//    //"teams_persons",
//    {"name": "forms", "ru-name": "Формы"},
//    {"name": "params", "ru-name": "Параметры"},
//    {"name": "forms_types", "ru-name": "Формы-Типы мероприятий"},
//    {"name": "param_values", "ru-name": "Значения параметров"},
//}
