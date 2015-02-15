package scheme

import (
    "fmt"
    "github.com/orc/db"
)

func Init() {
    for i, v := range db.Tables {
        db.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s CASCADE;", v), nil)
        db.Exec(fmt.Sprintf("DROP SEQUENCE IF EXISTS %s_id_seq;", v), nil)
        db.QueryCreateTable_(db.Tables[i])
    }
}
