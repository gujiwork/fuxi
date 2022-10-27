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

package config

import "time"

type Config struct {
	Http  HttpOptions  `mapstructure:"http"  json:"http" yaml:"http"`
	Mysql MysqlOptions `mapstructure:"mysql"  json:"mysql" yaml:"mysql"`
	Redis RedisOptions `mapstructure:"redis"  json:"redis" yaml:"redis"`
}

type HttpOptions struct {
	Mode   string `mapstructure:"mode" json:"mode" yaml:"mode"`
	Listen int    `mapstructure:"listen" json:"listen" yaml:"listen"`
}

type MysqlOptions struct {
	Host     string `mapstructure:"host" json:"host" yaml:"host"`
	Config   string `mapstructure:"config" json:"config" yaml:"config"`
	Dbname   string `mapstructure:"db-name" json:"dbname" yaml:"db-name"`
	Username string `mapstructure:"username" json:"username" yaml:"username"`
	Password string `mapstructure:"password" json:"password" yaml:"password"`
}

type RedisOptions struct {
	Addr         string        `mapstructure:"addr" json:"addr" yaml:"addr"`
	Username     string        `mapstructure:"username" json:"username" yaml:"username"`
	DB           int           `mapstructure:"db" json:"db" yaml:"db"`
	Password     string        `mapstructure:"password" json:"password" yaml:"password"`
	PoolSize     int           `mapstructure:"poolSize" json:"poolSize" yaml:"poolSize"`
	DialTimeout  time.Duration `mapstructure:"dialTimeout" json:"dialTimeout" yaml:"dialTimeout"`
	ReadTimeout  time.Duration `mapstructure:"readTimeout" json:"readTimeout" yaml:"readTimeout"`
	WriteTimeout time.Duration `mapstructure:"writeTimeout" json:"writeTimeout" yaml:"writeTimeout"`
}
