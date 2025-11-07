// Package database provides support for multiple database instances including MySQL and Redis.
package database

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	// 使用map存储多个MySQL实例
	mysqlInstances = make(map[string]*gorm.DB)
	// 保证线程安全的互斥锁
	mu sync.RWMutex
)

// MySQLOptions defines options for mysql database.
type MySQLOptions struct {
	Instance              string
	Addr                  string
	Username              string
	Password              string
	Database              string
	MaxIdleConnections    int
	MaxOpenConnections    int
	MaxConnectionLifeTime time.Duration
	// SlogLogger is the slog logger for MySQL operations
	Logger *slog.Logger
}

// DSN return DSN from MySQLOptions.
func (o *MySQLOptions) DSN() string {
	return fmt.Sprintf(`%s:%s@tcp(%s)/%s?charset=utf8&parseTime=%t&loc=%s`,
		o.Username,
		o.Password,
		o.Addr,
		o.Database,
		true,
		"Local")
}

// InitMySQL 初始化单个MySQL实例
func InitMySQL(opts *MySQLOptions) (*gorm.DB, error) {
	gormConfig := &gorm.Config{
		// PrepareStmt executes the given query in cached statement.
		// This can improve performance.
		PrepareStmt: true,
	}
	if opts.Logger != nil {
		gormConfig.Logger = &mysqlLogger{
			logger: opts.Logger.With(slog.String("component", "mysql")),
		}
	}

	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       opts.DSN(),
		DefaultStringSize:         256,   // string 类型字段的默认长度
		DisableDatetimePrecision:  true,  // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    true,  // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true,  // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false, // 根据当前 MySQL 版本自动配置
	}), gormConfig)
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	// SetMaxOpenConns sets the maximum number of open connections to the database.
	sqlDB.SetMaxOpenConns(opts.MaxOpenConnections)

	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	sqlDB.SetConnMaxLifetime(opts.MaxConnectionLifeTime)

	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sqlDB.SetMaxIdleConns(opts.MaxIdleConnections)

	mu.Lock()
	mysqlInstances[opts.Instance] = db
	mu.Unlock()

	return db, nil
}

// InitMySQLs 批量初始化MySQL实例
func InitMySQLs(opts []*MySQLOptions) error {
	for _, opt := range opts {
		if _, err := InitMySQL(opt); err != nil {
			return fmt.Errorf("failed to initialize MySQL instance '%s': %w", opt.Instance, err)
		}
	}
	return nil
}

// GetMySQL 获取指定名称的MySQL实例，如果不传入名称或传入空字符串，则返回默认实例（第一个实例）
func GetMySQL(instances ...string) *gorm.DB {
	mu.RLock()
	defer mu.RUnlock()

	instance := "default"
	if len(instances) > 0 && instances[0] != "" {
		instance = instances[0]
	}

	if db, exists := mysqlInstances[instance]; exists {
		return db
	}

	// 如果找不到指定实例，返回第一个实例作为默认实例
	for _, db := range mysqlInstances {
		return db
	}

	return nil
}

// GetMySQLInstances 获取所有MySQL实例名称
func GetMySQLInstances() []string {
	mu.RLock()
	defer mu.RUnlock()

	instances := make([]string, 0, len(mysqlInstances))
	for name := range mysqlInstances {
		instances = append(instances, name)
	}
	return instances
}

// CloseMySQL 关闭指定的MySQL实例连接
func CloseMySQL(_ context.Context, instances ...string) error {
	mu.Lock()
	defer mu.Unlock()

	// 如果没有指定实例，则关闭所有实例
	if len(instances) == 0 {
		for name, db := range mysqlInstances {
			if sqlDB, err := db.DB(); err == nil {
				if err = sqlDB.Close(); err != nil {
					return fmt.Errorf("failed to close MySQL instance '%s': %w", name, err)
				}
			}
			delete(mysqlInstances, name)
		}
		return nil
	}

	// 关闭指定实例
	for _, instance := range instances {
		if db, exists := mysqlInstances[instance]; exists {
			if sqlDB, err := db.DB(); err == nil {
				if err = sqlDB.Close(); err != nil {
					return fmt.Errorf("failed to close MySQL instance '%s': %w", instance, err)
				}
			}
			delete(mysqlInstances, instance)
		}
	}
	return nil
}
