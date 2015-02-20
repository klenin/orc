package db

import "fmt"

func Init() {
    for i, v := range Tables {
        Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s CASCADE;", v), nil)
        Exec(fmt.Sprintf("DROP SEQUENCE IF EXISTS %s_id_seq;", v), nil)
        QueryCreateTable_(Tables[i])
    }
}
