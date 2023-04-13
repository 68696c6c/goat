package query

import "net/url"

type testCase[I, E any] struct {
	input    I
	expected E
}

type filterResult struct {
	sql    string
	params []any
}

func mustMakeQuery(input string) url.Values {
	result, err := url.ParseQuery(input)
	if err != nil {
		panic(err)
	}
	return result
}
