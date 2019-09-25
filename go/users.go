package main

import (
	"database/sql"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

func getUser(r *http.Request) (user User, errCode int, errMsg string) {
	userID, errCode, errMsg := getUserID(r)
	if errMsg != "" {
		return
	}

	err := dbx.Get(&user, "SELECT * FROM `users` WHERE `id` = ?", userID)
	if err == sql.ErrNoRows {
		return user, http.StatusNotFound, "user not found"
	}
	return user, http.StatusOK, ""
}

func getUserSimple(r *http.Request) (user UserSimple, errCode int, errMsg string) {
	userID, errCode, errMsg := getUserID(r)
	if errMsg != "" {
		return
	}

	err := dbx.Get(&user, "SELECT id, account_name, num_sell_items FROM `users` WHERE `id` = ?", userID)
	if err == sql.ErrNoRows {
		return user, http.StatusNotFound, "user not found"
	}
	return user, http.StatusOK, ""
}

func getUserSimpleByID(q sqlx.Queryer, userID int64) (userSimple UserSimple, err error) {
	user := UserSimple{}
	err = sqlx.Get(q, &user, "SELECT id,account_name,num_sell_items FROM `users` WHERE `id` = ?", userID)
	return user, err
}

func getUserSimpleByIDs(q sqlx.Queryer, userIDs []int64) (userSimples map[int64]UserSimple, err error) {
	userSimples = make(map[int64]UserSimple)
	inQuery, inArgs, err := sqlx.In("SELECT id, account_name, num_sell_items FROM `users` WHERE `id` IN (?)", userIDs)
	if err != nil {
		return nil, err
	}
	rows, err := q.Queryx(inQuery, inArgs...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var user UserSimple
		err = rows.StructScan(&user)
		if err != nil {
			return nil, err
		}
		userSimples[user.ID] = user
	}
	return userSimples, nil
}
