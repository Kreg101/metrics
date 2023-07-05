package handler

import (
	"github.com/Kreg101/metrics/internal/server/constants"
	"github.com/Kreg101/metrics/internal/server/storage"
	"net/http"
	"reflect"
	"testing"
)

func TestMux_ServeHTTP(t *testing.T) {
	type fields struct {
		storage *storage.Storage
	}
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mux := Mux{
				storage: tt.fields.storage,
			}
			mux.ServeHTTP(tt.args.w, tt.args.r)
		})
	}
}

func TestNewMux(t *testing.T) {
	tests := []struct {
		name string
		want *Mux
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewMux(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewMux() = %v, want %v", got, tt.want)
			}
		})
	}
}

//func Test_removeEmptyElements(t *testing.T) {
//	tt := []struct {
//		name     string
//		data     *[]string
//		expected *[]string
//	}{
//		{name: "nothing to do", data: &[]string{"abc", "bc", "d"}, expected: &[]string{"abc", "bc", "d"}},
//	}
//	for _, tc := range tt {
//		t.Run(tc.name, func(t *testing.T) {
//			removeEmptyElements(tt.args.split)
//		})
//	}
//}

func Test_requestValidation(t *testing.T) {
	type args struct {
		split []string
	}
	tests := []struct {
		name  string
		args  args
		want  interface{}
		want1 constants.Error
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := requestValidation(tt.args.split)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("requestValidation() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("requestValidation() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
