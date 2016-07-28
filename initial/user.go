package initial

import (
	"errors"
	"github.com/klenin/orc/mailer"
	"github.com/klenin/orc/db"
	"github.com/klenin/orc/mvc/controllers"
	"github.com/klenin/orc/utils"
	"math/rand"
	"strconv"
	"time"
	"log"
)

const USER_COUNT = 20

func loadAdmin() {
	base := new(controllers.BaseController)
	date := time.Now().Format("2006-01-02T15:04:05Z00:00")

	result, regId := base.RegistrationController().Register("admin", "password", mailer.Admin_.Email, "admin")
	if result != "ok" {
		utils.HandleErr("[LoadAdmin]: "+result, nil, nil)

		return
	}

	query := `INSERT INTO param_values (param_id, value, date, user_id, reg_id)
        VALUES (4, $1, $2, NULL, $3);`
	db.Exec(query, []interface{}{mailer.Admin_.Email, date, regId})

	for k := 5; k < 8; k++ {
		query := `INSERT INTO param_values (param_id, value, date, user_id, reg_id)
            VALUES (`+strconv.Itoa(k)+`, 'admin', $1, NULL, $2);`
		db.Exec(query, []interface{}{date, regId})
	}

	query = `SELECT users.token FROM registrations
        INNER JOIN events ON registrations.event_id = events.id
        INNER JOIN faces ON faces.id = registrations.face_id
        INNER JOIN users ON users.id = faces.user_id
        WHERE events.id = $1 AND registrations.id = $2;`
	res := db.Query(query, []interface{}{1, regId})

	if len(res) == 0 {
		utils.HandleErr("[LoadAdmin]: ", errors.New("Data are not faund."), nil)

		return
	}

	token := res[0].(map[string]interface{})["token"].(string)
	base.RegistrationController().ConfirmUser(token)
}

func loadUsers() {
	base := new(controllers.BaseController)
	date := time.Now().Format("2006-01-02T15:04:05Z00:00")

	type FullNames struct {
		firstNames, lastNames, patronymics []string
	}

	male := FullNames{
		firstNames: readStringsFromFile(getResourcePath("first-name-male.txt")),
		lastNames: readStringsFromFile(getResourcePath("last-name-male.txt")),
		patronymics: readStringsFromFile(getResourcePath("patronymic-male.txt")),
	}
	female := FullNames{
		firstNames: readStringsFromFile(getResourcePath("first-name-female.txt")),
		lastNames: readStringsFromFile(getResourcePath("last-name-female.txt")),
		patronymics: readStringsFromFile(getResourcePath("patronymic-female.txt")),
	}

	for i := 0; i < USER_COUNT; i++ {
		userName := "user" + strconv.Itoa(i + 1)
		userEmail := userName + "@example.com"

		result, regId := base.RegistrationController().Register(userName, "password", userEmail, "user")
		if result != "ok" {
			log.Fatalln("[loadUsers]:", result)
		}

		query := `INSERT INTO param_values (param_id, value, date, reg_id)
            VALUES ($1, $2, $3, $4);`

		db.Exec(query, []interface{}{4, userEmail, date, regId})
		var fullNames FullNames
		if rand.Int() % 2 == 1 {
			fullNames = male
		} else {
			fullNames = female
		}
		db.Exec(query, []interface{}{6, fullNames.firstNames[rand.Intn(len(fullNames.firstNames))], date, regId})
		db.Exec(query, []interface{}{5, fullNames.lastNames[rand.Intn(len(fullNames.lastNames))], date, regId})
		db.Exec(query, []interface{}{7, fullNames.patronymics[rand.Intn(len(fullNames.patronymics))], date, regId})

		query = `SELECT users.token FROM registrations
            INNER JOIN events ON registrations.event_id = events.id
            INNER JOIN faces ON faces.id = registrations.face_id
            INNER JOIN users ON users.id = faces.user_id
            WHERE events.id = $1 AND registrations.id = $2;`
		res := db.Query(query, []interface{}{1, regId})

		if len(res) == 0 {
			log.Fatalln("[loadUsers]:", "Data are not found")
		}

		token := res[0].(map[string]interface{})["token"].(string)
		base.RegistrationController().ConfirmUser(token)
	}
}
