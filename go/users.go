package main

import (
	"database/sql"
	"log"
	"net/http"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var (
	mUserCache sync.Mutex
	userCache  = make(map[int64]*User)
)

func initUserCache() {
	mUserCache.Lock()
	defer mUserCache.Unlock()

	userCache = make(map[int64]*User)

	users := []User{}
	err := dbx.Select(&users, "SELECT * FROM `users`")
	if err != nil {
		log.Panic(err)
	}

	for _, u := range users {
		v := u
		userCache[u.ID] = &v
	}
}

func updateUserBump(userID int64, now time.Time) {
	mUserCache.Lock()
	defer mUserCache.Unlock()
	user := userCache[userID]

	if user != nil {
		user.LastBump = now
	}
}

func updateUserSell(userID int64, now time.Time, items int) {
	mUserCache.Lock()
	defer mUserCache.Unlock()
	user := userCache[userID]

	if user != nil {
		user.LastBump = now
		if user.NumSellItems+1 != items {
			log.Printf("userID=%v num=%v next=%v", userID, user.NumSellItems, items)
		}
		user.NumSellItems = items
	}
}

func getUser(r *http.Request) (user User, errCode int, errMsg string) {
	userID, errCode, errMsg := getUserID(r)
	if errMsg != "" {
		return
	}

	mUserCache.Lock()
	if u := userCache[userID]; u != nil {
		u := *u
		mUserCache.Unlock()
		return u, http.StatusOK, ""
	}
	mUserCache.Unlock()

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

	mUserCache.Lock()
	if user := userCache[userID]; user != nil {
		us := UserSimple{
			ID:           user.ID,
			AccountName:  user.AccountName,
			NumSellItems: user.NumSellItems,
		}
		mUserCache.Unlock()
		return us, http.StatusOK, ""
	}
	mUserCache.Unlock()

	err := dbx.Get(&user, "SELECT /* simple */ id, account_name, num_sell_items FROM `users` WHERE `id` = ?", userID)
	if err == sql.ErrNoRows {
		return user, http.StatusNotFound, "user not found"
	}
	return user, http.StatusOK, ""
}

func getUserSimpleByID(q sqlx.Queryer, userID int64) (userSimple UserSimple, err error) {
	mUserCache.Lock()
	defer mUserCache.Unlock()

	if user := userCache[userID]; user != nil {
		us := UserSimple{
			ID:           user.ID,
			AccountName:  user.AccountName,
			NumSellItems: user.NumSellItems,
		}
		return us, nil
	}

	user := User{}
	err = sqlx.Get(q, &user, "SELECT /* simple by id */ * FROM `users` WHERE id = ?", userID)
	if err != nil {
		return UserSimple{}, err
	}
	userCache[user.ID] = &user

	userSimple.ID = user.ID
	userSimple.AccountName = user.AccountName
	userSimple.NumSellItems = user.NumSellItems
	return
}
