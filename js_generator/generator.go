package js_generator

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type CacheModel struct {
	name  string
	value string
}

type JsGenerator struct {
	Models map[string]*CacheModel
	Writer io.StringWriter
}

func NewJsGenerator(writer io.StringWriter) (*JsGenerator, error) {
	gen := JsGenerator{
		Models: map[string]*CacheModel{},
		Writer: writer,
	}

	return &gen, nil
}

func (gen *JsGenerator) IterateFieldStruct(data interface{}) (ObjectTs, InterfaceTs, error) {
	tipes := reflect.TypeOf(data)
	values := reflect.ValueOf(data)

	objectVal := ObjectTs{}
	objectType := InterfaceTs{}

	for c := 0; c < tipes.NumField(); c++ {

		tipe := tipes.Field(c)
		val := values.Field(c)

		if tipe.Anonymous {
			var anonim any

			if tipe.Type.Kind() == reflect.Pointer {
				isNil := val.IsNil()
				elem := tipe.Type.Elem()
				val = reflect.Indirect(val)

				if isNil {
					anonim = reflect.Zero(elem).Interface()
				} else {
					anonim = val.Interface()
				}

			} else {
				anonim = val.Interface()
			}

			aobj, atipe, err := gen.IterateFieldStruct(anonim)

			if err != nil {
				return objectVal, objectType, err
			}

			objectVal = append(objectVal, aobj...)
			objectType = append(objectType, atipe...)

		}

		key := tipe.Tag.Get("json")
		if key == "" || key == "-" {
			continue
		}

		keys := strings.Split(key, ",")
		key = keys[0]

		fieldname := val.Type().Name()
		if gen.Models[fieldname] != nil {
			return objectVal, objectType, nil
		}

		// getting value
		valstr, tipestr, err := gen.GenerateFromStruct(val.Interface(), 0)
		if err != nil {
			return objectVal, objectType, nil
		}

		item := ObjectTsItem{
			Val: valstr,
			Key: key,
		}

		objectVal = append(objectVal, &item)
		objectType = append(objectType, InterfaceTsItem{
			Key: key,
			Val: tipestr,
		})
	}

	return objectVal, objectType, nil

}

type EnumString interface {
	EnumList() []string
}

type ProxyStruct interface {
	ProxyStruct() interface{}
}

func convertUnion(datas []string) string {
	hasil := []string{}

	for _, val := range datas {
		hasil = append(hasil, fmt.Sprintf(`"%s"`, val))
	}
	hasil = append(hasil,
		`({} & string)`,
	)
	return strings.Join(hasil, " | ")
}

func (gen *JsGenerator) GenerateFromStruct(data interface{}, level int) (string, string, error) {
	level += 1

	if data == nil {
		return "{}", "any", nil
	}

	tipes := reflect.TypeOf(data)
	values := reflect.ValueOf(data)

	switch tipes.Kind() {
	case reflect.String:
		dd := values.Interface()

		// check enum
		enumdat, okc := dd.(EnumString)
		if okc {
			enumList := enumdat.EnumList()
			union := convertUnion(enumList)
			return "`" + enumList[0] + "`", union, nil
		}

		cc, ok := dd.(string)
		if !ok {
			return "``", "string", nil
		}

		return "`" + cc + "`", "string", nil

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

		// log.Println(values.CanFloat(), "asdasdasd", values.Float())
		dat := values.Float()
		return strconv.FormatFloat(dat, 'f', 2, 32), tipestr, nil

	case reflect.Uint:
		dat := values.Uint()
		return strconv.FormatUint(dat, 10), "number", nil
	case reflect.Uint8:
		dat := values.Uint()
		return strconv.FormatUint(dat, 10), "number", nil
	case reflect.Uint16:
		dat := values.Uint()
		return strconv.FormatUint(dat, 10), "number", nil
	case reflect.Uint32:
		dat := values.Uint()
		return strconv.FormatUint(dat, 10), "number", nil
	case reflect.Uint64:
		dat := values.Uint()
		return strconv.FormatUint(dat, 10), "number", nil
	case reflect.Int:
		dat := values.Int()
		return strconv.Itoa(int(dat)), "number", nil
	case reflect.Int16:
		dat := values.Int()
		return strconv.Itoa(int(dat)), "number", nil
	case reflect.Int32:
		dat := values.Int()
		return strconv.Itoa(int(dat)), "number", nil
	case reflect.Int64:
		dat := values.Int()
		return strconv.Itoa(int(dat)), "number", nil
	case reflect.Int8:
		dat := values.Int()
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
			if valuestr != "" {
				valuestr = valuestr + " as " + tipe
			}
			return valuestr, tipe, err
		}

		valuedata := value.Interface()
		if valuedata == nil {
			valuedata = reflect.Zero(value.Type()).Interface()
		}

		valuestr, tipe, err := gen.GenerateFromStruct(valuedata, level)
		tipe = tipe + " | undefined"
		if valuestr != "" {
			valuestr = valuestr + " as " + tipe
		}
		return valuestr, tipe, err

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

		// remove undefined jare toni
		tipeArray = strings.ReplaceAll(tipeArray, " | undefined", "")

		tipeArray = "Array<" + tipeArray + ">"

		return arrayValue.GenerateTs(level) + " as " + tipeArray, tipeArray, nil

	case reflect.Struct:
		switch dd := data.(type) {
		case time.Time:
			return "`2021-12-01T07:00:00+07:00`", "string", nil
		case ProxyStruct:
			data = dd.ProxyStruct()
			return gen.GenerateFromStruct(data, level)
		}

		name := tipes.Name()
		name = DetectGeneric(name)

		importObject := gen.Models[name]

		if importObject != nil {
			return importObject.value, importObject.name, nil
		} else {
			gen.Models[name] = &CacheModel{
				name: name,
			}
		}

		objectVal, objectType, err := gen.IterateFieldStruct(data)

		if err != nil {
			return "", "", err
		}

		importObject = &CacheModel{
			name:  name,
			value: objectVal.GenerateTs(level),
		}
		gen.Models[name] = importObject
		tipestr := objectType.GenerateTs()
		gen.Writer.WriteString("export interface " + name + " " + tipestr + "\n\n")
		return importObject.value, name, nil

	}

	return `null`, "null", nil

}

func DetectGeneric(name string) string {
	rex := regexp.MustCompile(`\[(.+)\]`)
	dd := rex.FindStringSubmatch(name)

	generic := ""

	if len(dd) > 0 {

		matchend := dd[1]

		multidata := strings.Split(matchend, ",")
		log.Println(multidata)
		for _, gene := range multidata {
			datas := strings.Split(gene, ".")
			c := len(datas)
			generic += datas[c-1]
		}

		names := strings.Split(name, "[")
		name = names[0]
	}

	return name + generic
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
	mapunique := map[string]bool{}
	for _, data := range obj {
		if mapunique[data.Key] {
			continue
		} else {
			mapunique[data.Key] = true
		}
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
	mapunique := map[string]bool{}
	for _, data := range obj {
		if mapunique[data.Key] {
			continue
		} else {
			mapunique[data.Key] = true
		}

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
		if val == "" {
			continue
		}
		tabs := strings.Join(make([]string, level), "\t")
		hasil = append(hasil, tabs+val)
	}
	tabs := strings.Join(make([]string, level-1), "\t")
	hasil = append(hasil, tabs+"]")

	return strings.Join(hasil, "\n\t")
}
