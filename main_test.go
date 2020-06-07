package main

import (
	"net"
	"testing"
)

func Test_procesInput(t *testing.T) {
	type args struct {
		input string
		scope string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{name: "Ip in scope", args: args{
			input: "192.168.2.55",
			scope: "192.168.2.0/24",
		}, want: true},
		{name: "Ip out of scope", args: args{
			input: "192.168.2.55",
			scope: "192.168.2.22/32",
		}, want: false},
		{name: "Host in scope", args: args{
			input: "one.one.one.one",
			scope: "1.1.1.1/32",
		}, want: true},
		{name: "Host out of scope", args: args{
			input: "one.one.one.one",
			scope: "1.0.0.2/32",
		}, want: false},
		{name: "Url in scope", args: args{
			input: "https://one.one.one.one",
			scope: "1.1.1.1/32",
		}, want: true},
		{name: "Url out of scope", args: args{
			input: "https://one.one.one.one",
			scope: "8.8.8.8/32",
		}, want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, subnet, _ := net.ParseCIDR(tt.args.scope)
			currentScopes := []*net.IPNet{subnet}
			if got := procesInput(tt.args.input, currentScopes); got != tt.want {
				t.Errorf("Assertion (have: %t should be: %t) failed on scope(%s) contains (%s) with test name:\"%s\"", got, tt.want, tt.args.scope, tt.args.input, tt.name)
			}
		})
	}
}
