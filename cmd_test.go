package troff

import (
	"bytes"
	"testing"
)

func Test_cmdPrint(t *testing.T) {
	type args struct {
		cmd  Command
		args []interface{}
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{
			name: "increaseSize",
			args: args{
				cmd: increaseTypeSize,
			},
			want: []byte(".LG\n"),
		},
		{
			name: "decreaseSize",
			args: args{
				cmd: decreaseTypeSize,
			},
			want: []byte(".SM\n"),
		},
		{
			name: "indentParagraph",
			args: args{
				cmd: beginAndIndentParagraph,
				args: []interface{}{
					"first:",
					9,
				},
			},
			want: []byte(".IP first: 9\n"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := cmdPrint(tt.args.cmd, tt.args.args...); bytes.Compare(got, tt.want) != 0 {
				t.Errorf("cmdPrintf() = %v, want %v", got, tt.want)
			}
		})
	}
}
