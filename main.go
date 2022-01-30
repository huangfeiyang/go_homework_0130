package main

import (
	"database/sql"
	"fmt"
	_"github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
	"strings"
)

const (
	userName = "root"
	password = "root"
	ip = "127.0.0.1"
	dbName = "local"
	port = "3306"
)

var DB *sql.DB

type User_info struct {
	Name string
	Age int
	Test string
	Id int
}

type HANDTYPE int 

const (
	INIT HANDTYPE = iota
	INSERT
	SELECTBYNAME
)

func InitDB() (dbErr error) {
	dbPath := strings.Join([]string{userName, ":", password, "@tcp(",ip,":",port,")/", dbName, "?charset=utf8"}, "")

	DB, dbErr =sql.Open("mysql", dbPath)

	if dbErr != nil {
		return errors.Wrapf(dbErr, "main:InitDB open database fail:", dbPath)
	}

	DB.SetConnMaxLifetime(100)
	DB.SetConnMaxIdleTime(10)

	if dbErr = DB.Ping(); dbErr != nil {
		fmt.Println("sss")
		return errors.Wrapf(dbErr, "main: InitDB connect databse fail: %s", dbPath)
	}

	return nil
}

func QueryUserByName(name string) (interface{}, error) {
	var user User_info
	sql := "select * from user_info where name = ? "
	err := DB.QueryRow(sql, name).Scan(&user.Id, &user.Name, &user.Age, &user.Test)

	if err != nil {
		return nil, errors.Wrapf(err, "main:QueryUserByName Scan falied:%s,%s",sql, name)
	}
	return &user,nil
}

func InsertUser(user User_info) error {
	tx, err := DB.Begin()

	if err != nil {
		return errors.Wrap(err, "main:InsertUser db brgin falied")
	}

	stmt, err := tx.Prepare("insert into user_info(`name`, `age`, `test`) values (?, ?, ?)")
	if err != nil {
		return errors.Wrap(err, "mian:InsertUser Sql Prepare failed")
	}

	if _, errExec := stmt.Exec(user.Name, user.Age, user.Test); errExec != nil {
		return errors.Wrap(err, "main:InsertUser Exec failed")
	}

	tx.Commit()
	return nil
}

func HandleDbFunc(dType HANDTYPE, inInfo User_info) (outInfo interface{}, dbErr error) {
	switch dType {
	case INIT:
		dbErr = InitDB()
	case INSERT:
		dbErr = InsertUser(inInfo)
	case SELECTBYNAME:
		outInfo, dbErr = QueryUserByName(inInfo.Name)
	default:
		dbErr = errors.Errorf("main:HandleDbFunc unknow db type:%d", dType)
	}
	return outInfo, dbErr
}

func HandleDataFunc(inName, queryName string) (interface{}, error) {
	var outUser interface{}
	var err error = nil
	for i := HANDTYPE(0); i <= SELECTBYNAME; i++ {
		if i == SELECTBYNAME {
			outUser, err = HandleDbFunc(i, User_info{queryName, 22, "", 11})
		} else {
			_, err = HandleDbFunc(i, User_info{inName, 22, "", 11})
		}

		if err != nil {
			return nil, errors.WithMessagef(err, "main:HandleDataFunc handletype:%d", i)
		}
	}
	return outUser, nil
}

func PrintErr(err error) {
	if err != nil {
		fmt.Printf("original error: %+v\v", errors.Cause(err))
		fmt.Printf("stack trace:\n%+v\n", err)
		if errors.Cause(err) == sql.ErrNoRows {
			fmt.Println("sql res no rows!")
		}
	}
}

func main() {
	info, err := HandleDataFunc("walk", "walk")
	if err != nil {
		PrintErr(err)
		return
	}

	if user, ok := info.(*User_info); ok {
		fmt.Printf("name:%s, age:%d, id:%d", user.Name, user.Age, user.Id)
	}
}