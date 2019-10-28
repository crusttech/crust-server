module github.com/crusttech/crust-server

go 1.12

require (
	github.com/cortezaproject/corteza-server v0.0.0-20191028200743-9f090f355a67
	github.com/crusttech/permit v0.0.0-20190226221958-6c0c4bece8da // indirect
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/joho/godotenv v1.3.0
	github.com/spf13/cobra v0.0.3
	github.com/spf13/pflag v1.0.3
	github.com/spf13/viper v1.4.0 // indirect
	go.uber.org/zap v1.10.0
)

replace gopkg.in/Masterminds/squirrel.v1 => github.com/Masterminds/squirrel v1.1.0

replace github.com/cortezaproject/corteza-server => ../corteza-server
