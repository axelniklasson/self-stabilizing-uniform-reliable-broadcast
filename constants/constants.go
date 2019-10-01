package constants

// HostsFilePath tells the system where to find the file with all other hosts in the network
const HostsFilePath = "./hosts.txt"

// TestHostFilePath is used during testing for temp hosts
const TestHostFilePath = "./test_hosts.txt"

// ModuleRunSleepMs is the duration each module sleeps before one iteration of the do forever loop
const ModuleRunSleepMs = 250

// BufferUnitSize is used to control the number of messages allowed to be in the buffer for a processor
const BufferUnitSize = 100

// ThetafdW is the threshold used by the theta fd
const ThetafdW = 100

// ServerBufferSize is the size of the server buffer used when reading messages over the UDP socket
const ServerBufferSize = 1024

// UnitTestingEnvVar indicates that the system is performing unit tests
const UnitTestingEnvVar = "UNIT_TESTING"

// TravisEnvVar indicates that the system is running on travis ci
const TravisEnvVar = "TRAVIS_CI"

// IPEnvVar is set to allow for the IP address to be used when binding API
const IPEnvVar = "IP"

// Env is used to control what env is currently launching the app
const Env = "ENV"

// BufferUnitSizeEnvVar is used to set buffer size at runtime
const BufferUnitSizeEnvVar = "BUFFER_SIZE"

// ModuleRunSleepEnvVar is used to set module run sleep at runtime (ms)
const ModuleRunSleepEnvVar = "MOD_RUN_SLEEP"
