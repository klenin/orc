package initial

import (
	"github.com/klenin/orc/db"
	"github.com/klenin/orc/mvc/models"
	"math/rand"
	"time"
	"fmt"
)

func Init(resetDB, loadTestData bool) {
	if resetDB {
		clearDatabase()
		loadAdmin()
		loadParamTypes()
		loadForms()
		createRegistrationEvent()
	}

	if loadTestData {
		rand.Seed(time.Now().UnixNano())

		loadUsers()
		loadEvents()
		loadEventTypes()
	}
}

var base = new(models.ModelManager)

func clearDatabase() {
	for k, v := range db.Tables {
		db.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s CASCADE;", v), nil)
		db.Exec(fmt.Sprintf("DROP SEQUENCE IF EXISTS %s_id_seq;", v), nil)
		db.QueryCreateTable(base.GetModel(db.Tables[k]))
	}
}
