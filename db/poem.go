// user CRUD

package db

import (
	"database/sql"
	"log"
	"rutgo/util"
)

// NewPoem POST
func NewPoem(title, author, content string) error {
	stmtNew, err := db.Prepare(
		"INSERT INTO poem (title, author, content, uid) VALUES (?,?,?,?)")
	defer stmtNew.Close()
	if err != nil {
		return err
	}
	uid := "P" + util.GenSID()
	_, err = stmtNew.Exec(title, author, content, uid)
	if err != nil {
		return err
	}
	return nil
}

// UpdatePoem PUT
func UpdatePoem() error {
	return nil
}

// GetPoem GET
func GetPoem(uid string) (*Poem, error) {
	stmtGet, err := db.Prepare(
		"SELECT (title, author, content) FROM poem WHERE uid = ?")
	defer stmtGet.Close()
	if err != nil {
		return nil, err
	}

	var title, author, content string
	err = stmtGet.QueryRow(uid).Scan(&title, &author, &content)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	if err == sql.ErrNoRows {
		return nil, nil
	}
	res := &Poem{ID: 0, Title: title, Author: author, Content: content, UID: uid}

	return res, nil
}

// DeletePoem DELETE
func DeletePoem(uid string) error {
	stmtDel, err := db.Prepare("DELETE FROM poem WHERE uid =?")
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
