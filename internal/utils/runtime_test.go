package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetCallerFunctionName(t *testing.T) {
	type args struct {
		skip int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "no skip",
			args: args{
				skip: 0,
			},
			want: "func1",
		}, {
			name: "skip",
			args: args{
				skip: 1,
			},
			want: "tRunner",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetCallerFunctionName(tt.args.skip)
			assert.Equal(t, tt.want, got, "GetCallerFunctionName() = %v, want %v", got, tt.want)
		})
	}
}

func TestGetCallerBaseFunctionName(t *testing.T) {
	type args struct {
		skip int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "no skip",
			args: args{
				skip: 0,
			},
			want: "utils.TestGetCallerBaseFunctionName.func1",
		}, {
			name: "skip",
			args: args{
				skip: 1,
			},
			want: "testing.tRunner",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetCallerBaseFunctionName(tt.args.skip)
			assert.Equal(t, tt.want, got, "GetCallerBaseFunctionName() = %v, want %v", got, tt.want)
		})
	}
}