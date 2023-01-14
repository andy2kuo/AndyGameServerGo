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

var ErrCreateNewConfig error = errors.New("load config by create new one")
var ErrLoadConfig error = errors.New("load config on path")
var tempConfigMap map[string]IConfig

type IConfig interface {
	Name() string
}

func init() {
	os.MkdirAll("Config", 0755)
	tempConfigMap = make(map[string]IConfig)
}

func createConfig(config_data interface{}) (config_file *ini.File) {
	config_file = ini.Empty()

	data_value := reflect.ValueOf(config_data)
	data_type := data_value.Type()

	for i := 0; i < data_type.Elem().NumField(); i++ {
		section_value := loadSection(data_type.Elem().Field(i).Name, data_type.Elem().Field(i).Type, config_file)
		data_value.Elem().Field(i).Set(section_value.Elem())
	}

	return config_file
}

func loadConfig(config_data interface{}, config_file *ini.File) {
	data_value := reflect.ValueOf(config_data)
	data_type := data_value.Type()

	for i := 0; i < data_type.Elem().NumField(); i++ {
		section_value := loadSection(data_type.Elem().Field(i).Name, data_type.Elem().Field(i).Type, config_file)
		data_value.Elem().Field(i).Set(section_value.Elem())
	}
}

func loadSection(section_name string, section_type reflect.Type, config_file *ini.File) reflect.Value {
	new_section_value := reflect.New(section_type)
	var config_section_info *ini.Section
	var _is_section_exist bool = config_file.HasSection(section_name)

	if _is_section_exist {
		config_section_info, _ = config_file.GetSection(section_name)
	} else {
		config_section_info, _ = config_file.NewSection(section_name)
	}

	for i := 0; i < new_section_value.Elem().NumField(); i++ {
		_field_value := new_section_value.Elem().Field(i)
		_field_type := _field_value.Type()

		var _input_value string
		var _is_key_exist bool = false

		if _is_section_exist {
			_is_key_exist = config_section_info.HasKey(section_type.Field(i).Name)
		}

		if _is_key_exist {
			_key := config_section_info.Key(section_type.Field(i).Name)
			_input_value = _key.String()
		} else {
			_input_value = new_section_value.Elem().Type().Field(i).Tag.Get("default")
		}

		var err error = nil
		switch _field_type.Kind() {
		case reflect.Bool:
			var val bool
			if val, err = strconv.ParseBool(_input_value); err != nil {
				val = false
			}

			_field_value.Set(reflect.ValueOf(val).Convert(_field_type))
			if !_is_key_exist {
				config_section_info.NewKey(section_type.Field(i).Name, fmt.Sprint(val))
			}

		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			var val int64
			if val, err = strconv.ParseInt(_input_value, 10, 64); err != nil {
				val = 0
			}

			_field_value.Set(reflect.ValueOf(val).Convert(_field_type))
			if !_is_key_exist {
				config_section_info.NewKey(section_type.Field(i).Name, fmt.Sprint(val))
			}

		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			var val uint64
			if val, err = strconv.ParseUint(_input_value, 10, 64); err != nil {
				val = 0
			}

			_field_value.Set(reflect.ValueOf(val).Convert(_field_type))
			if !_is_key_exist {
				config_section_info.NewKey(section_type.Field(i).Name, fmt.Sprint(val))
			}

		case reflect.Float32, reflect.Float64:
			var val float64
			if val, err = strconv.ParseFloat(_input_value, 64); err != nil {
				val = 0
			}

			_field_value.Set(reflect.ValueOf(val).Convert(_field_type))
			if !_is_key_exist {
				config_section_info.NewKey(section_type.Field(i).Name, fmt.Sprint(val))
			}

		case reflect.String:
			if _input_value == "" || _input_value == "-" {
				_input_value = "empty"
			}
			_field_value.Set(reflect.ValueOf(_input_value))
			if !_is_key_exist {
				config_section_info.NewKey(section_type.Field(i).Name, _input_value)
			}

		case reflect.Struct:
			switch _field_type {
			case reflect.TypeOf(time.Time{}):
				var val int64
				if val, err = strconv.ParseInt(_input_value, 10, 64); err != nil {
					val = time.Now().Unix()
				}

				_field_value.Set(reflect.ValueOf(time.Unix(val, 0)))
				if !_is_key_exist {
					config_section_info.NewKey(section_type.Field(i).Name, fmt.Sprint(val))
				}
			}
		}
	}
	return new_section_value
}

func GetConfig(env string, config_data IConfig) (err error) {
	if reflect.TypeOf(config_data).Kind() != reflect.Ptr {
		return fmt.Errorf("not a pointer")
	}

	// 檢查此設定檔是否已經讀取過
	if _, isLoaded := tempConfigMap[config_data.Name()]; isLoaded {
		config_data = tempConfigMap[config_data.Name()]
		return ErrLoadConfig
	}

	env_path := fmt.Sprintf("Config/%v", env)
	os.MkdirAll(env_path, 0755)

	var ini_name string = fmt.Sprintf("%v.%v", config_data.Name(), "ini")
	var ini_path string = path.Join(env_path, ini_name)
	var ini_file *ini.File

	if _, err = os.Stat(ini_path); errors.Is(err, os.ErrNotExist) {
		_file, _ := os.Create(ini_path)
		_file.Close()

		ini_file = createConfig(config_data)
		ini_file.SaveTo(ini_path)
		err = ErrCreateNewConfig
	} else {
		ini_file, err = ini.Load(ini_path)
		if err != nil {
			return err
		}

		loadConfig(config_data, ini_file)
		ini_file.SaveTo(ini_path)
		err = ErrLoadConfig
	}

	tempConfigMap[config_data.Name()] = config_data
	return err
}

func IsCreateNew(err error) bool {
	return errors.Is(err, ErrCreateNewConfig)
}

func IsLoadOnPath(err error) bool {
	return errors.Is(err, ErrLoadConfig)
}
