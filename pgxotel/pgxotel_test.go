package pgxotel

import "testing"

func TestSpanOperationName(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		stmt string
		want string
	}{
		{name: "select", stmt: "select * from cards", want: "SELECT"},
		{name: "line comment", stmt: "-- name: ListCards :many\nselect * from cards", want: "SELECT"},
		{name: "block comment", stmt: "/* sqlc */\ninsert into cards values ($1)", want: "INSERT"},
		{name: "only comment", stmt: "-- name only", want: "UNKNOWN"},
		{name: "empty", stmt: " ", want: "UNKNOWN"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := SpanOperationName(tt.stmt); got != tt.want {
				t.Fatalf("SpanOperationName(%q) = %q, want %q", tt.stmt, got, tt.want)
			}
		})
	}
}

func TestTrimLeadingSQLComments(t *testing.T) {
	t.Parallel()

	stmt := "/* one */\n-- two\nupdate cards set name = $1"
	if got := TrimLeadingSQLComments(stmt); got != "update cards set name = $1" {
		t.Fatalf("TrimLeadingSQLComments() = %q", got)
	}
}
