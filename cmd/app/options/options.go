/*
Copyright (c) 2022 The DnsJia Authors.
WebSite:  https://github.com/dnsjia/fuxi
Email:    OpenSource@dnsjia.com

MIT License

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.

*/

package options

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"

	"github.com/dnsjia/fuxi/cmd/app/config"
	"github.com/dnsjia/fuxi/pkg/db"
)

const (
	defaultConfigFile = "etc/fuxi.yaml"
)

type Options struct {
	DB              *gorm.DB
	Redis           *redis.Client
	LOG             *zap.Logger
	GinEngine       *gin.Engine
	ComponentConfig *config.Config
	Factory         db.ShareDaoFactory
	ConfigFile      string
}

func NewOptions() (*Options, error) {
	return &Options{
		ConfigFile: defaultConfigFile,
	}, nil
}

func (o *Options) Viper(path ...string) *viper.Viper {
	var configFile string
	if len(path) == 0 {
		flag.StringVar(&configFile, "c", "", "choose configFile file.")
		flag.Parse()
		if configFile == "" { // 优先级: 命令行 > 环境变量 > 默认值
			if configEnv := os.Getenv("FUXI_CONFIG"); configEnv == "" {
				configFile = "etc/fuxi.yaml"
				fmt.Printf("您正在使用config的默认值, 配置路径为%v\n", "etc/configFile.yaml")
			} else {
				configFile = configEnv
				fmt.Printf("您正在使用FUXI_CONFIG环境变量, 配置路径为%v\n", configFile)
			}
		} else {
			fmt.Printf("您正在使用命令行的-c参数传递的值, 配置路径为%v\n", configFile)
		}
	} else {
		configFile = path[0]
		fmt.Printf("您正在使用func Viper()传递的值,config的路径为%v\n", configFile)
	}

	v := viper.New()
	v.SetConfigFile(configFile)
	err := v.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error configFile file: %s \n", err))
	}
	v.WatchConfig()

	v.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("configFile file changed:", e.Name)
		if err := v.Unmarshal(&o.ComponentConfig); err != nil {
			fmt.Println(err)
		}
	})

	if err := v.Unmarshal(&o.ComponentConfig); err != nil {
		fmt.Println(err)
	}
	return v
}

func (o *Options) BindFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&o.ConfigFile, "config", "", "Please specify the configuration file path of fuxi")
}

func (o *Options) Database() error {
	m := o.ComponentConfig.Mysql
	dsn := m.Username + ":" + m.Password + "@tcp(" + m.Host + ")/" + m.Dbname + "?" + "charset=utf8mb4&parseTime=True&loc=Local"
	mysqlConfig := mysql.Config{
		DSN:                       dsn,   // DSN data source name
		DefaultStringSize:         191,   // string 类型字段的默认长度
		DisableDatetimePrecision:  true,  // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    true,  // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true,  // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false, // 根据版本自动配置
	}
	var err error
	if o.DB, err = gorm.Open(mysql.New(mysqlConfig), &gorm.Config{
		Logger: gormlogger.Default.LogMode(gormlogger.Info),
	}); err != nil {
		o.LOG.Error("无法连接到数据库")
		return err
	}

	sqlDB, err := o.DB.DB()
	if err != nil {
		return err
	}

	if err := sqlDB.Ping(); err != nil {
		return err
	}

	o.Factory = db.NewDaoFactory(o.DB)
	// 初始化表
	db.InitMysqlTables(o.DB)
	return nil
}

func (o *Options) RedisCache() *redis.Client {
	o.Redis = redis.NewClient(&redis.Options{
		Addr:         o.ComponentConfig.Redis.Addr,
		Password:     o.ComponentConfig.Redis.Password,
		PoolSize:     o.ComponentConfig.Redis.PoolSize,
		DialTimeout:  o.ComponentConfig.Redis.DialTimeout,
		ReadTimeout:  o.ComponentConfig.Redis.ReadTimeout,
		WriteTimeout: o.ComponentConfig.Redis.WriteTimeout,
	})
	return o.Redis
}

func (o Options) Set(key string, value interface{}, expiration time.Duration) error {
	return o.Redis.Set(context.TODO(), key, value, expiration).Err()
}

func (o Options) Get(key string) (string, error) {
	return o.Redis.Get(context.TODO(), key).Result()
}

func (o Options) Del(key string) error {
	return o.Redis.Del(context.TODO(), key).Err()
}

func (o Options) SetNX(key string, value interface{}, expiration time.Duration) error {
	return o.Redis.SetNX(context.TODO(), key, value, expiration).Err()
}

func (o Options) Expire(key string, expiration time.Duration) error {
	return o.Redis.Expire(context.TODO(), key, expiration).Err()
}

func (o Options) Exists(key string) error {
	return o.Redis.Exists(context.TODO(), key).Err()
}

func (o Options) HSet(key string, values ...interface{}) error {
	return o.Redis.HSet(context.TODO(), key, values...).Err()
}

func (o Options) HGet(key, field string) (string, error) {
	return o.Redis.HGet(context.TODO(), key, field).Result()
}

func (o Options) HDel(key, field string) error {
	return o.Redis.HDel(context.TODO(), key, field).Err()
}

func (o Options) MSet(values ...interface{}) error {
	return o.Redis.MSet(context.TODO(), values...).Err()
}

func (o Options) MGet(keys ...string) ([]interface{}, error) {
	return o.Redis.MGet(context.TODO(), keys...).Result()
}

func (o *Options) Complete() error {
	o.GinEngine = gin.Default()
	o.Viper()

	if err := o.register(); err != nil {
		return err
	}

	return nil
}

func (o *Options) register() error {
	if err := o.Database(); err != nil {
		return err
	}

	o.RedisCache()
	return nil
}

func (o *Options) Run(stopCh <-chan struct{}) {
	_ = o.GinEngine.Run(fmt.Sprintf(":%d", o.ComponentConfig.Http.Listen))
}
