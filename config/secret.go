package config

// vol中对应的秘钥配置
type Secret struct {
	User string `mapstructure:"User" json:"User" yaml:"User"`
}
