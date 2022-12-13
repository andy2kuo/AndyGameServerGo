package main

import (
	"fmt"
	"reflect"
	"strconv"
)

type Ainter interface {
	Test()
}

type A struct {
	I int    `default0:"4" default1:"42"`
	S string `default0:"Some String..." default1:"Some Other String..."`
}

func (A) Test() {

}

func setField(field reflect.Value, defaultVal string) error {

	if !field.CanSet() {
		return fmt.Errorf("can't set value")
	}

	switch field.Kind() {

	case reflect.Int:
		if val, err := strconv.ParseInt(defaultVal, 10, 64); err == nil {
			field.Set(reflect.ValueOf(int(val)).Convert(field.Type()))
		}
	case reflect.String:
		field.Set(reflect.ValueOf(defaultVal).Convert(field.Type()))
	}

	return nil
}

func Set(ptr Ainter, tag string) error {
	if reflect.TypeOf(ptr).Kind() != reflect.Ptr {
		return fmt.Errorf("not a pointer")
	}

	
	fmt.Printf("'%v'\n", reflect.TypeOf(ptr).Elem().Name())

	v := reflect.ValueOf(ptr).Elem()
	t := v.Type()

	for i := 0; i < t.NumField(); i++ {
		if defaultVal := t.Field(i).Tag.Get(tag); defaultVal != "-" {
			if err := setField(v.Field(i), defaultVal); err != nil {
				return err
			}

		}
	}

	return nil
}

func main() {

	ptr := &A{}

	if err := Set(ptr, "default0"); err == nil {
		//fmt.Printf("ptr.I=%d ptr.S=%s\n", ptr.I, ptr.S)
		// ptr.I=3 ptr.S=Some String...
	} else {
		fmt.Println(err)
	}

	if err := Set(ptr, "default1"); err == nil {
		//fmt.Printf("ptr.I=%d ptr.S=%s\n", ptr.I, ptr.S)
		// ptr.I=42 ptr.S=Some Other String...
	} else {
		fmt.Println(err)
	}

}
