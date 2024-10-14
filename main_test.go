package main

import (
	"testing"
)

func resetEnv(t *testing.T) {
	t.Setenv("environment", "")
	t.Setenv("port", "")
	t.Setenv("dbName", "")
	t.Setenv("MongoVaultSideCar", "")
	t.Setenv("logLevel", "")
	t.Setenv("printDBQueries", "")
}
