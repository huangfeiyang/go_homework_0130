package mysql

import "testing"

func TestMysql(t *testing.T) {
	err := InitDb()
	if err != nil {
		t.Log(err)
	}
	var sql = "SELECT name FROM user_info WHERE name = '乐乐';"
	values := make([]interface{}, 0)
	row := QueryRow(sql, values)
	var name string
	if rowErr := row.Scan(&name); rowErr != nil || name != "乐乐" {
		t.Fatal("fail")
	}
}