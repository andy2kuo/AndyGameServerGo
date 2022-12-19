package config

import (
	"errors"
	"fmt"
	"os"
	"path"
	"reflect"

	"gopkg.in/ini.v1"
)

type ConfigBase interface {
	Name() string
}

func init() {
	os.MkdirAll("Config", 0755)
}

func defaultConfig(config ConfigBase) (default_file *ini.File) {
	default_file = ini.Empty()

	v := reflect.ValueOf(config).Elem()
	t := v.Type()

	for i := 0; i < t.NumField(); i++ {
		m_name := t.Field(i).Name

		section, err := default_file.NewSection(m_name)
		if err != nil {
			continue
		}

		if t.Field(i).Type.Kind() == reflect.Struct {
			for j := 0; j < t.Field(i).Type.NumField(); j++ {
				f := t.Field(i).Type.Field(j)
				if defaultVal := f.Tag.Get("default"); defaultVal != "-" && defaultVal != "" {
					section.NewKey(f.Name, defaultVal)
				} else {
					switch f.Type.Kind() {
					case reflect.Int:
						section.NewKey(f.Name, "0")
					case reflect.Int8:
						section.NewKey(f.Name, "0")
					case reflect.Int16:
						section.NewKey(f.Name, "0")
					case reflect.Int32:
						section.NewKey(f.Name, "0")
					case reflect.Int64:
						section.NewKey(f.Name, "0")
					case reflect.String:
						section.NewKey(f.Name, "Empty")
					}
				}
			}
		}
	}

	return default_file
}

func LoadConfig(config_data ConfigBase) error {
	if reflect.TypeOf(config_data).Kind() != reflect.Ptr {
		return fmt.Errorf("not a pointer")
	}

	var ini_path string = path.Join("Config", config_data.Name())
	var err error = nil
	if _, err = os.Stat(ini_path); errors.Is(err, os.ErrNotExist) {
		_file, _ := os.Create(ini_path)
		_file.Close()

		_ini_file := defaultConfig(config_data)
		_ini_file.SaveTo(ini_path)
	}

	var ini_file *ini.File
	ini_file, _ = ini.Load(ini_path)
	return ini_file.MapTo(config_data)
}
