package api

import (
	"errors"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

// Config 用于映射 YAML 文件的结构
type Config struct {
	WxAppletPath string `yaml:"wx_applet_path"`
	WxOutputPath string `yaml:"wx_output_path"`
	TimeOut      int    `yaml:"time_out"`
}

// LoadConfig 读取 YAML 文件并返回配置
func LoadConfig(filePath string) (*Config, error) {
	if filePath == "" {
		return nil, errors.New("file path cannot be empty")
	}

	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
