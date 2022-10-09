package Config

import "time"

type DbTimeoutConfig struct {
	// Time in seconds to wait
	Timeout time.Duration
}

type DbConfig struct {
	Url       string
	Username  string
	Password  string
	Database  string
	Namespace string
	// If true, this will automatically call "sign in" for the database
	AutoLogin bool
	// If true, this will automatically call "use" for the database and namespace
	AutoUse bool
	// This will configure context automatically and setup timeouts
	// Set this to nil to disable auto timeout configuration via ctx
	Timeouts *DbTimeoutConfig
}
