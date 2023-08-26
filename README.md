# v2_gots_sdk

Untuk standarisasi wiring pdc dengan bahasa golang. Sdk ini untuk mengkonversi dari golang ke typescript langsung.

#### Wiring Support

- api
- websocket



#### cara wiring untuk api

Sementara dependency yang dipakai adalah [gin framework](https://gin-gonic.com/)

- create generator api sdk

  `````
  sdk := v2_gots_sdk.NewApiSdk(gin.Default())
  save := sdk.GenerateSdkFunc("sdk.ts", true)
  
  
  // registering api
  
  type PayloadDataDD struct {
  	Name string
  }
  
  sdk.Register(&v2_gots_sdk.Api{
      Payload:      PayloadDataDD{},
      Method:       http.MethodPost,
      RelativePath: "/users",
  }, func(ctx *gin.Context) {
  	// handle api nya
  })
  `````

  

- handle group dari url

  ``````
  member := sdk.Group("/member")
  
  member.Register(&v2_gots_sdk.Api{
      Method: http.MethodGet,
  }, func(ctx *gin.Context) {
  	// handle api nya
  })
  
  member.Register(&v2_gots_sdk.Api{
      Method:       http.MethodPost,
      RelativePath: "create",
      Response:     ResponseData{},
  }, func(ctx *gin.Context) {})
  
  ``````

  

- generate typescript sdk dari api yang sudah di register

  `````
  save := sdk.GenerateSdkFunc("sdk.ts", true) 
  
  save() // save generated sdk.ts
  `````

  

#### cara wiring untuk websocket

library utama untuk websocket pakai [nhooyr.io/websocket](https://github.com/nhooyr/websocket)