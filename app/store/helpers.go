package store

import (
	"errors"
	"strings"

	"modernc.org/sqlite"

	"github.com/nilBora/servers-manager/app/enum"
)

// isUniqueViolation checks if error is a unique constraint violation
func isUniqueViolation(err error) bool {
	if err == nil {
		return false
	}

	// sqlite: SQLITE_CONSTRAINT_UNIQUE = 2067, SQLITE_CONSTRAINT_PRIMARYKEY = 1555
	var sqliteErr *sqlite.Error
	if errors.As(err, &sqliteErr) {
		code := sqliteErr.Code()
		return code == 2067 || code == 1555
	}

	// fallback: check error message
	return strings.Contains(err.Error(), "UNIQUE constraint")
}

func parseProviderType(s string) (enum.ProviderType, error) {
	return enum.ParseProviderType(s)
}

func parseAccountType(s string) (enum.AccountType, error) {
	return enum.ParseAccountType(s)
}

func parseServerStatus(s string) (enum.ServerStatus, error) {
	return enum.ParseServerStatus(s)
}

func parseServerType(s string) (enum.ServerType, error) {
	return enum.ParseServerType(s)
}

func parseLogAction(s string) (enum.LogAction, error) {
	return enum.ParseLogAction(s)
}
