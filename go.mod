module github.com/crusttech/crust-server

go 1.12

require (
	github.com/cortezaproject/corteza-server v0.0.0-20191025131936-b9b164663371
	github.com/crusttech/permit v0.0.0-20190226221958-6c0c4bece8da
	github.com/pkg/errors v0.8.1
	go.uber.org/zap v1.10.0
)

replace gopkg.in/Masterminds/squirrel.v1 => github.com/Masterminds/squirrel v1.1.0
