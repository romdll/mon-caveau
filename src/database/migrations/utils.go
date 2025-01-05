package migrations

import (
	"database/sql"
	"strings"
	"time"
)

func parseSQLDateTime(raw []byte) string {
	if len(raw) == 0 {
		return ""
	}
	parsed, err := time.Parse("2006-01-02 15:04:05", string(raw))
	if err != nil {
		logger.Warnf("Failed to parse datetime: %v", err)
		return string(raw)
	}
	return parsed.Format("2006-01-02 15:04:05")
}

func nullStringToString(ns sql.NullString) string {
	if ns.Valid {
		return ns.String
	}
	return ""
}

func getCharacterSet(collation string) string {
	if parts := strings.Split(collation, "_"); len(parts) > 0 {
		return parts[0]
	}
	return ""
}
