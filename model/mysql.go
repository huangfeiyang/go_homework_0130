package mysql

import (
	"database/sql"
	"go_homework_0130/common/log"
	"strings"
	_ "github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
)

//不是一个连接，是数据库抽象接口。可以根据driver打开关闭数据库连接，管理连接池
var (
	DB *sql.DB
)

const (
	userName = "root"
	password = "root"
	ip       = "127.0.0.1"
	dbName   = "local"
	port     = "3306"
)

//初始化连接数据库
func InitDb() (err error) {
	dbPath := strings.Join([]string{userName, ":", password, "@tcp(", ip, ":", port, ")/", dbName, "?charset=utf8"}, "")
	DB, err = sql.Open("mysql", dbPath)
	if err != nil {
		return errors.Wrapf(err, "Database open fail: dbPath is %s", dbPath)
	}

	DB.SetConnMaxLifetime(100)
	DB.SetConnMaxIdleTime(10)

	err = DB.Ping()
	if err != nil {
		return errors.Wrapf(err, "Database connect fail: dbPath is %s", dbPath)
	}
	return nil
}

func QueryRow(str string, values []interface{}) *sql.Row {
	log.Info(str, values)
	return DB.QueryRow(str, values...)
}

func Exec(sql string, values []interface{}) (result sql.Result, execErr error) {
	log.Info(sql, values)
	if result, execErr = DB.Exec(sql, values...); execErr != nil {

		//空查询处理，只记录
		if ok := errors.Is(errors.Cause(execErr), errors.New("sql: no rows in result set")); !ok {
			log.Infof("no data: ", execErr)
			return nil, nil
		}

		return nil, errors.Wrap(execErr, "sql exec fail")
	}

	return result, nil
}