package config

import "time"

var (
	WRITETIMEOUT               time.Duration = time.Second * 10
	READTIMEOUT                time.Duration = time.Second * 10
	JWT_ACCESS_EXPIRATION_TIME time.Duration = time.Minute * 15
)
