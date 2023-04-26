/*
Metrics server.

Usage:

	cmd/server [flags]

Flags:

	-a server address (env ADDRESS)
	-i store interval (env STORE_INTERVAL)
	-f store file (env STORE_FILE)
	-r should restore (env RESTORE)
	-k hmac sign key (env KEY)
	-d database dsn (env DATABASE_DSN)

Additional environment variables:

	LOG_LEVEL - logging level
	LOG_FILE  - file path to log
*/
package main
