// user CRUD

package db

import (
	"database/sql"
	"log"
	"rutgo/util"
)

// NewUser POST
func NewUser(name, psw, email string) error {
	stmtNew, err := db.Prepare(
		"INSERT INTO user (name, pswmd, email, uid) VALUES (?,?,?,?)")
	defer stmtNew.Close()
	if err != nil {
		return err
	}
	pswMd := util.Cipher(psw)
	uid := util.GenSID()
	_, err = stmtNew.Exec(name, pswMd, email, uid)
	if err != nil {
		return err
	}
	return nil
}

// UpdateUser PUT
func UpdateUser() error {
	return nil
}

// GetUser GET
func GetUser(uid, name string) (*User, error) {
	stmtGet, err := db.Prepare(
		"SELECT (name, email, uid) FROM user WHERE uid = ? or name = ?")
	defer stmtGet.Close()
	if err != nil {
		return nil, err
	}

	var email string
	err = stmtGet.QueryRow(name, uid).Scan(&name, &email, &uid)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	if err == sql.ErrNoRows {
		return nil, nil
	}
	res := &User{ID: 0, Pswmd: "nil", Name: name, Email: email, UID: uid}

	return res, nil
}

// DeleteUser DELETE
func DeleteUser(uid string) error {
	stmtDel, err := db.Prepare("DELETE FROM user WHERE uid =?")
	defer stmtDel.Close()
	if err != nil {
		log.Printf("%s", err)
		return err
	}
	_, err = stmtDel.Exec(uid)
	if err != nil {
		return err
	}

	return nil
}
