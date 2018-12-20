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
)

type Rows map[string]interface{}

type context struct {
	sql string
	db *sql.DB
}

type DataAccess interface {
	Query() (Rows, error)
	QueryCount() (int8, error)
	Update() (int8, error)
	Delete() (int8, error)
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
		username: "",
		password: "",
	},
	"test": {},
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
		update(updateValue)
	} else if selectValue != "" {
		query(selectValue)
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
	From(sql).Query()
	return nil, nil
}

// SQL更新或者删除
func update(sql string) (int8, error) {
	From(sql).Update()
	return 0, nil
}

var dbOnce sync.Once
var pool *sql.DB

func initDataSourceIfNecessary() {
	dbOnce.Do(func() {
		pool_, err := sql.Open("mysql", "user:password@tcp(localhost:5555)/dbname?charset=utf8mb4,utf8&autocommit=false&parseTime=true")
		if err != nil {
			detail := fmt.Sprintf("数据库连接失败: [%s].", err)
			Logger.Warn(detail)
			panic(detail)
		}

		pool = pool_
	})
}

func From(sql string) DataAccess {
	return &context{
		sql: sql,
		db:  pool,
	}
}

func (db *context) Query() (Rows, error) {
	initDataSourceIfNecessary()
	result, err := pool.Exec(db.sql)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	return nil, nil
}

func (db *context) QueryCount() (int8, error) {
	initDataSourceIfNecessary()
	return 0, nil
}

func (db *context) Update() (int8, error) {
	initDataSourceIfNecessary()
	return 0, nil
}

func (db *context) Delete() (int8, error) {
	initDataSourceIfNecessary()
	return 0, nil
}



