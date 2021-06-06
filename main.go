package main

import (
	"database/sql"
	"fmt"
	errors "github.com/pkg/errors"
	"jksjHomework/utils"
)

func main() {
	user, err := GetUserByID()
	if err != nil {
		fmt.Printf("%+v", err)
		return
	}
	if user != nil {
		fmt.Println(user)
	} else {
		fmt.Println("this user not exist")
	}

}
func GetUserByID() (*User, error) {

	user, err := FindUserByID("4")
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		} else {
			return nil, err
		}
	}
	return user, err
}

func FindUserByID(userID string) (*User, error) {
	sqlStr := "select id,userName from users where id = 100"
	user := &User{}
	err := utils.Db.QueryRow(sqlStr).Scan(&user.ID, &user.UserName)
	if err != nil {
		return nil, errors.Wrap(err, "FindUserByID err")
	}
	return user, nil
}

// User 结构体
type User struct {
	ID       int
	UserName string
}
