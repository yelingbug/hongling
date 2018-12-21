package utility

import (
	"gopkg.in/urfave/cli.v2"
	"fmt"
	"database/sql"
	"errors"
	"sync"
	"time"
	"reflect"
	_ "github.com/go-sql-driver/mysql"
	"strings"
	"regexp"
)

type Rows []map[string]interface{}

type context struct {
	sql string
	params []interface{}
	db *sql.DB
}

type DataAccess interface {
	Query() (Rows, error)
	QueryCount() (int64, error)
	Update() (int64, error)
	Delete() (int64, error)
}

var usageValue bool
var selectValue, updateValue, envValue string

type DataSourceProperties struct {
	username string
	password string
	host string
	port int
	schema string

	connMaxLifetime time.Duration
	maxIdleConns int
	maxOpenConns int
}

var environments = map[string]DataSourceProperties {
	"dev" : {
		username: "C65E183D16B1A210",
		password: "646D541782E6D78C51923A35A3D19D3E",
	},
	"test": {
		username: "",
		password: "C3940641962CA5BDFF24F7D079E9D5D31385ACCB7B67EB68",
	},
	"prod": {},
}

var DrdsCommand = &cli.Command{
	Name: "drds",
	Category:"DRDS数据库操作",
	Aliases: []string{"mysql", "db"},
	Usage:"hl [global options] drds/mysql/db [command options] [arguments...]",
	Action: drds,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:"environment",
			Aliases: []string{"env", "e"},
			Usage: "选择环境`dev/test/prod`, 缺省为dev.",
			Value: "dev",
			Destination: &envValue,
		},
		&cli.BoolFlag{
			Name:"usage",
			Aliases:[]string{"u"},
			Usage: "查看drds命令的用法.",
			Destination: &usageValue,
		},
		&cli.StringFlag{
			Name:"select",
			Aliases:[]string{"query", "q", "s"},
			Usage: "`SQL`查询返回结果以JSON格式返回.",
			Destination: &selectValue,
		},
		&cli.StringFlag{
			Name:"update",
			Aliases:[]string{"set", "modify"},
			Usage: "`SQL`更新返回影响结果行数.",
			Destination: &updateValue,
		},
	},
}

// 通过命令行传入SQL执行数据库查询和更新.
func drds(c *cli.Context) error {
	if _, ok := environments[envValue]; !ok {
		return errors.New(fmt.Sprintf("环境只能是%s其中之一,默认为dev.", reflect.ValueOf(environments).MapKeys()))
	}

	if updateValue != "" {
		if rows, err := update(updateValue); err != nil {
			Logger.Info(fmt.Sprintf("更新失败:[%s].", err))
		} else {
			Logger.Info(fmt.Sprintf("更新结果影响%d行.", rows))
		}
	} else if selectValue != "" {
		if rows, err := query(selectValue); err != nil {
			Logger.Info(fmt.Sprintf("查询失败:[%s].", err))
		} else {
			Logger.Info(fmt.Sprintf("查询结果:[%s]", rows))
		}

	} else if usageValue {
		fmt.Println(`用法:
  查询类:
    hl drds/mysql/db -select/-query/-q/-s [SQL语句]
      例如:hl drds/mysql/db -select "select * from uc_user where user_id = 12345"
  更新/删除类:
    hl drds/mysql/db -update/-set/-modify [SQL更新语句,支持update/delete]
      例如:hl drds/mysql/db -update "update uc_user set login_username='abc' where user_id = 12345"`)
	}
	return nil
}

// SQL查询
func query(sql string) (Rows, error) {
	return From(sql).Query()
}

// SQL更新或者删除
func update(sql string) (int64, error) {
	sqlToLower := strings.ToLower(sql)
	matchedForUpdate, errForUpdate := regexp.MatchString("\\s*update\\s*\\w\\s*set\\s*.*", sqlToLower)
	matchedForDelete, errForDelete := regexp.MatchString("\\s*delete\\s*from\\s*.*", sqlToLower)
	if errForUpdate != nil || errForDelete != nil {
		return 0, errors.New(fmt.Sprintf("正则表达式匹配错误:[%s/%s].", errForUpdate, errForDelete))
	}

	if !matchedForDelete && !matchedForUpdate {
		return 0, errors.New("SQL语句非更新或者删除语句.")
	}

	return From(sql).Update()
}

var dbOnce sync.Once
var pool *sql.DB

func initDataSourceIfNecessary() {
	dbOnce.Do(func() {
		pool_, err := sql.Open("mysql", environments[envValue].username + ":" + environments[envValue].password + "@tcp(" + environments[envValue].host + ":" + string(environments[envValue].port) + ")/" + environments[envValue].schema + "?charset=utf8mb4,utf8&autocommit=false&parseTime=true")
		if err != nil {
			detail := fmt.Sprintf("数据库连接失败: [%s].", err)
			Logger.Warn(detail)
			panic(detail)
		}

		pool = pool_

		pool.SetConnMaxLifetime(environments[envValue].connMaxLifetime)
		pool.SetMaxIdleConns(environments[envValue].maxIdleConns)
		pool.SetMaxOpenConns(environments[envValue].maxOpenConns)
	})
}

func From(sql string) DataAccess {
	return &context{
		sql: sql,
		db:  pool,
	}
}

func FromByParams(sql string, params ...interface{}) DataAccess {
	return &context{
		sql: sql,
		params: params,
		db:  pool,
	}
}

func (db *context) Query() (Rows, error) {
	initDataSourceIfNecessary()
	rows, err := pool.Query(db.sql)
	if err != nil {
		return nil, err
	}

	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	values := make([]interface{}, len(columns))
	valuesPtr := make([]interface{}, len(columns))
	for i := 0; i < len(columns); i++ {
		valuesPtr[ i ] = &values[i]
	}

	results := make([]map[string]interface{}, 0)
	for rows.Next() {
		rows.Scan(valuesPtr...)

		row := make(map[string]interface{})
		for i, column := range columns {
			v, ok := values[ i ].([]byte)
			if ok {
				row[ column ] = string(v)
			} else {
				row[ column ] = v
			}

		}
		results = append(results, row)

	}
	return results, nil
}

func (db *context) QueryCount() (int64, error) {
	initDataSourceIfNecessary()
	row := pool.QueryRow(db.sql)

	var value int64
	row.Scan(&value)
	return value, nil
}

func (db *context) Update() (int64, error) {
	initDataSourceIfNecessary()
	tx, err := pool.Begin()
	if err != nil {
		return 0, err
	}

	result, err := pool.Exec(db.sql)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	tx.Commit()

	return result.RowsAffected()
}

func (db *context) Delete() (int64, error) {
	initDataSourceIfNecessary()

	tx, err := pool.Begin()
	if err != nil {
		return 0, err
	}

	result, err := pool.Exec(db.sql)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	tx.Commit()

	return result.RowsAffected()
}





