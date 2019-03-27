package settings

import (
	"testing"

	"github.com/jmoiron/sqlx/types"
)

func TestKV_Bool(t *testing.T) {
	type args struct {
		k string
	}
	tests := []struct {
		name  string
		kv    KV
		args  args
		wantV bool
	}{
		{
			name:  "True value should return true",
			kv:    KV{"true-value": types.JSONText(`true`)},
			args:  args{k: "true-value"},
			wantV: true,
		},
		{
			name:  "Null value should return false",
			kv:    KV{"null-value": types.JSONText(`null`)},
			args:  args{k: "null-value"},
			wantV: false,
		},
		{
			name:  "Unexisting value should return false",
			kv:    KV{},
			args:  args{k: "unexisting"},
			wantV: false,
		},
		{
			name:  "Invalid KV should return false",
			kv:    nil,
			args:  args{k: "invalid-kv"},
			wantV: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotV := tt.kv.Bool(tt.args.k); gotV != tt.wantV {
				t.Errorf("KV.Bool() = %v, want %v", gotV, tt.wantV)
			}
		})
	}
}
