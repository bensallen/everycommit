package main

import (
	"reflect"
	"testing"
)

func Test_parseURL(t *testing.T) {
	type args struct {
		rawurl string
	}
	tests := []struct {
		name    string
		args    args
		want    *pullRequest
		wantErr bool
	}{
		{
			name:    "Empty URL",
			args:    args{rawurl: ""},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Example Github PR URL",
			args:    args{rawurl: "https://github.com/bensallen/everycommit/pull/1"},
			want:    &pullRequest{owner: "bensallen", project: "everycommit", id: 1},
			wantErr: false,
		},
		{
			name:    "Bad ID",
			args:    args{rawurl: "https://github.com/bensallen/everycommit/pull/1a"},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Bad URL format",
			args:    args{rawurl: "https://github.com/pull/bensallen/everycommit/1"},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Not enough path elements",
			args:    args{rawurl: "https://github.com/bensallen/everycommit/pull"},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Junk URL",
			args:    args{rawurl: "asfdhjfdsglhasfarlkjhsakklbasdf"},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseURL(tt.args.rawurl)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseURL() = %v, want %v", got, tt.want)
			}
		})
	}
}
