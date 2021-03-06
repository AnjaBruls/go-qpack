package qpack

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func TestPack(t *testing.T) {
	var m = make(map[interface{}]interface{})
	m["Names"] = []string{"Iris", "Sasha"}

	type TestQp struct {
		Name string `qp:"myname"`
	}
	testQp := TestQp{Name: "Iris"}

	var dat interface{}
	decoder := json.NewDecoder(strings.NewReader(`{"num":[1, 1.0, 1.1]}`))
	decoder.UseNumber()
	if err := decoder.Decode(&dat); err != nil {
		panic(err)
	}

	cases := []struct {
		in   interface{}
		want []byte
		err  error
	}{
		{" Hi Qpack", []byte{
			140, 239, 163, 159, 32, 72, 105, 32, 81, 112, 97, 99, 107}, nil},
		{testQp, []byte{244, 134, 109, 121, 110, 97, 109, 101, 132, 73, 114, 105, 115}, nil},
		{dat, []byte{244, 131, 110, 117, 109, 240, 1, 127, 236, 154, 153,
			153, 153, 153, 153, 241, 63}, nil},
		{true, []byte{249}, nil},
		{false, []byte{250}, nil},
		{nil, []byte{251}, nil},
		{-1, []byte{64}, nil},
		{-60, []byte{123}, nil},
		{-61, []byte{232, 195}, nil},
		{0, []byte{0}, nil},
		{1, []byte{1}, nil},
		{int64(4), []byte{4}, nil},
		{63, []byte{63}, nil},
		{64, []byte{232, 64}, nil},
		{-1.0, []byte{125}, nil},
		{0.0, []byte{126}, nil},
		{1.0, []byte{127}, nil},
		{-120, []byte{232, 136}, nil},
		{-0xfe, []byte{233, 2, 255}, nil},
		{-0xfedcba, []byte{234, 70, 35, 1, 255}, nil},
		{-0xfedcba9876, []byte{235, 138, 103, 69, 35, 1, 255, 255, 255}, nil},
		{120, []byte{232, 120}, nil},
		{0xfe, []byte{233, 254, 0}, nil},
		{0xfedcba, []byte{234, 186, 220, 254, 0}, nil},
		{0xfedcba9876, []byte{235, 118, 152, 186, 220, 254, 0, 0, 0}, nil},
		{-1.234567, []byte{236, 135, 136, 155, 83, 201, 192, 243, 191}, nil},
		{123.4567, []byte{236, 83, 5, 163, 146, 58, 221, 94, 64}, nil},
		{[]float64{0.0, 1.1, 2.2}, []byte{
			240, 126, 236, 154, 153, 153, 153, 153, 153, 241, 63,
			236, 154, 153, 153, 153, 153, 153, 1, 64}, nil},
		{[]int{10, 20, 30, 40, 50}, []byte{242, 10, 20, 30, 40, 50}, nil},
		{[]int{10, 20, 30, 40, 50, 60}, []byte{
			252, 10, 20, 30, 40, 50, 60, 254}, nil},
		{[]interface{}{0, m}, []byte{
			239, 0, 244, 133, 78, 97, 109, 101, 115, 239, 132, 73, 114, 105,
			115, 133, 83, 97, 115, 104, 97}, nil},
	}
	for _, c := range cases {
		if c.err == nil {
			got, err := Pack(c.in)
			if err != nil {
				t.Errorf(
					"Pack(%q) returned an unexpexted error: %s", c.in, err)
			} else if !bytes.Equal(got, c.want) {
				t.Errorf("Pack(%v) == %v, want %v", c.in, got, c.want)
			}
		} else {
			_, err := Pack(c.in)
			t.Errorf("Error: %v", err)
		}
	}
}
