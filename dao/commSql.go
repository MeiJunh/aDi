package dao

import (
	"aDi/log"
	"aDi/util"
	"database/sql"
	"errors"
	"github.com/jmoiron/sqlx"
)

type MSqlClient struct {
	*sqlx.DB
}

// FindOneWithNull 允许为空
func (mc MSqlClient) FindOneWithNull(result interface{}, rawSql string, args ...interface{}) error {
	err := mc.Get(result, rawSql, args...)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		log.Errorf("find one fail,query:%s,param:%+v,err:%s", rawSql, args, err.Error())
		return err
	}
	return nil
}

// FindList 查找列表
func (mc MSqlClient) FindList(listPtr interface{}, rawSql string, args ...interface{}) error {
	err := mc.Select(listPtr, rawSql, args...)
	if err != nil {
		log.Errorf("find list fail,query:%s,param:%+v,err:%s", rawSql, args, err.Error())
		return err
	}
	return err
}

// Transaction 需要使用事物进行操作，添加监控上报
func Transaction(client *sqlx.DB, task func(tx *sqlx.Tx) error) (err error) {
	callerName := util.GetCallerName()
	// 开启事物
	mTX, err := client.Beginx()
	if err != nil {
		log.Errorf("%s begin fail,err:%s", callerName, err.Error())
		return err
	}
	defer func() {
		// catch 住panic,免得因为panic导致事务没有提交和回滚导致锁住
		panicErr := recover()
		if err != nil || panicErr != nil {
			// 如果出错则进行回滚
			errT := mTX.Rollback()
			if errT != nil {
				log.Errorf("rollback fail,err:%s", errT)
			}
		}
	}()

	err = task(mTX)
	if err != nil {
		log.Errorf("task fail,err:%s", err.Error())
		return err
	}

	// 提交错误直接返回错误
	err = mTX.Commit()
	if err != nil {
		log.Errorf("%s commit fail,err:%s", callerName, err.Error())
		return err
	}
	return err
}
