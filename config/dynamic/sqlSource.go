package dynamic

import (
	"aDi/log"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"time"
)

type sqlSource struct {
	sqlQuery  string   // db 查询语句
	dbHandler *sqlx.DB // db handler
}

// NewSQLConfSourceByURL 通过db url提供mysql source来源
func NewSQLConfSourceByURL(dbURL string, opt ...Option) (*sqlSource, error) {
	// 初始化db
	db, err := sqlx.Connect("mysql", dbURL)
	if err != nil {
		log.Errorf("init config source db fail,db url:%s,err:%s", dbURL, err.Error())
		return nil, err
	}
	log.Debugf("init sql config source success,db url:%s", dbURL)

	// 设置最大连接数等信息
	db.SetMaxOpenConns(8)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(25 * time.Minute)
	go func() {
		for {
			log.Debugf("config refresh db status:%+v", db.Stats())
			time.Sleep(time.Minute * 10)
		}
	}()

	return NewSQLConfSourceByDBHandler(db, opt...), nil
}

// NewSQLConfSourceByDBHandler 直接通过db handler进行初始化
func NewSQLConfSourceByDBHandler(mDB *sqlx.DB, opt ...Option) *sqlSource {
	source := &sqlSource{
		sqlQuery:  DefSqlQuery,
		dbHandler: mDB,
	}

	for _, o := range opt {
		o(source)
	}
	return source
}

// Get sql get func
func (s *sqlSource) Get(serviceName, key string) (value string, err error) {
	value = ""
	rows, err := s.dbHandler.Query(s.sqlQuery, serviceName, key)
	if err != nil {
		log.Errorf("query fail,service:%s,key:%s,err:%s", serviceName, key, err.Error())
		return value, err
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&value)
		if err != nil {
			log.Errorf("scan fail,service:%s,key:%s,err:%s", serviceName, key, err.Error())
			return value, err
		}
		return value, nil
	}

	return
}

const (
	DefSqlQuery = "select config_value from sys_config where service_name = ? and config_key = ?"
)

type Option func(options *sqlSource)

// SqlQuery 自定义sql语句
func SqlQuery(sqlQuery string) Option {
	return func(options *sqlSource) {
		options.sqlQuery = sqlQuery
	}
}
