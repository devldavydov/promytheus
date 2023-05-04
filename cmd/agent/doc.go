/*
Agent for sending metrics.

Usage:

	agent [flags]

Flags:

	-a metrics server address (env ADDRESS)
	-r report interval (env REPORT_INTERVAL)
	-p poll interval (env POLL_INTERVAL)
	-k hmac sign key (env KEY)
	-l rate limit (env RATE_LIMIT)

Additional environment variables:

	LOG_LEVEL - logging level
	LOG_FILE  - file path to log
*/
package main
