package helpers

import (
	"os"
	"reflect"
	"testing"

	"gotest.tools/assert"

	"github.com/axelniklasson/self-stabilizing-uniform-reliable-broadcast/constants"
	"github.com/axelniklasson/self-stabilizing-uniform-reliable-broadcast/models"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func createHostsFile(hostsString string) {
	f, err := os.Create(constants.TestHostFilePath)
	check(err)
	defer f.Close()

	f.WriteString(hostsString)
}

func TestCanParseValidFile(t *testing.T) {
	createHostsFile("0,localhost,127.0.0.1\n1,localhost,127.0.0.1\n")
	SetUnitTestingEnv()

	// parse file and make sure that processors are parsed correctly
	processors, err := ParseHostsFile()
	assert.NilError(t, err)
	assert.Equal(t, len(processors), 2)
	assert.Assert(t, reflect.DeepEqual(processors[0], models.Processor{ID: 0, Hostname: "localhost", IPString: "127.0.0.1", IP: []byte{127, 0, 0, 1}}))
	assert.Assert(t, reflect.DeepEqual(processors[1], models.Processor{ID: 1, Hostname: "localhost", IPString: "127.0.0.1", IP: []byte{127, 0, 0, 1}}))

	os.Remove(constants.TestHostFilePath)
}

func TestFailsOnMalformedLine(t *testing.T) {
	createHostsFile("0,localhost,127.0.0.1,5\n1,localhost,127.0.0.1\n")
	SetUnitTestingEnv()

	// make sure it fails appropriately due to malformed hosts file
	processors, err := ParseHostsFile()
	assert.Assert(t, processors == nil)
	assert.Error(t, err, "Malformed line in hosts file: 0,localhost,127.0.0.1,5\n")

	os.Remove(constants.TestHostFilePath)
}

func TestFailsOnIdNotBeingAnInt(t *testing.T) {
	createHostsFile("asd,localhost,127.0.0.1\n1,localhost,127.0.0.1\n")
	SetUnitTestingEnv()

	// make sure it fails appropriately due to malformed hosts file
	processors, err := ParseHostsFile()
	assert.Assert(t, processors == nil)
	assert.Error(t, err, "Could not parse ID. Got error: strconv.Atoi: parsing \"asd\": invalid syntax")

	os.Remove(constants.TestHostFilePath)
}
