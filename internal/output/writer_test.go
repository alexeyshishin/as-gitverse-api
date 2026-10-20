package output

import (
	"bytes"
	"strings"
	"testing"
)

func TestTableWriter_WriteTable(t *testing.T) {
	var buf bytes.Buffer
	w := New(&buf, false)

	if err := w.WriteTable([]string{"NAME", "VALUE"}, [][]string{{"a", "1"}, {"b", "2"}}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "NAME") || !strings.Contains(out, "a") || !strings.Contains(out, "2") {
		t.Fatalf("got %q", out)
	}
}

func TestJSONWriter_WriteTable(t *testing.T) {
	var buf bytes.Buffer
	w := New(&buf, true)

	if err := w.WriteTable([]string{"NAME"}, [][]string{{"a"}}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), `"NAME":"a"`) {
		t.Fatalf("got %q", buf.String())
	}
}

func TestJSONWriter_WriteBulkResults(t *testing.T) {
	var buf bytes.Buffer
	w := New(&buf, true)

	err := w.WriteBulkResults([]ItemResult{
		{Name: "repo1", OK: true},
		{Name: "repo2", OK: false, Error: "boom"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, `"name":"repo1"`) || !strings.Contains(out, `"error":"boom"`) {
		t.Fatalf("got %q", out)
	}
}

func TestTableWriter_WriteBulkResults(t *testing.T) {
	var buf bytes.Buffer
	w := New(&buf, false)

	err := w.WriteBulkResults([]ItemResult{
		{Name: "repo1", OK: true},
		{Name: "repo2", OK: false, Error: "boom"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "repo1") || !strings.Contains(out, "error") || !strings.Contains(out, "boom") {
		t.Fatalf("got %q", out)
	}
}
