package model

import "testing"

func TestFunctionsReturnsThree(t *testing.T) {
	fns := Functions()
	if len(fns) != 3 {
		t.Fatalf("want 3 functions, got %d", len(fns))
	}
}

func TestValidate(t *testing.T) {
	cases := []struct {
		name    string
		typ     TaskType
		params  map[string]any
		wantErr bool
	}{
		{"ok translate", TypeTranslate, map[string]any{"text": "hi", "from": "en", "to": "zh"}, false},
		{"empty text", TypeTranslate, map[string]any{"text": "", "from": "en", "to": "zh"}, true},
		{"bad dir", TypeTranslate, map[string]any{"text": "hi", "from": "en", "to": "fr"}, true},
		{"same dir", TypeTranslate, map[string]any{"text": "hi", "from": "en", "to": "en"}, true},
		{"ok summarize", TypeSummarize, map[string]any{"text": "long"}, false},
		{"summarize empty", TypeSummarize, map[string]any{"text": ""}, true},
		{"unknown type", TaskType("x"), map[string]any{"text": "hi"}, true},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := Validate(c.typ, c.params)
			if (err != nil) != c.wantErr {
				t.Fatalf("err=%v wantErr=%v", err, c.wantErr)
			}
		})
	}
}
