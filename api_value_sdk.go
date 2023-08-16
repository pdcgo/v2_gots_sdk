package v2_gots_sdk

import (
	"encoding/json"
	"io"
	"os"
	"reflect"
	"strconv"
	"strings"
)

func (api *Api) GenerateClientValueSdk() {

}

type JsGenerator struct {
	Base      string
	Models    map[string]string
	ApiValues ObjectTs
	Writer    io.StringWriter
}

func NewJsGenerator(base string) (*JsGenerator, func() error, error) {
	f, err := os.OpenFile(base, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
	if err != nil {
		return nil, func() error { return nil }, err
	}

	gen := JsGenerator{
		Base:      base,
		Models:    map[string]string{},
		Writer:    f,
		ApiValues: ObjectTs{},
	}

	return &gen, func() error {

		f.WriteString("\n\n")

		f.WriteString("export const client = ")
		f.WriteString(gen.ApiValues.GenerateTs())

		f.Close()
		return nil
	}, nil
}

func (gen *JsGenerator) AddApi(api *Api) {

	apiObject := ObjectTs{}

	level := 2

	query := ObjectTsItem{
		Key:   "query",
		Val:   "undefined",
		Level: level,
	}
	if api.Query != nil {
		val, _, err := gen.GenerateFromStruct(api.Query, level+1)
		if err != nil {
			panic(err)
		}
		query.Val = val
	}

	response := ObjectTsItem{
		Key:   "response",
		Val:   "undefined",
		Level: level,
	}
	if api.Response != nil {
		val, _, err := gen.GenerateFromStruct(api.Response, level+1)
		if err != nil {
			panic(err)
		}
		response.Val = val
	}

	payload := ObjectTsItem{
		Key:   "body",
		Val:   "undefined",
		Level: level,
	}
	if api.Payload != nil {

		val, _, err := gen.GenerateFromStruct(api.Payload, level+1)
		if err != nil {
			panic(err)
		}
		payload.Val = val
	}

	apiObject = append(apiObject, query)
	apiObject = append(apiObject, response)
	apiObject = append(apiObject, payload)

	apiObject = append(apiObject, ObjectTsItem{
		Key:    "method",
		Val:    `'` + strings.ToUpper(api.Method) + `'`,
		Suffix: "as const",
		Level:  level,
	})

	apiObject = append(apiObject, ObjectTsItem{
		Key:    "url",
		Val:    `'` + strings.TrimPrefix(api.GetFullUriPath(), `/`) + `'`,
		Suffix: "as const",
		Level:  level,
	})

	gen.ApiValues = append(gen.ApiValues, ObjectTsItem{
		Key:   api.GetKeyName(),
		Val:   apiObject.GenerateTs(),
		Level: 2,
	})

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
			valstr, tipestr, err := gen.GenerateFromStruct(val.Interface(), level)
			if err != nil {
				return "", objectType.GenerateTs(), err
			}

			item.Val = valstr
			item.Key = key
			item.Level = level

			objectVal = append(objectVal, item)
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

		return objectVal.GenerateTs(), name, nil

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
	tabs := make([]string, obj.Level)
	return strings.Join(tabs, "\t")
}

type ObjectTs []ObjectTsItem

func (obj ObjectTs) GenerateTs() string {
	hasil := []string{}
	var level int = 1
	for _, data := range obj {
		suffix := data.Suffix
		if suffix != "" {
			suffix = " " + suffix
		}
		hasil = append(hasil, data.GetLevel()+data.Key+": "+data.Val+suffix)
		level = data.Level
	}
	tab := make([]string, level-1)
	return "{\n" + strings.Join(hasil, ",\n") + "\n" + strings.Join(tab, "\t") + "}"
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
