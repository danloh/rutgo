// user CRUD

package db

import (
	"rutgo/util"
)

// NewUser POST
func NewUser(name, psw, email string) error {
	stmtNew, err := db.Prepare(
		"INSERT INTO user (name, pswmd, email) VALUES (?,?,?)")
	defer stmtNew.Close()
	if err != nil {
		return err
	}
	pswMd := util.Cipher(psw)
	_, err = stmtNew.Exec(name, pswMd, email)
	if err != nil {
		return err
	}
	return nil
}

// UpdateUser PUT
func UpdateUser() error {
	return nil
}
