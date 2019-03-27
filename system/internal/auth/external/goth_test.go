package external

import (
	"fmt"
	"testing"

	"github.com/jmoiron/sqlx/types"

	"github.com/crusttech/crust/internal/settings"
	"github.com/crusttech/crust/internal/test"
)

func Test_extractOpenIDProviders(t *testing.T) {
	type args struct {
		s       settings.KV
		handler func(name, url, clientKey, secret string)
	}

	var (
		handlerOutput string
		handler       = func(name, url, clientKey, secret string) {
			handlerOutput += fmt.Sprintf("%s / %s / %s / %s | ", name, url, clientKey, secret)
		}
	)

	tests := []struct {
		name    string
		args    args
		handled string
	}{
		{
			args: args{
				s: settings.KV{
					"auth.providers.openid-connect.foo":        types.JSONText("true"),
					"auth.providers.openid-connect.foo.url":    types.JSONText(`"url"`),
					"auth.providers.openid-connect.foo.key":    types.JSONText(`"key"`),
					"auth.providers.openid-connect.foo.secret": types.JSONText(`"secret"`),
					"auth.providers.openid-connect.bar":        types.JSONText("true"),
					"auth.providers.openid-connect.bar.url":    types.JSONText(`"url"`),
					"auth.providers.openid-connect.bar.key":    types.JSONText(`"key"`),
					"auth.providers.openid-connect.bar.secret": types.JSONText(`"secret"`),
					"auth.providers.openid-connect.baz":        types.JSONText("false"),
					"auth.providers.openid-connect.baz.url":    types.JSONText(`"url"`),
					"auth.providers.openid-connect.baz.key":    types.JSONText(`"key"`),
					"auth.providers.openid-connect.baz.secret": types.JSONText(`"secret"`),
				},
				handler: handler,
			},
			handled: "bar / url / key / secret | foo / url / key / secret | ",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handlerOutput = "" // reset
			extractOpenIDProviders(tt.args.s, tt.args.handler)

			test.Assert(t,
				handlerOutput == tt.handled,
				"Expecting extracted values: %q, got %q",
				tt.handled,
				handlerOutput)
		})
	}
}
