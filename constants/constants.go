package constants

// HostsFilePath tells the system where to find the file with all other hosts in the network
const HostsFilePath = "./hosts.txt"

// ModuleRunSleepSeconds is the amount of seconds each module sleeps before one iteration of the do forever loop
const ModuleRunSleepSeconds = 1

// ThetafdW is the amount of messages considered
const ThetafdW = 10

// ServerBufferSize is the size of the server buffer used when reading messages over the UDP socket
const ServerBufferSize = 1024
