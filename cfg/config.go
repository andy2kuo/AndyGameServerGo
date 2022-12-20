package config

import (
	"errors"
	"fmt"
	"os"
	"path"
	"reflect"
	"strconv"
	"time"

	"gopkg.in/ini.v1"
)

type ConfigBase interface {
	Name() string
}

func init() {
	os.MkdirAll("Config", 0755)
}

func newDefaultConfig(config_data interface{}) (config_file *ini.File) {
	config_file = ini.Empty()

	data_value := reflect.ValueOf(config_data)
	data_type := data_value.Type()

	for i := 0; i < data_type.Elem().NumField(); i++ {
		section_value := newDefaultSection(data_type.Elem().Field(i).Type, config_file)
		data_value.Elem().Field(i).Set(section_value.Elem())
	}

	return config_file
}

func newDefaultSection(section_type reflect.Type, config_file *ini.File) reflect.Value {
	new_section_value := reflect.New(section_type)
	config_section_info, _ := config_file.NewSection(section_type.Name())

	for i := 0; i < new_section_value.Elem().NumField(); i++ {
		_field_value := new_section_value.Elem().Field(i)
		_field_type := _field_value.Type()
		_default_value := new_section_value.Elem().Type().Field(i).Tag.Get("default")

		var err error = nil
		switch _field_type.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			var val int64
			if val, err = strconv.ParseInt(_default_value, 10, 64); err != nil {
				val = 0
			}

			_field_value.Set(reflect.ValueOf(val).Convert(_field_type))
			config_section_info.NewKey(section_type.Field(i).Name, fmt.Sprint(val))

		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			var val uint64
			if val, err = strconv.ParseUint(_default_value, 10, 64); err != nil {
				val = 0
			}

			_field_value.Set(reflect.ValueOf(val).Convert(_field_type))
			config_section_info.NewKey(section_type.Field(i).Name, fmt.Sprint(val))

		case reflect.Float32, reflect.Float64:
			var val float64
			if val, err = strconv.ParseFloat(_default_value, 64); err != nil {
				val = 0
			}

			_field_value.Set(reflect.ValueOf(val).Convert(_field_type))
			config_section_info.NewKey(section_type.Field(i).Name, fmt.Sprint(val))

		case reflect.String:
			if _default_value == "" || _default_value == "-" {
				_default_value = "empty"
			}
			_field_value.Set(reflect.ValueOf(_default_value))
			config_section_info.NewKey(section_type.Field(i).Name, _default_value)

		case reflect.Struct:
			switch _field_type {
			case reflect.TypeOf(time.Time{}):
				var val int64
				if val, err = strconv.ParseInt(_default_value, 10, 64); err != nil {
					val = time.Now().Unix()
				}

				_field_value.Set(reflect.ValueOf(time.Unix(val, 0)))
				config_section_info.NewKey(section_type.Field(i).Name, fmt.Sprint(val))
			}
		}
	}
	return new_section_value
}

func LoadConfig(config_data ConfigBase) error {
	if reflect.TypeOf(config_data).Kind() != reflect.Ptr {
		return fmt.Errorf("not a pointer")
	}

	var ini_path string = path.Join("Config", config_data.Name())
	var err error = nil
	var ini_file *ini.File
	if _, err = os.Stat(ini_path); errors.Is(err, os.ErrNotExist) {
		_file, _ := os.Create(ini_path)
		_file.Close()

		ini_file = newDefaultConfig(config_data)
		ini_file.SaveTo(ini_path)
	} else {
		ini_file, _ = ini.Load(ini_path)
	}

	return nil
}
