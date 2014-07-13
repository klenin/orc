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

var Forms = `CREATE TABLE IF NOT EXISTS forms (
    id   int  NOT NULL PRIMARY KEY DEFAULT NEXTVAL('forms_id_seq'),
    name text NOT NULL
);`

var Params = `CREATE TABLE IF NOT EXISTS params (
    id         int  NOT NULL PRIMARY KEY DEFAULT NEXTVAL('params_id_seq'),
    name       text NOT NULL UNIQUE,
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
