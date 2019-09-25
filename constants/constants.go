package constants

import "time"

// HostsFilePath tells the system where to find the file with all other hosts in the network
const HostsFilePath = "./hosts.txt"

// TestHostFilePath is used during testing for temp hosts
const TestHostFilePath = "./test_hosts.txt"

// ModuleRunSleepDuration is the duration each module sleeps before one iteration of the do forever loop
const ModuleRunSleepDuration = 250 * time.Millisecond

// ThetafdW is the threshold used by the theta fd
const ThetafdW = 10

// ServerBufferSize is the size of the server buffer used when reading messages over the UDP socket
const ServerBufferSize = 1024

// UnitTestingEnvVar indicates that the system is performing unit tests
const UnitTestingEnvVar = "UNIT_TESTING"

// TravisEnvVar indicates that the system is running on travis ci
const TravisEnvVar = "TRAVIS_CI"

// BufferUnitSize is used to control the number of messages allowed to be in the buffer for a processor
const BufferUnitSize = 100
