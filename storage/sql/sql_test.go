package sql

import "testing"

func TestTranslate(t *testing.T) {
	tests := []struct {
		testCase string
		flavor   flavor
		query    string
		exp      string
	}{
		{
			"sqlite3 query bind replacement",
			flavorMySQL,
			`select foo from bar where foo.zam = $1;`,
			`select foo from bar where foo.zam = ?;`,
		},
		{
			"sqlite3 query bind replacement at newline",
			flavorMySQL,
			`select foo from bar where foo.zam = $1`,
			`select foo from bar where foo.zam = ?`,
		},
		{
			"sqlite3 query true",
			flavorMySQL,
			`select foo from bar where foo.zam = true`,
			`select foo from bar where foo.zam = true`,
		},
		{
			"sqlite3 query false",
			flavorMySQL,
			`select foo from bar where foo.zam = false`,
			`select foo from bar where foo.zam = false`,
		},
		{
			"sqlite3 bytea",
			flavorMySQL,
			`"connector_data" bytea not null,`,
			`"connector_data" blob not null,`,
		},
		{
			"sqlite3 now",
			flavorMySQL,
			`now(),`,
			`now(),`,
		},
	}

	for _, tc := range tests {
		if got := tc.flavor.translate(tc.query); got != tc.exp {
			t.Errorf("%s: want=%q, got=%q", tc.testCase, tc.exp, got)
		}
	}
}
