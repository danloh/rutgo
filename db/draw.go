// user CRUD

package db

import (
	"database/sql"
	"log"
	"rutgo/util"
	"time"
)

// NewDraw POST
func NewDraw(title, author, url string) error {
	stmtNew, err := db.Prepare(
		"INSERT INTO draw (title, author, url, createat, uid) VALUES (?,?,?,?,?)")
	defer stmtNew.Close()
	if err != nil {
		return err
	}
	t := time.Now()
	uid := "D" + util.GenSID()
	_, err = stmtNew.Exec(title, author, url, t, uid)
	if err != nil {
		return err
	}
	return nil
}

// UpdateDraw PUT
func UpdateDraw() error {
	return nil
}

// GetDraw GET
func GetDraw(uid string) (*Draw, error) {
	stmtGet, err := db.Prepare(
		"SELECT (title, author, url, createat) FROM draw WHERE uid = ?")
	defer stmtGet.Close()
	if err != nil {
		return nil, err
	}

	var title, author, url, createat string
	err = stmtGet.QueryRow(uid).Scan(&title, &author, &url, &createat)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	if err == sql.ErrNoRows {
		return nil, nil
	}
	res := &Draw{ID: 0, Title: title, Author: author, URL: url, CreateAt: createat, UID: uid}

	return res, nil
}

// DeleteDraw DELETE
func DeleteDraw(uid string) error {
	stmtDel, err := db.Prepare("DELETE FROM draw WHERE uid =?")
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
