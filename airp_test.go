package airp

import (
	"io/ioutil"
	"reflect"
	"testing"
)

func TestLexer(t *testing.T) {
	tests := []struct {
		have string
		want []token
	}{
		{`{"a": null}`, []token{
			{Type: objectOToken},
			{Type: stringToken, Value: "a"},
			{Type: colonToken},
			{Type: nullToken},
			{Type: objectCToken},
		}},
		{`[false, -31.2, 5, "ab\"cd"]`, []token{
			{Type: arrayOToken},
			{Type: falseToken},
			{Type: commaToken},
			{Type: numberToken, Value: "-31.2"},
			{Type: commaToken},
			{Type: numberToken, Value: "5"},
			{Type: commaToken},
			{Type: stringToken, Value: "ab\\\"cd"},
			{Type: arrayCToken},
		}},
		{`{"a": 20, "b": [true, null]}`, []token{
			{Type: objectOToken},
			{Type: stringToken, Value: "a"},
			{Type: colonToken},
			{Type: numberToken, Value: "20"},
			{Type: commaToken},
			{Type: stringToken, Value: "b"},
			{Type: colonToken},
			{Type: arrayOToken},
			{Type: trueToken},
			{Type: commaToken},
			{Type: nullToken},
			{Type: arrayCToken},
			{Type: objectCToken},
		}},
	}
	for _, test := range tests {
		ch := lex(test.have)
		for _, w := range test.want {
			tk := <-ch
			if tk != w {
				t.Errorf("have %v, got %s, want %s",
					test.have, tk, test.want)
			}
		}
		tk, ok := <-ch
		if ok {
			t.Errorf("expected nothing, got %s", tk)
		}
	}
}

func TestParser(t *testing.T) {
	tests := []struct {
		have string
		want Node
	}{
		{`{"a": null}`, Node{
			jsonType: Object,
			Children: []Node{
				{key: "a", jsonType: Null},
			},
		}},
		{`[false, -31.2, 5, "ab\"cd"]`, Node{
			jsonType: Array,
			Children: []Node{
				{jsonType: Bool, value: "false"},
				{jsonType: Number, value: "-31.2"},
				{jsonType: Number, value: "5"},
				{jsonType: String, value: "ab\\\"cd"},
			},
		}},
		{`{"a": 20, "b": [true, null]}`, Node{
			jsonType: Object,
			Children: []Node{
				{key: "a", jsonType: Number, value: "20"},
				{key: "b", jsonType: Array, Children: []Node{
					{jsonType: Bool, value: "true"},
					{jsonType: Null},
				}},
			},
		}},
	}
	for _, test := range tests {
		ch := lex(test.have)
		ast, err := parse(ch)
		if err != nil {
			t.Error(err)
		}
		if !eqNode(ast, &test.want) {
			t.Errorf("for %v got %v", test.have, ast)
		}
	}
}

func TestFile(t *testing.T) {
	want := &Node{key: "", jsonType: 6, value: "", Children: []Node{
		{key: "bool", jsonType: 2, value: "true"},
		{key: "obj", jsonType: 6, value: "", Children: []Node{
			{key: "v", jsonType: 1, value: ""}}},
		{key: "values", jsonType: 5, value: "", Children: []Node{
			{key: "", jsonType: 6, value: "", Children: []Node{
				{key: "a", jsonType: 3, value: "5"},
				{key: "b", jsonType: 4, value: "hi"},
				{key: "c", jsonType: 3, value: "5.8"},
				{key: "d", jsonType: 1, value: ""},
				{key: "e", jsonType: 2, value: "true"}}},
			{key: "", jsonType: 6, value: "", Children: []Node{
				{key: "a", jsonType: 5, value: "", Children: []Node{
					{key: "", jsonType: 3, value: "5"},
					{key: "", jsonType: 3, value: "6"},
					{key: "", jsonType: 3, value: "7"},
					{key: "", jsonType: 3, value: "8"}}},
				{key: "b", jsonType: 4, value: "hi2"},
				{key: "c", jsonType: 3, value: "5.9"},
				{key: "d", jsonType: 6, value: "", Children: []Node{
					{key: "f", jsonType: 4, value: "Hello there!"}}},
				{key: "e", jsonType: 2, value: "false"}}}}}}}
	data, err := ioutil.ReadFile("../testfiles/test.json")
	if err != nil {
		t.Error(err)
	}
	n, err := parse(lex(string(data)))
	if err != nil {
		t.Error(err)
	}
	if !eqNode(want, n) {
		t.Error("WRONG!")
	}
}

func TestValue(t *testing.T) {
	tests := []struct {
		have string
		want interface{}
	}{
		{`{"a": null}`, map[string]interface{}{"a": nil}},
		{`[false, -31.2, 5, "ab\"cd"]`, []interface{}{
			false, -31.2, float64(5), "ab\\\"cd",
		}},
		{`{"a": 20, "b": [true, null]}`, map[string]interface{}{
			"a": float64(20), "b": []interface{}{true, nil},
		}},
	}
	for _, test := range tests {
		ast, _ := parse(lex(test.have))
		itf, err := ast.Value()
		if err != nil {
			t.Error(err)
		}
		if !reflect.DeepEqual(itf, test.want) {
			t.Errorf("want %v, got %v", test.want, itf)
		}
	}
}

func TestASTStringer(t *testing.T) {
	tests := []struct {
		want string
		have Node
	}{
		{`{"a":null}`, Node{
			jsonType: Object,
			Children: []Node{
				{key: "a", jsonType: Null},
			},
		}},
		{`[false,-31.2,5,"ab\"cd"]`, Node{
			jsonType: Array,
			Children: []Node{
				{jsonType: Bool, value: "false"},
				{jsonType: Number, value: "-31.2"},
				{jsonType: Number, value: "5"},
				{jsonType: String, value: "ab\\\"cd"},
			},
		}},
		{`{"a":20,"b":[true,null]}`, Node{
			jsonType: Object,
			Children: []Node{
				{key: "a", jsonType: Number, value: "20"},
				{key: "b", jsonType: Array, Children: []Node{
					{jsonType: Bool, value: "true"},
					{jsonType: Null},
				}},
			},
		}},
	}
	for _, test := range tests {
		got := test.have.String()
		if got != test.want {
			t.Errorf("want: %s, got: %s", test.want, got)
		}
	}
}