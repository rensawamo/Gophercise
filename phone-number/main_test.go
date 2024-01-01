package main

import "testing"

type normalizeTestCase struct {
	input string
	want  string
}

func TestNormalize(t *testing.T) {
	testCases := []struct {
		input string
		want  string
	}{
		{"1234567890", "1234567890"},
		{"123 456 7891", "1234567891"},
		{"(123) 456 7892", "1232567892"},
		{"(123) 456-7893", "1234567893"},
		{"123-456-7894", "1234567894"},
		{"(123)456-7892", "1234567892"},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			actual := normalize(tc.input)
			if actual != tc.want { // ここがテスト構文  で 下のt.Errof が呼び出されない場合に okが帰るテスト構文
				t.Errorf("got %s; want %s", actual, tc.want)
			}
		})
	}
}
