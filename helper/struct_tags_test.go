package helper

import (
	"reflect"
	"testing"
)

type cat struct {
	Name string
	Type int `json:"type" id:"100"`
}

func Test_readTag(t *testing.T) {

	typeOfCat := reflect.TypeOf(cat{})

	if catType, ok := typeOfCat.FieldByName("Type"); ok {
		t.Log(catType.Tag.Get("json"), catType.Tag.Get("id"))
	}
}


/*
http://c.biancheng.net/view/111.html

任意值通过 `reflect.TypeOf()` 获得反射对象信息后，如果是结构体，
可以通过反射值对象（reflect.Type) 的 NumField() 和 Field() 方法
获得结构体成员的详细信息。

 */
func Test_structInfo(t *testing.T) {

	ins := cat{Name: "mimi", Type: 1}

	// 获取结构体实例的反射类型
	typeOfCat := reflect.TypeOf(ins)

	// 遍历结构体的所有成员
	for i := 0; i < typeOfCat.NumField(); i++ {
		// 获取每个成员的结构体字段类型
		fieldType := typeOfCat.Field(i)
		// 输出成员名和 tag
		t.Logf("name => %v, tag => '%v'\n", fieldType.Name, fieldType.Tag)
	}

	// http://c.biancheng.net/view/4407.html
	var a int
	typeOfA := reflect.TypeOf(a)
	// typeOfA.Name => int, Kind => int
	t.Logf("typeOfA.Name => %s, Kind => %v\n", typeOfA.Name(), typeOfA.Kind())

	ins2 := &cat{}

	typeOfCat2 := reflect.TypeOf(ins2)
	// Name => '', Kind => 'ptr'
	t.Logf("Name => '%v', Kind => '%v'\n", typeOfCat2.Name(), typeOfCat2.Kind())

	elemOfCat2 := typeOfCat2.Elem()
	// Name => 'cat', Kind => 'struct'
	t.Logf("Name => '%v', Kind => '%v'\n", elemOfCat2.Name(), elemOfCat2.Kind())


}

/**
reflect.ValueOf 返回 reflect.Value 类型，包含有 rawValue 的值信息。
reflect.Value 与原值间可以通过值包装和值获取互相转化。
reflect.Value 是一些反射操作的重要类型，如反射调用函数。
 */
func Test_valueInfo(t *testing.T) {

	var a int
	a = 1024
	valueOfA := reflect.ValueOf(a)

	canInterface := valueOfA.CanInterface()
	t.Logf("canInterface => '%v'\n", canInterface)

	var getA int
	getA = valueOfA.Interface().(int)

	var getA2 int64 = int64(valueOfA.Int())

	t.Logf("getA => %v, getA2 => %v\n", getA, getA2)
}

func Test_GetStructValue(t *testing.T) {
	type dummy struct {
		a int
		b string
		// 嵌入字段
		float32
		bool

		next *dummy
	}

	d := reflect.ValueOf(dummy{
		next : &dummy{},
	})

	t.Log("numfield =>", d.NumField())

	// 获取索引为 2 的字段 （float32 ）
	floatField := d.Field(2)

	t.Log("Field =>", floatField.Type())

	t.Log("b: findByName => ", d.FieldByName("b").Type())

	t.Log("FieldByIndex =>", d.FieldByIndex([]int{4, 0}).Type())
}