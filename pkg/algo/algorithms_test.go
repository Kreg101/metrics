package algo

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Float2String(t *testing.T) {
	tt := []struct {
		name string
		args float64
		want string
	}{
		{
			name: "no trim",
			args: 1.235,
			want: "1.235",
		},
		{
			name: "trim 1 digit",
			args: 1.230,
			want: "1.23",
		},
		{
			name: "trim 2 digits",
			args: 1.200,
			want: "1.2",
		},
		{
			name: "integer",
			args: 1.00000,
			want: "1",
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.want, Float2String(tc.args))
		})
	}
}
