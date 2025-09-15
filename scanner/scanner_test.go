package scanner

import "testing"

func TestScannerTokenType(t *testing.T) {
	testMap := map[string][]TokenType{
		"5 5 5":    {TOK_INTEGER, TOK_INTEGER, TOK_INTEGER},
		"if if if": {TOK_IF, TOK_IF, TOK_IF},
		`
		if x < 5
			x = x + 1;
		end if;
		`: {
			TOK_IF,
			TOK_IDENT,
			TOK_LT,
			TOK_INTEGER,
			TOK_IDENT,
			TOK_EQ,
			TOK_IDENT,
			TOK_PLUS,
			TOK_INTEGER,
			TOK_SEMI,
			TOK_END,
			TOK_IF,
			TOK_SEMI,
		},
		`1.5 "Hello, World!" while`: {TOK_FLOAT, TOK_STRING, TOK_WHILE},
	}

	tokStringsLen := len(TokStrings)

	for s, ttypeList := range testMap {
		scan := NewScanner()

		scan.Tokenize(s)

		for idx, ttype := range ttypeList {
			tok := scan.Advance()

			expected := ""

			if int(ttype) < tokStringsLen {
				expected = TokStrings[ttype]
			}

			found := tok.Text

			if int(tok.TType) < tokStringsLen {
				found = TokStrings[tok.TType]
			}

			if ttype != tok.TType {
				t.Errorf(
					"For string \"%s\" at index %d expected %s (%d) but got %s (%d)",
					s,
					idx,
					expected,
					ttype,
					found,
					tok.TType,
				)
			}
		}
	}
}

func TestScanFloat(t *testing.T) {
	validFloats := [...]string{
		"1.5",
		"0.5",
		".5E10",
		"5.E10",
		"1.5E+1",
		"1.5E-1",
	}
	invalidFloats := [...]string{
		"1",
		"hello",
		"",
		".",
		".E5",
		"1.5E",
	}

	for _, s := range validFloats {
		_, foundS, found := scanFloat(s)

		if !found {
			t.Errorf("Expected to find float in \"%s\"", s)
		}
		if s != foundS {
			t.Errorf(
				"Expected to receive \"%s\", got \"%s\"",
				s,
				foundS,
			)
		}
	}
	for _, s := range invalidFloats {
		_, _, found := scanFloat(s)

		if found {
			t.Errorf("Expected not to find float in \"%s\"", s)
		}
	}
}

func TestScanString(t *testing.T) {
	validStrings := [...]string{
		`"Hello"`,
		`"Hello\"\"\"\\"`,
		`"Hello, World!!!!\n\n\t\r"`,
	}
	invalidStrings := [...]string{
		`"Hello`,
		`"Hello\"`,
		`Hello`,
		`9.5`,
		`9`,
		`if x < 5`,
	}

	for _, s := range validStrings {
		_, foundS, found := scanString(s)

		if !found {
			t.Errorf("Expected to find string literal in \"%s\"", s)
		}
		if s != foundS {
			t.Errorf(
				"Expected to receive \"%s\", got \"%s\"",
				s,
				foundS,
			)
		}
	}
	for _, s := range invalidStrings {
		_, _, found := scanString(s)

		if found {
			t.Errorf("Expected not to find string literal in \"%s\"", s)
		}
	}
}
