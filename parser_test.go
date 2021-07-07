package godb

import (
	"reflect"
	"testing"
)

func TestExecuteStatement(t *testing.T) {
	type args struct {
		s *Statement
		t *Table
	}
	tests := []struct {
		name string
		args args
		want ExecuteResult
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ExecuteStatement(tt.args.s, tt.args.t); got != tt.want {
				t.Errorf("ExecuteStatement() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParse(t *testing.T) {
	type args struct {
		input string
	}
	tests := []struct {
		name string
		args args
		want Status
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Parse(tt.args.input); got != tt.want {
				t.Errorf("Parse() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPrepareStatement(t *testing.T) {
	type args struct {
		input string
	}
	tests := []struct {
		name string
		args args
		want *Statement
	}{
		{
			name: "basic",
			args: args{
				input: "insert 1 elo elo",
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := PrepareStatement(tt.args.input); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PrepareStatement() = %v, want %v", got, tt.want)
			}
		})
	}
}
