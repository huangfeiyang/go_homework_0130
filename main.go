package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
	"strings"
)

const (
	userName = "root"
	password = "root"
	ip = "127.0.0.1"
	dbName = ""
)