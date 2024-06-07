package v2_gots_sdk_test

import (
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/pdcgo/v2_gots_sdk"
	"github.com/pdcgo/v2_gots_sdk/pdc_api"
	"github.com/stretchr/testify/assert"
)

type IntEnum int
type UintEnum uint

type PayloadDataDD struct {
	Name string `json:"name"`
}

type ResponseData struct {
	Data  string   `json:"data"`
	Page  IntEnum  `json:"page"`
	Page2 UintEnum `json:"page2"`
}

func TestCreateSdkCustomTemplate(t *testing.T) {
	t.Run("testing dengan custom template file", func(t *testing.T) {
		sdk := v2_gots_sdk.NewApiSdk(gin.Default())

		save, err := sdk.GenerateSdkFunc("sdk_custom.ts", "custom.template")
		assert.Nil(t, err)

		sdk.Register(&pdc_api.Api{
			Payload:      PayloadDataDD{},
			Method:       http.MethodPost,
			RelativePath: "/users",
		}, func(ctx *gin.Context) {

		})

		datag := sdk.Group("/data")
		datag.Register(&pdc_api.Api{
			Method: http.MethodGet,
		}, func(ctx *gin.Context) {

		})

		usrg := datag.Group("/user")
		usrg.Register(&pdc_api.Api{
			Method:       http.MethodPost,
			RelativePath: "create",
			Response:     ResponseData{},
		}, func(ctx *gin.Context) {})

		sdk.RegisterGroup("/product", func(group *gin.RouterGroup, register v2_gots_sdk.RegisterFunc) {
			register(&pdc_api.Api{
				Payload:      PayloadDataDD{},
				Method:       http.MethodPost,
				RelativePath: "/create",
			})
		})

		sdk.RegisterGroup("/product_data", func(group *gin.RouterGroup, register v2_gots_sdk.RegisterFunc) {
			register(&pdc_api.Api{
				Payload:      []*PayloadDataDD{},
				Method:       http.MethodPost,
				Response:     []string{},
				RelativePath: "/create",
			})
		})

		save()
	})
}

func TestCreateSDK(t *testing.T) {
	sdk := v2_gots_sdk.NewApiSdk(gin.Default())

	save, err := sdk.GenerateSdkFunc("sdk.ts", "")
	assert.Nil(t, err)

	sdk.Register(&pdc_api.Api{
		Payload:      PayloadDataDD{},
		Method:       http.MethodPost,
		RelativePath: "/users",
	}, func(ctx *gin.Context) {

	})

	datag := sdk.Group("/data")
	datag.Register(&pdc_api.Api{
		Method: http.MethodGet,
	}, func(ctx *gin.Context) {

	})

	usrg := datag.Group("/user")
	usrg.Register(&pdc_api.Api{
		Method:       http.MethodPost,
		RelativePath: "create",
		Response:     ResponseData{},
	}, func(ctx *gin.Context) {})

	sdk.RegisterGroup("/product", func(group *gin.RouterGroup, register v2_gots_sdk.RegisterFunc) {
		register(&pdc_api.Api{
			Payload:      PayloadDataDD{},
			Method:       http.MethodPost,
			RelativePath: "/create",
		})
	})

	sdk.RegisterGroup("/product_data", func(group *gin.RouterGroup, register v2_gots_sdk.RegisterFunc) {
		register(&pdc_api.Api{
			Payload:      []*PayloadDataDD{},
			Method:       http.MethodPost,
			Response:     []string{},
			RelativePath: "/create",
		})
	})

	save()
}
