module github.com/crusttech/crust-server

go 1.12

require (
	github.com/cortezaproject/corteza-server v0.0.0-20191212081655-25793a923008
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/joho/godotenv v1.3.0
	github.com/kr/pretty v0.1.0 // indirect
	github.com/prometheus/client_golang v0.9.3 // indirect
	github.com/spf13/cobra v0.0.3
	go.uber.org/zap v1.10.0
	gopkg.in/check.v1 v1.0.0-20180628173108-788fd7840127 // indirect
)

replace gopkg.in/Masterminds/squirrel.v1 => github.com/Masterminds/squirrel v1.1.0
