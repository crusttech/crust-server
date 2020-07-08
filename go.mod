module github.com/crusttech/crust-server

go 1.12

require (
	github.com/DestinyWang/cronexpr v0.0.0-20140423231348-a557574d6c02 // indirect
	github.com/cortezaproject/corteza-server v0.0.0-20200630185829-f2b227954713
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/joho/godotenv v1.3.0
	github.com/kr/pretty v0.1.0 // indirect
	github.com/prometheus/client_golang v0.9.3 // indirect
	github.com/spf13/cobra v0.0.3
	go.uber.org/zap v1.13.0
)

replace gopkg.in/Masterminds/squirrel.v1 => github.com/Masterminds/squirrel v1.1.0
