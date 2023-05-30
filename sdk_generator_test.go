package v2_gots_sdk_test

import (
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/pdcgo/v2_gots_sdk"
)

type PayloadDataDD struct {
	Name string
}

func TestCreateSDK(t *testing.T) {
	sdk := v2_gots_sdk.NewApiSdk(gin.Default())

	save := sdk.GenerateSdkFunc("sdk.ts")

	sdk.Register(&v2_gots_sdk.Api{
		Payload:      PayloadDataDD{},
		Method:       http.MethodPost,
		RelativePath: "/users",
	}, func(ctx *gin.Context) {

	})

	datag := sdk.Group("/data")
	datag.Register(&v2_gots_sdk.Api{
		Method: http.MethodGet,
	}, func(ctx *gin.Context) {

	})

	usrg := datag.Group("/user")
	usrg.Register(&v2_gots_sdk.Api{
		Method: http.MethodPost,
	}, func(ctx *gin.Context) {})

	sdk.RegisterGroup("/product", func(group *gin.RouterGroup, register v2_gots_sdk.RegisterFunc) {
		register(&v2_gots_sdk.Api{
			Payload:      PayloadDataDD{},
			Method:       http.MethodPost,
			RelativePath: "/create",
		})
	})

	save()
}