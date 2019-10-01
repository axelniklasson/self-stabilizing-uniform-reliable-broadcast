package helpers

import (
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/axelniklasson/self-stabilizing-uniform-reliable-broadcast/constants"
)

// SetUnitTestingEnv sets an env var to indicate that unit testing is being performed
func SetUnitTestingEnv() {
	os.Setenv(constants.UnitTestingEnvVar, "true")
}

// IsUnitTesting returns true if the unit testing env var is set
func IsUnitTesting() bool {
	_, isSet := os.LookupEnv(constants.UnitTestingEnvVar)
	return isSet
}

// IsRunningOnTravis returns true if the travis ci env var is set
func IsRunningOnTravis() bool {
	_, isSet := os.LookupEnv(constants.TravisEnvVar)
	return isSet
}

// GetIP returns the current IP address defined by IPEnvVar, if not set it returns ""
func GetIP() string {
	if IsDevEnv() {
		return "127.0.0.1"
	}

	val, isSet := os.LookupEnv(constants.IPEnvVar)
	if !isSet {
		return ""
	}

	return val
}

// IsDevEnv returns true if env var is set to DEV
func IsDevEnv() bool {
	val, isSet := os.LookupEnv(constants.Env)
	if !isSet {
		return false
	}
	return val == "DEV"
}

// GetBufferUnitSize returns the env var buffer size if applicable, otherwise -1
func GetBufferUnitSize() int {
	val, isSet := os.LookupEnv(constants.BufferUnitSizeEnvVar)
	if !isSet {
		return constants.BufferUnitSize
	}

	x, err := strconv.Atoi(val)
	if err != nil {
		log.Fatal(err)
	}
	return x
}

// GetModRunSleepDuration retuns the Duration to sleep for main module
func GetModRunSleepDuration() time.Duration {
	val, isSet := os.LookupEnv(constants.ModuleRunSleepEnvVar)
	if !isSet {
		return time.Duration(rand.Int31n(int32(constants.ModuleRunSleepMs))) * time.Millisecond
	}

	x, err := strconv.Atoi(val)
	if err != nil {
		log.Fatal(err)
	}

	return time.Duration(rand.Int31n(int32(x))) * time.Millisecond
}
