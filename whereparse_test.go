package whereparse

import (
	"fmt"
	"testing"

	sq "github.com/Masterminds/squirrel"
)

func TestParse(t *testing.T) {
	want := []string{
		"SELECT * WHERE (Alice.Name ~ ? OR Bob.LastName !~ ?)",
		"SELECT * WHERE (Field1 = ? AND (Field2 <> ? OR Field3 > ?))",
		"SELECT * WHERE Bar.Alpha = ?",
		"SELECT * WHERE (Foo.Bar.Beta > ? AND Alpha.Bar <> ?)",
		"SELECT * WHERE (Alice.IsActive AND Bob.LastHash = ?)",
		"No query",
		"The query contains a LIMIT expression",
	}
	input := []string{
		"Alice.Name ~ 'A.*`' OR Bob.LastName !~ 'Bill.*`'",
		"Field1 = \"foo\" AND Field2 != 7 OR Field3 > 11.7",
		"Bar.Alpha = 7",
		"Foo.Bar.Beta > 21 AND Alpha.Bar != 'hello'",
		"SELECT * WHERE Alice.IsActive AND Bob.LastHash = 'ab5534b'",
		"",
		"SELECT * WHERE Alice.IsActive AND Bob.LastHash = 'ab5534b' LIMIT 3",
	}

	for i, data := range input {
		builder := sq.Select("*")
		newbuilder, err := Parse(data, builder)
		if err != nil {
			if err.Error() == want[i] {
				continue
			}
			t.Errorf("Error at %d", i)
			t.Errorf("%s\n", err.Error())
			continue
		}
		str, _, _ := newbuilder.ToSql()
		if str != want[i] {
			fmt.Println(str)
			t.Errorf("Wrong value %d", i)
		}
	}
}
