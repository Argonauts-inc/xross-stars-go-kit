package pgxotel

import "strings"

func SpanOperationName(stmt string) string {
	stmt = TrimLeadingSQLComments(stmt)
	fields := strings.Fields(stmt)
	if len(fields) == 0 {
		return "UNKNOWN"
	}
	return strings.ToUpper(fields[0])
}

func TrimLeadingSQLComments(stmt string) string {
	stmt = strings.TrimSpace(stmt)
	for {
		switch {
		case strings.HasPrefix(stmt, "--"):
			end := strings.IndexByte(stmt, '\n')
			if end < 0 {
				return ""
			}
			stmt = strings.TrimSpace(stmt[end+1:])
		case strings.HasPrefix(stmt, "/*"):
			end := strings.Index(stmt, "*/")
			if end < 0 {
				return ""
			}
			stmt = strings.TrimSpace(stmt[end+2:])
		default:
			return stmt
		}
	}
}
