package scheme

import (
	"fmt"
	"github.com/orc/db"
	"github.com/orc/mvc/controllers"
)

func Init() {
	//Drop()
	for i, _ := range db.Tables {
		controllers.GetModel(db.Tables[i]).Create()
	}
}

func Drop() {
	for _, v := range db.Tables {
		db.Query(fmt.Sprintf("DROP TABLE IF EXISTS %s CASCADE;", v), nil)
		db.Query(fmt.Sprintf("DROP SEQUENCE IF EXISTS %s_id_seq;", v), nil)
	}
}
