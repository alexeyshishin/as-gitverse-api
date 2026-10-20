package output

import (
	"encoding/json"
	"fmt"
	"io"
	"text/tabwriter"
)

type ItemResult struct {
	Name  string `json:"name"`
	OK    bool   `json:"ok"`
	Error string `json:"error,omitempty"`
}

type Writer interface {
	WriteTable(headers []string, rows [][]string) error
	WriteBulkResults(results []ItemResult) error
}

func New(w io.Writer, jsonMode bool) Writer {
	if jsonMode {
		return &jsonWriter{w: w}
	}
	return &tableWriter{w: w}
}

type tableWriter struct{ w io.Writer }

func (t *tableWriter) WriteTable(headers []string, rows [][]string) error {
	tw := tabwriter.NewWriter(t.w, 0, 4, 2, ' ', 0)
	fmt.Fprintln(tw, joinTab(headers))
	for _, row := range rows {
		fmt.Fprintln(tw, joinTab(row))
	}
	return tw.Flush()
}

func (t *tableWriter) WriteBulkResults(results []ItemResult) error {
	rows := make([][]string, len(results))
	for i, r := range results {
		status := "ok"
		errMsg := ""
		if !r.OK {
			status = "error"
			errMsg = r.Error
		}
		rows[i] = []string{r.Name, status, errMsg}
	}
	return t.WriteTable([]string{"NAME", "STATUS", "ERROR"}, rows)
}

func joinTab(fields []string) string {
	out := ""
	for i, f := range fields {
		if i > 0 {
			out += "\t"
		}
		out += f
	}
	return out
}

type jsonWriter struct{ w io.Writer }

func (j *jsonWriter) WriteTable(headers []string, rows [][]string) error {
	records := make([]map[string]string, len(rows))
	for i, row := range rows {
		rec := make(map[string]string, len(headers))
		for c, h := range headers {
			if c < len(row) {
				rec[h] = row[c]
			}
		}
		records[i] = rec
	}
	return json.NewEncoder(j.w).Encode(records)
}

func (j *jsonWriter) WriteBulkResults(results []ItemResult) error {
	return json.NewEncoder(j.w).Encode(results)
}
