package v2_gots_sdk

import (
	"path/filepath"
	"reflect"
	"strings"

	"github.com/iancoleman/strcase"
)

var templateFunc = `
export async function #FuncName#(query: #Query#): Promise<#Response#> {
    let res = await clientSdk.client.Method<any, AxiosResponse<#Response#, any>, any>('#Url#', {
        params: query,
    });
    return res.data;
}
`

var templateFuncWithBody = `
export async function #FuncName#(query: #Query#, data: #Payload#): Promise<#Response#> {
    let res = await clientSdk.client.Method<any, AxiosResponse<#Response#, any>, #Payload#>('#Url#', data, {
        params: query,
    });
    return res.data;
}
`

var templateImportHead = `import { AxiosInstance, AxiosResponse } from "axios";`
var templateClassApi = `
class ClientSdk {
    client!: AxiosInstance

}
const clientSdk = new ClientSdk()

export function SetClient(client: AxiosInstance) {
    clientSdk.client = client
}
`

func getStructName(data interface{}) string {
	if data == nil {
		return "any"
	}

	if reflect.ValueOf(data).Kind() == reflect.Ptr {
		return reflect.Indirect(reflect.ValueOf(data)).Type().Name()
	}
	return reflect.TypeOf(data).Name()
}

type Api struct {
	Method       string
	RelativePath string
	Payload      interface{}
	Response     interface{}
	Query        interface{}
	GroupPath    string
}

func (api *Api) replaceFuncName(template string) string {
	name := filepath.Join(api.GroupPath, api.RelativePath)
	name = strings.ReplaceAll(name, `\`, `/`)
	template = strings.ReplaceAll(template, "#Url#", name)

	fnname := ""
	funcname := filepath.Join(strings.ToLower(api.Method), name)
	funcs := strings.Split(funcname, `\`)

	for _, u := range funcs {
		fnname += strcase.ToCamel(u)
	}

	template = strings.ReplaceAll(template, "#FuncName#", fnname)

	return template
}

func (api *Api) GenerateTs() string {
	var template string

	if api.Payload != nil {
		template = templateFuncWithBody
	} else {
		template = templateFunc
	}

	template = strings.ReplaceAll(template, "Method", strings.ToLower(api.Method))
	template = strings.ReplaceAll(template, "#Query#", getStructName(api.Query))
	template = strings.ReplaceAll(template, "#Response#", getStructName(api.Response))
	template = strings.ReplaceAll(template, "#Payload#", getStructName(api.Payload))

	template = api.replaceFuncName(template)

	return template
}
