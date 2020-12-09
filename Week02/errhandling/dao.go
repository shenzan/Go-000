package main

import (
	"database/sql"
	"log"
)

type ErrorNotFound struct {
	error
	msg string
}

func (e *ErrorNotFound) IsNotFound() bool { return e.error == sql.ErrNoRows }
func (e *ErrorNotFound) ErrorMsg() string { return e.msg }

func findUserNameById(id int) (string, error) {
	var name string
	//err := db.QueryRow("select name from users where id = ?", 1).Scan(&name)
	err := sql.ErrNoRows
	if err != nil {
		if err == sql.ErrNoRows {
			err = &ErrorNotFound{err, "user not found"}
		} else {
			log.Fatal(err)
		}
	}
	return name, err
}
