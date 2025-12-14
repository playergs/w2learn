package def

// Health Status Def
const (
	HealthStatusNoCheck    = -1
	HealthStatusCheckError = 0
	HealthStatusCheckOK    = 1
)

// Health Status Flag
const (
	HealthStatusRequestFlagNoCheck  = 0b00
	HealthStatusRequestFlagDatabase = 0b01
	HealthStatusRequestFlagAllCheck = 0b11
)
