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
