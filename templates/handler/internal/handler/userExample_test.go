package handler

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/zhufuyi/pkg/mysql/query"
)

// 打印序列化请求和返回参数的json数据
func TestCovert2json(t *testing.T) {
	type args struct {
		obj interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "CreateUserExampleRequest",
			args:    args{&CreateUserExampleRequest{}},
			wantErr: false,
		},
		{
			name:    "UpdateUserExampleByIDRequest",
			args:    args{&UpdateUserExampleByIDRequest{}},
			wantErr: false,
		},
		{
			name: "GetUserExamplesRequest",
			args: args{&GetUserExamplesRequest{
				Params: query.Params{
					Page: 0,
					Size: 0,
					Sort: "a",
					Columns: []query.Column{
						{},
					},
				},
			},
			},
			wantErr: false,
		},
		{
			name:    "GetUserExampleByIDRespond",
			args:    args{GetUserExampleByIDRespond{}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out, err := json.MarshalIndent(tt.args.obj, "", "    ")
			if err != nil {
				t.Fatal(err)
			}

			fmt.Println(string(out))
		})
	}
}
