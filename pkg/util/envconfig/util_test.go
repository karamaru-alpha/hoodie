package envconfig

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestEnvConfig_Load(t *testing.T) {
	for name, tt := range map[string]struct {
		in  any
		out any
		env map[string]string
	}{
		"BOOL": {
			in:  &struct{ Bool bool }{},
			out: &struct{ Bool bool }{Bool: true},
			env: map[string]string{"BOOL": "true"},
		},
		"INT32": {
			in:  &struct{ Int32 int32 }{},
			out: &struct{ Int32 int32 }{Int32: 1},
			env: map[string]string{"INT32": "1"},
		},
		"INT64": {
			in:  &struct{ Int64 int64 }{},
			out: &struct{ Int64 int64 }{Int64: 10},
			env: map[string]string{"INT64": "10"},
		},
		"FLOAT64": {
			in:  &struct{ Float64 float64 }{},
			out: &struct{ Float64 float64 }{Float64: 1.1},
			env: map[string]string{"FLOAT64": "1.1"},
		},
		"STRING": {
			in:  &struct{ String string }{},
			out: &struct{ String string }{String: "string value"},
			env: map[string]string{"STRING": "string value"},
		},
		"TIME": {
			in:  &struct{ Time time.Time }{},
			out: &struct{ Time time.Time }{Time: time.Date(2000, 1, 23, 4, 5, 6, 0, time.UTC)},
			env: map[string]string{"TIME": "2000-01-23T04:05:06Z"},
		},
		"DURATION": {
			in:  &struct{ Duration time.Duration }{},
			out: &struct{ Duration time.Duration }{Duration: 3 * time.Hour},
			env: map[string]string{"DURATION": "3h"},
		},
		"BOOL_SLICE": {
			in:  &struct{ BoolSlice []bool }{},
			out: &struct{ BoolSlice []bool }{BoolSlice: []bool{true, false, true}},
			env: map[string]string{"BOOL_SLICE": "true,false,true"},
		},
		"INT32_SLICE": {
			in:  &struct{ Int32Slice []int32 }{},
			out: &struct{ Int32Slice []int32 }{Int32Slice: []int32{1, 2, 3}},
			env: map[string]string{"INT32_SLICE": "1,2,3"},
		},
		"INT64_SLICE": {
			in:  &struct{ Int64Slice []int64 }{},
			out: &struct{ Int64Slice []int64 }{Int64Slice: []int64{10, 20, 30}},
			env: map[string]string{"INT64_SLICE": "10,20,30"},
		},
		"FLOAT64_SLICE": {
			in:  &struct{ Float64Slice []float64 }{},
			out: &struct{ Float64Slice []float64 }{Float64Slice: []float64{1.1, 2.2, 3.3}},
			env: map[string]string{"FLOAT64_SLICE": "1.1,2.2,3.3"},
		},
		"STRING_SLICE": {
			in:  &struct{ StringSlice []string }{},
			out: &struct{ StringSlice []string }{StringSlice: []string{"string value1", "string value2"}},
			env: map[string]string{"STRING_SLICE": "string value1,string value2"},
		},
		"DURATION_SLICE": {
			in:  &struct{ DurationSlice []time.Duration }{},
			out: &struct{ DurationSlice []time.Duration }{DurationSlice: []time.Duration{3 * time.Hour, 2 * time.Minute}},
			env: map[string]string{"DURATION_SLICE": "3h,2m"},
		},
	} {
		t.Run(name, func(t *testing.T) {
			for k, v := range tt.env {
				t.Setenv(k, v)
			}
			err := Load(tt.in)
			assert.NoError(t, err)
			assert.Equal(t, tt.out, tt.in)
		})
	}
}
