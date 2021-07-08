package godb

import (
	"reflect"
	"testing"
)

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
