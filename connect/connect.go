package connect

import (
	"asd-swap/common/mysql"
	"asd-swap/swap-api/config"
	"fmt"
	"gorm.io/gorm"
	"sync"
)

type DefaultPool struct {
	Mysql *gorm.DB
	once  sync.Once
}

var _default = &DefaultPool{}

func Default() *DefaultPool {
	_default.once.Do(initialize)
	return _default
}

func initialize() {
	var err error
	_default.Mysql, err = mysql.NewMysqlClient(*config.Config().Mysql)
	if err != nil {
		fmt.Printf("ERROR initialize connect pool mysql error: %v\n", err)
	}
}
