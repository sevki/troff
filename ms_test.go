package troff

import (
	"bytes"
	"testing"
)

func Test_cmdPrint(t *testing.T) {
	type args struct {
		ms   Macro
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
				ms: increaseTypeSize,
			},
			want: []byte(".LG\n"),
		},
		{
			name: "decreaseSize",
			args: args{
				ms: decreaseTypeSize,
			},
			want: []byte(".SM\n"),
		},
		{
			name: "indentParagraph",
			args: args{
				ms: beginAndIndentParagraph,
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
			if got := msPrint(tt.args.ms, tt.args.args...); bytes.Compare(got, tt.want) != 0 {
				t.Errorf("cmdPrintf() = %v, want %v", got, tt.want)
			}
		})
	}
}
