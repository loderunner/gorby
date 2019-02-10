package main

import (
	"github.com/loderunner/popt"
)

var options = []popt.Option{
	{
		Name:    "address",
		Default: ":8080",
		Usage:   "address to listen for HTTP requests",
		Flag:    "address",
		Short:   "a",
		Env:     "GORBY_ADDRESS",
	},
	{
		Name:    "https-address",
		Default: "",
		Usage:   "address to listen for HTTPS requests (same as HTTP if left blank)",
		Flag:    "https-address",
		Short:   "s",
		Env:     "GORBY_HTTPS_ADDRESS",
	},
	{
		Name:    "api-address",
		Default: ":8081",
		Usage:   "address to listen for gorby API requests",
		Flag:    "api-address",
		Env:     "GORBY_API_ADDRESS",
	},
	{
		Name:    "db",
		Default: "",
		Usage:   "location of the recorded traffic database file (in-memory database, if empty)",
		Flag:    "db",
		Short:   "d",
		Env:     "GORBY_DATABASE_LOCATION",
	},
	{
		Name:    "verbose",
		Default: false,
		Usage:   "more output",
		Flag:    "verbose",
		Short:   "v",
	},
	{
		Name:    "quiet",
		Default: false,
		Usage:   "silence output",
		Flag:    "quiet",
		Short:   "q",
	},
}
