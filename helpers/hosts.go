package helpers

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/axelniklasson/self-stabilizing-uniform-reliable-broadcast/constants"
	"github.com/axelniklasson/self-stabilizing-uniform-reliable-broadcast/models"
)

// Processors is a slice of all processors
var Processors []models.Processor

func getHostsPath() string {
	if IsUnitTesting() {
		return constants.TestHostFilePath
	}
	return constants.HostsFilePath
}

// IPStringToSlice converts a string of the form X.X.X.X to an IP struct []byte
func IPStringToSlice(ipString string) []byte {
	parts := strings.Split(ipString, ".")
	ip := []byte{}
	for _, p := range parts {
		x, _ := strconv.Atoi(p)
		ip = append(ip, byte(x))
	}

	return ip
}

// ParseHostsFile parses a host file at the given path and returns a slice of corresponding processors
func ParseHostsFile() ([]models.Processor, error) {
	// parse file and exit if error
	file, err := os.Open(getHostsPath())
	if err != nil {
		return nil, err
	}

	// close file after we're done
	defer file.Close()

	reader := bufio.NewReader(file)
	var line string
	processors := []models.Processor{}

	// read hosts file line by line
	for {
		line, err = reader.ReadString('\n')

		if err != nil {
			break
		}

		// parse out processor information from line and att to slice
		if len(line) > 0 {
			parts := strings.Split(line, ",")
			// check for malformed line
			if len(parts) != 3 {
				return nil, fmt.Errorf("Malformed line in hosts file: %s", line)
			}
			id, err := strconv.Atoi(parts[0])
			if err != nil {
				return nil, fmt.Errorf("Could not parse ID. Got error: %v", err)
			}
			hostname := parts[1]
			ipString := strings.TrimSuffix(parts[2], "\n")

			// build net.IP struct
			ip := IPStringToSlice(ipString)
			p := models.Processor{ID: id, Hostname: hostname, IPString: ipString, IP: ip}
			processors = append(processors, p)
		}
	}

	// EOF error is expected
	if err != io.EOF {
		return nil, err
	}

	Processors = processors
	return processors, nil
}

