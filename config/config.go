package config

import (
	"fmt"
	"os"
	"path"
	"reflect"
	"strconv"

	"gopkg.in/ini.v1"
)

type ConfigBase interface {
	DirPath() string
	File() string
	Section() string
}

func setField(field reflect.Value, value string) error {

	if !field.CanSet() {
		return fmt.Errorf("'%v' can't set value", field)
	}

	switch field.Kind() {

	case reflect.Int:
		if val, err := strconv.ParseInt(value, 10, 64); err == nil {
			field.Set(reflect.ValueOf(int(val)).Convert(field.Type()))
		}
	case reflect.String:
		field.Set(reflect.ValueOf(value).Convert(field.Type()))
	}

	return nil
}

func DefaultConfig(config ConfigBase) error {
	if reflect.TypeOf(config).Kind() != reflect.Ptr {
		return fmt.Errorf("'%v' not a pointer", reflect.TypeOf(config).Name())
	}

	v := reflect.ValueOf(config).Elem()
	t := v.Type()

	for i := 0; i < t.NumField(); i++ {

		if defaultVal := t.Field(i).Tag.Get("default"); defaultVal != "-" {
			if err := setField(v.Field(i), defaultVal); err != nil {
				return err
			}

		}
	}

	return nil
}

func LoadConfig(config ConfigBase) error {
	if reflect.TypeOf(config).Kind() != reflect.Ptr {
		return fmt.Errorf("not a pointer")
	}

	return nil
}

func (config *Config) Load() (err error) {
	config.data_set = make(map[string]map[string]string)
	err = nil
	_cfg, err := ini.Load(config.data_path)
	if err != nil {
		return err
	}

	_sections := _cfg.Sections()
	for _section_index := range _sections {
		_section := _sections[_section_index]

		_keys := _section.Keys()
		for _key_index := range _keys {
			_key := _keys[_key_index]
			_key.Duration()
		}
	}
}

func (config *Config) Get(field string) {

}

// 建立新的設定檔
func NewConfig() (config *Config) {
	config = &Config{
		data_set:  make(map[string]map[string]string),
		data_path: default_config_path,
	}

	return config
}
