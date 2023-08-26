package js_generator

import (
	"encoding/json"
	"io"
	"reflect"
	"strconv"
	"strings"
)

type JsGenerator struct {
	Models map[string]string
	Writer io.StringWriter
}

func NewJsGenerator(writer io.StringWriter) (*JsGenerator, error) {
	gen := JsGenerator{
		Models: map[string]string{},
		Writer: writer,
	}

	return &gen, nil
}

func (gen *JsGenerator) GenerateFromStruct(data interface{}, level int) (string, string, error) {
	level += 1

	tipes := reflect.TypeOf(data)
	values := reflect.ValueOf(data)

	switch tipes.Kind() {
	case reflect.String:
		return "`" + values.Interface().(string) + "`", "string", nil

	case reflect.Bool:
		dat := values.Interface().(bool)
		tipestr := "boolean"
		if dat {
			return "true", tipestr, nil
		}
		return "false", tipestr, nil

	case reflect.Float32:
		tipestr := "number"
		dat := values.Interface().(float32)
		return strconv.FormatFloat(float64(dat), 'f', 2, 32), tipestr, nil

	case reflect.Float64:
		tipestr := "number"
		dat := values.Interface().(float64)
		return strconv.FormatFloat(dat, 'f', 2, 32), tipestr, nil

	case reflect.Int:
		dat := values.Interface().(int)
		return strconv.Itoa(dat), "number", nil
	case reflect.Int16:
		dat := values.Interface().(int16)
		return strconv.Itoa(int(dat)), "number", nil
	case reflect.Int32:
		dat := values.Interface().(int32)
		return strconv.Itoa(int(dat)), "number", nil
	case reflect.Int64:
		dat := values.Interface().(int64)
		return strconv.Itoa(int(dat)), "number", nil
	case reflect.Int8:
		dat := values.Interface().(int8)
		return strconv.Itoa(int(dat)), "number", nil

	case reflect.Map:
		tipestr := []string{}

		valkey := reflect.Zero(tipes.Key()).Interface()
		_, tipe, err := gen.GenerateFromStruct(valkey, level)
		if err != nil {
			return "any", tipe, err
		}
		tipestr = append(tipestr,
			"[key: "+tipe+"]",
		)

		valvat := reflect.Zero(tipes.Elem()).Interface()
		_, tipe, err = gen.GenerateFromStruct(valvat, level)
		if err != nil {
			return "any", tipe, err
		}

		tipestr = append(tipestr, tipe)

		tipestring := "{" + strings.Join(tipestr, ": ") + "}"

		if values.IsNil() {

			return "{}", tipestring, nil
		}

		valstr, _ := json.Marshal(values.Interface())

		return string(valstr), tipestring, nil

	case reflect.Pointer:
		value := reflect.Indirect(values)

		if values.IsNil() {
			val := reflect.Zero(tipes.Elem()).Interface()
			valuestr, tipe, err := gen.GenerateFromStruct(val, level)
			tipe = tipe + " | undefined"
			return valuestr + " as " + tipe, tipe, err
		}

		valuedata := value.Interface()
		if valuedata == nil {
			valuedata = reflect.Zero(value.Type()).Interface()
		}

		valuestr, tipe, err := gen.GenerateFromStruct(valuedata, level)
		tipe = tipe + " | undefined"
		return valuestr + " as " + tipe, tipe, err

	case reflect.Slice:
		arrayValue := ArrayTs{}

		lendata := values.Len()
		var tipeArray string
		if lendata == 0 {

			elemzero := reflect.Zero(tipes.Elem()).Interface()
			value, tipe, err := gen.GenerateFromStruct(elemzero, level)
			tipeArray = tipe
			if err != nil {
				return "any", tipe, err
			}

			arrayValue = append(arrayValue, value)
		}

		for i := 0; i < lendata; i++ {
			value, tipe, err := gen.GenerateFromStruct(values.Index(i).Interface(), level)
			if err != nil {
				return "any", tipe, err
			}
			arrayValue = append(arrayValue, value)
			tipeArray = tipe
		}

		// val := values.Elem().Interface()

		// value, err := GenerateFromStruct(val)
		// data = append(data, value)
		return arrayValue.GenerateTs(level), "Array<" + tipeArray + ">", nil

	case reflect.Struct:

		objectVal := ObjectTs{}
		objectType := InterfaceTs{}

		for c := 0; c < tipes.NumField(); c++ {
			item := ObjectTsItem{}

			tipe := tipes.Field(c)
			key := tipe.Tag.Get("json")
			if key == "" {
				continue
			}

			// getting value
			val := values.Field(c)
			valstr, tipestr, err := gen.GenerateFromStruct(val.Interface(), 0)
			if err != nil {
				return "", objectType.GenerateTs(), err
			}

			item.Val = valstr
			item.Key = key

			objectVal = append(objectVal, &item)
			objectType = append(objectType, InterfaceTsItem{
				Key: key,
				Val: tipestr,
			})
		}

		name := tipes.Name()
		importObject := gen.Models[name]
		if importObject == "" {
			importObject = objectType.GenerateTs()
			gen.Models[name] = importObject
			gen.Writer.WriteString("export interface " + name + " " + importObject + "\n\n")
		}
		return objectVal.GenerateTs(level), name, nil

	}

	return `null`, "null", nil

}

type InterfaceTsItem struct {
	Key string
	Val string
}

type InterfaceTs []InterfaceTsItem

func (obj InterfaceTs) GenerateTs() string {
	hasil := []string{
		"{",
	}

	for _, data := range obj {
		hasil = append(hasil, "\t"+data.Key+": "+data.Val)
	}

	hasil = append(hasil, "}")
	return strings.Join(hasil, "\n")
}

type ObjectTsItem struct {
	Key    string
	Val    string
	Suffix string
	Level  int
}

func (obj *ObjectTsItem) GetLevel() string {
	tabstring := ""
	if obj.Level == 0 {
		tabstring = ""
	} else {
		tab := make([]string, obj.Level)
		tabstring = strings.Join(tab, "\t")
	}
	return tabstring
}

type ObjectTs []*ObjectTsItem

func (obj ObjectTs) GenerateTs(level int) string {
	tabstring := ""
	if level == 0 {
		tabstring = ""
	} else {
		tab := make([]string, level)
		tabstring = strings.Join(tab, "\t")
	}

	hasil := []string{}
	for _, data := range obj {
		suffix := data.Suffix
		if suffix != "" {
			suffix = " " + suffix
		}

		vals := strings.Split(data.Val, "\n")
		newvals := make([]string, len(vals))

		for ind, val := range vals {
			if ind == 0 {
				newvals[ind] = val
				continue
			}
			newvals[ind] = tabstring + val
		}
		stringsvals := strings.Join(newvals, "\n")

		hasil = append(hasil, tabstring+"\t"+data.Key+": "+stringsvals+suffix)
	}

	return "{\n" + strings.Join(hasil, ",\n") + "\n" + tabstring + "}"
}

type ArrayTs []string

func (arr ArrayTs) GenerateTs(level int) string {
	hasil := []string{"["}
	if level <= 0 {
		level = 2
	}

	for _, val := range arr {
		tabs := strings.Join(make([]string, level), "\t")
		hasil = append(hasil, tabs+val)
	}
	tabs := strings.Join(make([]string, level-1), "\t")
	hasil = append(hasil, tabs+"]")

	return strings.Join(hasil, "\n\t")
}
