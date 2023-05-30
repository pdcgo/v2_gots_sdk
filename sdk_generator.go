package v2_gots_sdk

import (
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/tkrajina/typescriptify-golang-structs/typescriptify"
)

type AddSdkFunc func(api *Api)

type ApiSdk struct {
	R     *gin.Engine
	Model *typescriptify.TypeScriptify
	toSdk AddSdkFunc
}

func (sdk *ApiSdk) GenerateSdkFunc(fname string) (createSdkJs func()) {

	funcscripts := []string{}

	sdk.toSdk = func(api *Api) {
		query := api.Query
		payload := api.Payload
		response := api.Response

		if query != nil {
			sdk.Model.Add(query)
		}
		if payload != nil {
			sdk.Model.Add(payload)
		}
		if response != nil {
			sdk.Model.Add(response)
		}
		funcscripts = append(funcscripts, api.GenerateTs())
	}

	return func() {
		basepath := filepath.Join(fname)
		os.Remove(basepath)

		f, err := os.OpenFile(basepath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0755)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()

		model, _ := sdk.Model.Convert(map[string]string{})

		f.WriteString(templateImportHead)
		f.WriteString("\n\n")
		f.WriteString(model)
		f.WriteString("\n")
		f.WriteString(templateClassApi)
		f.WriteString("\n")
		f.WriteString(strings.Join(funcscripts, "\n"))
	}

}

type RegisterFunc func(api *Api, handlers ...gin.HandlerFunc) gin.IRoutes

func (sdk *ApiSdk) Register(api *Api, handlers ...gin.HandlerFunc) gin.IRoutes {
	sdk.toSdk(api)

	return sdk.R.Handle(api.Method, api.RelativePath, handlers...)
}

type SdkGroup struct {
	sdk      *ApiSdk
	G        *gin.RouterGroup
	Basepath string
}

func (grp *SdkGroup) Register(api *Api, handlers ...gin.HandlerFunc) gin.IRoutes {
	api.GroupPath = grp.Basepath
	grp.sdk.toSdk(api)
	return grp.G.Handle(api.Method, api.RelativePath, handlers...)
}

func (grp *SdkGroup) Group(path string) *SdkGroup {
	base, _ := url.JoinPath(grp.Basepath, path)
	newGroup := SdkGroup{
		sdk:      grp.sdk,
		G:        grp.G.Group(path),
		Basepath: base,
	}

	return &newGroup
}

func (sdk *ApiSdk) Group(relativePath string) *SdkGroup {
	newGroup := SdkGroup{
		sdk:      sdk,
		G:        sdk.R.Group(relativePath),
		Basepath: relativePath,
	}

	return &newGroup
}

func (sdk *ApiSdk) RegisterGroup(relativePath string, groupHandler func(group *gin.RouterGroup, register RegisterFunc)) {
	r := sdk.R.Group(relativePath)

	var registfn RegisterFunc = func(api *Api, handlers ...gin.HandlerFunc) gin.IRoutes {
		api.GroupPath = relativePath

		sdk.toSdk(api)

		return r.Handle(api.Method, api.RelativePath, handlers...)
	}
	groupHandler(r, registfn)
}

func NewApiSdk(r *gin.Engine) *ApiSdk {
	sdk := &ApiSdk{
		Model: typescriptify.New(),
		toSdk: func(api *Api) {},
		R:     r,
	}
	sdk.Model.CreateInterface = true
	sdk.Model.CreateConstructor = false
	return sdk
}
