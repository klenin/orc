package db

var Events = `CREATE TABLE IF NOT EXISTS events (
    id         int  NOT NULL PRIMARY KEY DEFAULT NEXTVAL('events_id_seq'),
    name       text NOT NULL UNIQUE,
    date_start date NOT NULL,
    date_end   date NOT NULL,
	time       time NOT NULL,
	url        text
)`

var Event_types = `CREATE TABLE IF NOT EXISTS event_types (
    id          int     NOT NULL PRIMARY KEY DEFAULT NEXTVAL('event_types_id_seq'),
    name        text    NOT NULL UNIQUE,
	description text    NOT NULL,
	topicality  boolean NOT NULL
);`

var Events_types = `CREATE TABLE IF NOT EXISTS events_types (
    id       int  NOT NULL PRIMARY KEY DEFAULT NEXTVAL('events_types_id_seq'),
    event_id int  NOT NULL REFERENCES events(id) ON DELETE CASCADE,
    type_id  int  NOT NULL REFERENCES event_types(id) ON DELETE CASCADE
);`

var Teams = `CREATE TABLE IF NOT EXISTS teams (
    id   int  NOT NULL PRIMARY KEY DEFAULT NEXTVAL('teams_id_seq'),
    name text NOT NULL UNIQUE
);`

var Persons = `CREATE TABLE IF NOT EXISTS persons (
    id            int  NOT NULL PRIMARY KEY DEFAULT NEXTVAL('persons_id_seq'),
    fname         text NOT NULL,
    lname         text NOT NULL,
    pname         text NOT NULL
);`

var Users = `CREATE TABLE IF NOT EXISTS users (
    id        int  NOT NULL PRIMARY KEY DEFAULT NEXTVAL('users_id_seq'),
    login     text NOT NULL,
    pass      text NOT NULL,
    salt      text NOT NULL,
	role      text NOT NULL DEFAULT 'user',
	person_id int  NOT NULL DEFAULT '-1' REFERENCES persons(id) ON DELETE CASCADE
);`

var Teams_persons = `CREATE TABLE IF NOT EXISTS teams_persons (
    id        int NOT NULL PRIMARY KEY DEFAULT NEXTVAL('teams_persons_id_seq'),
    team_id   int NOT NULL REFERENCES teams(id) ON DELETE CASCADE,
    person_id int NOT NULL REFERENCES persons(id) ON DELETE CASCADE
);`

var Teams_users = `CREATE TABLE IF NOT EXISTS teams_users (
    id      int NOT NULL PRIMARY KEY DEFAULT NEXTVAL('teams_users_id_seq'),
    team_id int NOT NULL REFERENCES teams(id) ON DELETE CASCADE,
    user_id int NOT NULL REFERENCES users(id) ON DELETE CASCADE
);`

/*var Persons_users = `CREATE TABLE IF NOT EXISTS persons_users (
    id        int NOT NULL PRIMARY KEY DEFAULT NEXTVAL('persons_users_id_seq'),
    person_id int NOT NULL REFERENCES persons(id) ON DELETE CASCADE,
    user_id   int NOT NULL REFERENCES users(id) ON DELETE CASCADE
);`*/

var Forms = `CREATE TABLE IF NOT EXISTS forms (
    id   int  NOT NULL PRIMARY KEY DEFAULT NEXTVAL('forms_id_seq'),
    name text NOT NULL
);`

var Params = `CREATE TABLE IF NOT EXISTS params (
    id         int  NOT NULL PRIMARY KEY DEFAULT NEXTVAL('params_id_seq'),
    name       text NOT NULL,
    type       text NOT NULL,
	form_id    int  NOT NULL REFERENCES forms(id) ON DELETE CASCADE, 
	identifier text NOT NULL
);`

var Forms_types = `CREATE TABLE IF NOT EXISTS forms_types (
    id            int NOT NULL PRIMARY KEY DEFAULT NEXTVAL('forms_types_id_seq'),
    form_id       int NOT NULL REFERENCES forms(id) ON DELETE CASCADE,
    type_id       int NOT NULL REFERENCES event_types(id) ON DELETE CASCADE, 
	serial_number int NOT NULL
);`

var Param_values = `CREATE TABLE IF NOT EXISTS param_values (
    id        int  NOT NULL PRIMARY KEY DEFAULT NEXTVAL('param_values_id_seq'),
    person_id int  NOT NULL REFERENCES persons(id) ON DELETE CASCADE,
    event_id  int  NOT NULL REFERENCES events(id) ON DELETE CASCADE, 
    param_id  int  NOT NULL REFERENCES params(id) ON DELETE CASCADE, 
	value     text NOT NULL
);`

var Persons_events = `CREATE TABLE IF NOT EXISTS persons_events (
    id        int  NOT NULL PRIMARY KEY DEFAULT NEXTVAL('persons_events_id_seq'),
    person_id int  NOT NULL REFERENCES persons(id) ON DELETE CASCADE,
    event_id  int  NOT NULL REFERENCES events(id) ON DELETE CASCADE, 
    reg_date  date NOT NULL,
    last_date date NOT NULL
);`

//var Select1 = "set client_encoding='WIN866';"
//var Select1 = "set client_encoding='win1251';"
var Select1 = "set client_encoding='utf8';"
var Select2 = "SET TIME ZONE +10;"

var Insert0 = `INSERT INTO events (name, date_start, date_end, time, url) VALUES ('Персональные данные', '1000-01-01', '1000-01-01', '00:00', '/');`
var Insert1 = `INSERT INTO events (name, date_start, date_end, time, url) VALUES ('Математика', '2014-01-12', '2014-02-12', '13:20', '/');`
var Insert2 = `INSERT INTO events (name, date_start, date_end, time, url) VALUES ('Информатика', '2014-03-13', '2014-04-12', '15:10', '/');`
var Insert3 = `INSERT INTO events (name, date_start, date_end, time, url) VALUES ('Литература', '2014-05-14', '2014-06-12', '16:40', '/');`
var Insert4 = `INSERT INTO events (name, date_start, date_end, time, url) VALUES ('Физика', '2014-07-15', '2014-08-12', '11:50', '/');`
var Insert5 = `INSERT INTO events (name, date_start, date_end, time, url) VALUES ('Химия', '2014-09-16', '2014-10-12', '10:10', '/');`

var Insert6 = `INSERT INTO event_types (name, description, topicality) VALUES ('Олимпиада', 'Mumbo-jumbo', true);`
var Insert7 = `INSERT INTO event_types (name, description, topicality) VALUES ('Турнир', 'Mumbo-jumbo', true);`
var Insert8 = `INSERT INTO event_types (name, description, topicality) VALUES ('Соревнование', 'Mumbo-jumbo', true);`
var Insert9 = `INSERT INTO event_types (name, description, topicality) VALUES ('Состязание', 'Mumbo-jumbo', true);`
var Insert10 = `INSERT INTO event_types (name, description, topicality) VALUES ('Семинар', 'Mumbo-jumbo', true);`

var Insert11 = `INSERT INTO teams (name) VALUES ('Жулики');`
var Insert12 = `INSERT INTO teams (name) VALUES ('Умники');`
var Insert13 = `INSERT INTO teams (name) VALUES ('Клоуны');`
var Insert14 = `INSERT INTO teams (name) VALUES ('Эксперты');`
var Insert15 = `INSERT INTO teams (name) VALUES ('Новички');`

var Insert16 = `INSERT INTO persons (fname, lname, pname) VALUES ('Башмачкин', 'Акакий', 'Акакиевич');`
var Insert17 = `INSERT INTO persons (fname, lname, pname) VALUES ('Пушкин', 'Александр', 'Сергеевич');`
var Insert18 = `INSERT INTO persons (fname, lname, pname) VALUES ('Гоголь', 'Николай', 'Васильевич');`
var Insert19 = `INSERT INTO persons (fname, lname, pname) VALUES ('Толстой', 'Лев', 'Николаевич');`
var Insert20 = `INSERT INTO persons (fname, lname, pname) VALUES ('Лермонтов', 'Михаил', 'Юрьевич');
`
var Insert21 = `INSERT INTO forms (name) VALUES ('№1');`
var Insert22 = `INSERT INTO forms (name) VALUES ('№2');`
var Insert23 = `INSERT INTO forms (name) VALUES ('№3');`
var Insert24 = `INSERT INTO forms (name) VALUES ('№4');`
var Insert25 = `INSERT INTO forms (name) VALUES ('№5');`

var Insert26 = `INSERT INTO params (name, type, form_id, identifier) VALUES ('Дата рождения', '0', '1', 'date');`

var Insert27 = `INSERT INTO params (name, type, form_id, identifier) VALUES ('Регион', '1', '2', 'region');`
var Insert28 = `INSERT INTO params (name, type, form_id, identifier) VALUES ('Район', '2', '2', 'district');`
var Insert29 = `INSERT INTO params (name, type, form_id, identifier) VALUES ('Город', '3', '2', 'city');`
var Insert30 = `INSERT INTO params (name, type, form_id, identifier) VALUES ('Улица', '4', '2', 'street');`
var Insert31 = `INSERT INTO params (name, type, form_id, identifier) VALUES ('Дом', '5', '2', 'building');`

var Insert32 = `INSERT INTO params (name, type, form_id, identifier) VALUES ('Класс', '6', '3', 'class');`
var Insert33 = `INSERT INTO params (name, type, form_id, identifier) VALUES ('Наставники', '9', '3', 'teachers');`

var Insert34 = `INSERT INTO params (name, type, form_id, identifier) VALUES ('M', '7', '4', 'sex');`
var Insert35 = `INSERT INTO params (name, type, form_id, identifier) VALUES ('Ж', '7', '4', 'sex');`
var Insert36 = `INSERT INTO params (name, type, form_id, identifier) VALUES ('Использовать контактные данные для рассылки', '8', '4', 'rights');`

var Insert37 = `INSERT INTO events_types (event_id, type_id) VALUES (1, 1);`
var Insert38 = `INSERT INTO events_types (event_id, type_id) VALUES (1, 5);`
var Insert39 = `INSERT INTO events_types (event_id, type_id) VALUES (4, 1);`
var Insert40 = `INSERT INTO events_types (event_id, type_id) VALUES (4, 5);`

var Insert41 = `INSERT INTO forms_types (form_id, type_id, serial_number) VALUES (1, 1, 1);`
var Insert42 = `INSERT INTO forms_types (form_id, type_id, serial_number) VALUES (3, 1, 2);`

var Insert43 = `INSERT INTO forms_types (form_id, type_id, serial_number) VALUES (2, 5, 3);`
var Insert44 = `INSERT INTO forms_types (form_id, type_id, serial_number) VALUES (4, 5, 4);`
