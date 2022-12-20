package main

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"time"

	config "github.com/andy2kuo/AndyGameServerGo/cfg"
)

type Finter interface {
	Test()
}

type F struct {
	A int       `default:"4"`
	S string    `default:"Some String..."`
	T time.Time `default:"1671076216"`
}

func (F) Test() {

}

func setField(field reflect.Value, val string) error {

	if !field.CanSet() {
		return fmt.Errorf("can't set value")
	}

	switch field.Kind() {

	case reflect.Int:
		if val, err := strconv.ParseInt(val, 10, 64); err == nil {
			field.Set(reflect.ValueOf(int(val)).Convert(field.Type()))
		}
	case reflect.String:
		field.Set(reflect.ValueOf(val).Convert(field.Type()))
	case reflect.Struct:
		switch field.Type() {
		case reflect.TypeOf(time.Time{}):
			if val, err := strconv.ParseInt(val, 10, 64); err == nil {
				_time := time.Unix(val, 0)
				field.Set(reflect.ValueOf(_time).Convert(field.Type()))
			}
		}
	}

	return nil
}

func Set(ptr Finter, tag string) error {
	if reflect.TypeOf(ptr).Kind() != reflect.Ptr {
		return fmt.Errorf("not a pointer")
	}

	fmt.Printf("'%v'\n", reflect.TypeOf(ptr).Elem().Name())

	v := reflect.ValueOf(ptr).Elem()
	t := v.Type()

	for i := 0; i < t.NumField(); i++ {
		if val := t.Field(i).Tag.Get(tag); val != "-" {
			if err := setField(v.Field(i), val); err != nil {
				return err
			}
		}
	}

	return nil
}

type A struct {
	ATest1 int64     `default:"9223372036854775807"`
	ATest2 time.Time `default:"1671529441"`
}

func (a A) GetTime() time.Time {
	return time.Unix(a.ATest1, 0)
}

type B struct {
	BTest1 int    `default:"2"`
	BTest2 string `default:"-"`
}

type TestConfig struct {
	A A
	B B
}

func (TestConfig) Name() string {
	return "Test.ini"
}

func main() {
	_cfg := &TestConfig{}

	config.LoadConfig(_cfg)
	return
	t1 := &TestConfig{
		A: A{
			ATest1: 100104789,
		},
	}
	fmt.Println("First", t1)
	t1_type := reflect.TypeOf(t1).Elem()
	t1_value := reflect.ValueOf(t1).Elem()

	t2_type := t1_type.Field(0).Type
	fmt.Println("01", t2_type)
	fmt.Println("02", reflect.New(t2_type).Type())
	fmt.Println("03", reflect.New(t2_type).Type().Kind())
	fmt.Println("04", reflect.New(t2_type).Type().Elem())
	fmt.Println("05", reflect.New(t2_type).Type().Elem().Kind())

	new_value := reflect.New(t2_type)
	fmt.Println("1", new_value)
	fmt.Println("2", new_value.Elem())
	fmt.Println("3", new_value.Elem().Field(0))
	fmt.Println("4", t1_value)

	val, _ := strconv.ParseInt("123", 10, 64)
	new_value.Elem().Field(0).Set(reflect.ValueOf(val).Convert(new_value.Elem().Field(0).Type()))
	//new_value.Set(reflect.ValueOf(int64(123)))

	t1_value.Field(0).Set(new_value.Elem())
	fmt.Println(t1_value)

	fmt.Println("Final", t1)
	return
	t := &TestConfig{}
	fmt.Println(t)
	err := config.LoadConfig(t)
	fmt.Println(err)
	fmt.Println(t)
	fmt.Println(t.A.GetTime())

	return
	_file, err := os.OpenFile("test.txt", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0755)
	fmt.Println(err)

	_info, err := _file.Stat()
	fmt.Println(err)
	fmt.Println(_info.Size())

	return
	ptr := &F{}

	if err := Set(ptr, "default"); err == nil {
		//fmt.Printf("ptr.I=%d ptr.S=%s\n", ptr.I, ptr.S)
		// ptr.I=3 ptr.S=Some String...
	} else {
		fmt.Println(err)
	}

	fmt.Println(ptr)
}
