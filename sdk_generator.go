package v2_gots_sdk

import (
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/pdcgo/v2_gots_sdk/pdc_api"
	"github.com/tkrajina/typescriptify-golang-structs/typescriptify"
)

type AddSdkFunc func(api *pdc_api.Api)

type ApiSdk struct {
	R     *gin.Engine
	Model *typescriptify.TypeScriptify
	toSdk AddSdkFunc
}

func (sdk *ApiSdk) GenerateSdkFunc(fname string) (createSdkJs func(), err error) {

	template, err := pdc_api.NewV2SdkTemplating(fname)

	if err != nil {
		return createSdkJs, err
	}

	sdk.toSdk = func(api *pdc_api.Api) {
		err := template.Register(api)
		if err != nil {
			panic(err)
		}

	}

	return func() {
		template.Save()
	}, nil

}

type RegisterFunc func(api *pdc_api.Api, handlers ...gin.HandlerFunc) gin.IRoutes

func (sdk *ApiSdk) Register(api *pdc_api.Api, handlers ...gin.HandlerFunc) gin.IRoutes {
	sdk.toSdk(api)

	return sdk.R.Handle(api.Method, api.RelativePath, handlers...)
}

type SdkGroup struct {
	sdk      *ApiSdk
	G        *gin.RouterGroup
	Basepath string
}

func (grp *SdkGroup) Register(api *pdc_api.Api, handlers ...gin.HandlerFunc) gin.IRoutes {
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

	var registfn RegisterFunc = func(api *pdc_api.Api, handlers ...gin.HandlerFunc) gin.IRoutes {
		api.GroupPath = relativePath

		sdk.toSdk(api)

		return r.Handle(api.Method, api.RelativePath, handlers...)
	}
	groupHandler(r, registfn)
}

func NewApiSdk(r *gin.Engine) *ApiSdk {
	sdk := &ApiSdk{
		Model: typescriptify.New(),
		toSdk: func(api *pdc_api.Api) {},
		R:     r,
	}
	sdk.Model.CreateInterface = true
	sdk.Model.CreateConstructor = false
	return sdk
}
