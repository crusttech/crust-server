module github.com/crusttech/crust-server

go 1.12

require (
	github.com/cortezaproject/corteza-server v0.0.0-20200710123055-c12d537f0093
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/joho/godotenv v1.3.0
	go.uber.org/zap v1.13.0
)

replace gopkg.in/Masterminds/squirrel.v1 => github.com/Masterminds/squirrel v1.1.0
