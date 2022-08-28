package sql2code

import "testing"

func TestGetGormStr(t *testing.T) {
	type args struct {
		args *Args
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "sql form param",
			args: args{args: &Args{
				Sql: `
create table user
(
    id         bigint unsigned auto_increment
        primary key,
    created_at datetime        null,
    updated_at datetime        null,
    deleted_at datetime        null,
    name       char(50)        not null comment '用户名',
    password   char(100)       not null comment '密码',
    email      char(50)        not null comment '邮件',
    phone      bigint unsigned not null comment '手机号码',
    age        tinyint         not null comment '年龄',
    gender     tinyint         not null comment '性别，1:男，2:女，3:未知',
    constraint user_email_uindex
        unique (email)
);
`,
			}},
			wantErr: false,
		},
		{
			name: "sql from file",
			args: args{args: &Args{
				InputFile: "test.sql",
			}},
			want:    "",
			wantErr: false,
		},
		{
			name: "sql from db",
			args: args{args: &Args{
				DBDsn:   "root:123456@(192.168.3.37:3306)/test",
				DBTable: "user",
			}},
			want:    "",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetGormCode(tt.args.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetGormCode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Log(got)
		})
	}
}
