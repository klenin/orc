package controllers

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/orc/db"
	"github.com/orc/sessions"
	"github.com/orc/utils"
	"regexp"
	"strconv"
	"time"
)

func MatchRegexp(pattern, str string) bool {
	result, _ := regexp.MatchString(pattern, str)
	return result
}

func GetMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

func (this *Handler) HandleRegister(login, password, role, fname, lname, pname string) string {
	result := map[string]string{"result": "ok"}
	salt := time.Now().Unix()
	hash := GetMD5Hash(password + strconv.Itoa(int(salt)))

	passHasInvalidChars := false
	for i := 0; i < len(password); i++ {
		if strconv.IsPrint(rune(password[i])) == false {
			passHasInvalidChars = true
			break
		}
	}

	isExist := db.IsExists("users", "login", login)
	if isExist == true {
		result["result"] = "loginExists"
	} else if !MatchRegexp("^[a-zA-Z0-9]{2,36}$", login) {
		result["result"] = "badLogin"
	} else if !MatchRegexp("^.{6,36}$", password) && !passHasInvalidChars {
		result["result"] = "badPassword"
	} else {
		query := db.QueryInsert("persons", []string{"fname", "lname", "pname"})
		db.Query(query, []interface{}{fname, lname, pname})

		db.GetNextId("persons")
		p_id, err := strconv.Atoi(db.GetCurrId("persons"))
		utils.HandleErr("[Haldler.Index]: strconv.Atoi", err, nil)

		fmt.Println("curr :", p_id)

		query = db.QueryInsert("users", []string{"login", "pass", "salt", "role", "person_id"})
		db.Query(query, []interface{}{login, hash, salt, role, p_id - 1})
	}

	response, err := json.Marshal(result)
	utils.HandleErr("[HandleRegister] json.Marshal: ", err, nil)
	return string(response)
}

func (this *Handler) HandleLogin(login, pass string) string {
	result := map[string]interface{}{"result": "invalidCredentials"}
	isExist := db.IsExists("users", "login", login)
	if isExist {
		query := db.QuerySelect("users", "login=$1", []string{"id", "pass", "salt"})
		row := db.QueryRow(query, []interface{}{login})
		var id, hash, salt string
		row.Scan(&id, &hash, &salt)
		if hash == GetMD5Hash(pass+salt) {
			result["result"] = "ok"
			sessions.SetSession(id, login, this.Response)
			fmt.Println("HandleLogin time: ", sessions.GetValue("time", this.Request))
		}
	}
	response, err := json.Marshal(result)
	utils.HandleErr("[HandleLogin] json.Marshal: ", err, nil)
	return string(response)
}

func (this *Handler) HandleLogout() string {
	result := map[string]string{"result": "ok"}
	sessions.ClearSession(this.Response)
	fmt.Println("HandleLogout time: ", sessions.GetValue("time", this.Request))
	response, err := json.Marshal(result)
	utils.HandleErr("[HandleLogout] json.Marshal: ", err, nil)
	return string(response)
}
