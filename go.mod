module github.com/crusttech/crust-server

go 1.12

require (
	github.com/cortezaproject/corteza-server v0.0.0-20191025131936-b9b164663371
	github.com/crusttech/permit v0.0.0-20190226221958-6c0c4bece8da
	github.com/davecgh/go-spew v1.1.1
	github.com/joho/godotenv v1.3.0
	github.com/pkg/errors v0.8.1
	github.com/spf13/cobra v0.0.3
	github.com/spf13/viper v1.4.0 // indirect
	go.uber.org/zap v1.10.0
	golang.org/x/mobile v0.0.0-20190312151609-d3739f865fa6
)

replace gopkg.in/Masterminds/squirrel.v1 => github.com/Masterminds/squirrel v1.1.0

replace github.com/cortezaproject/corteza-server => ../corteza-server
