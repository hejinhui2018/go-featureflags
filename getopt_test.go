package getopt

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSimpleCase(t *testing.T) {
	assert := assert.New(t)
	opts, i, err := Getopts([]string{
		"test_bin", "-afo", "output-file", "normal arg"}, "afo:")
	assert.Nil(err, "Expected err to be nil")
	assert.Equal(len(opts), 3, "Expected 3 options to be parsed")
	assert.Equal(i, 3, "Expected non-option args to start at index 2")
	assert.Equal(opts[0], Option{'a', ""})
	assert.Equal(opts[1], Option{'f', ""})
	assert.Equal(opts[2], Option{'o', "output-file"})
}

func TestShortFormArgument(t *testing.T) {
	assert := assert.New(t)
	opts, i, err := Getopts([]string{
		"test_bin", "-afooutput-file", "normal arg"}, "afo:")
	assert.Nil(err, "Expected err to be nil")
	assert.Equal(len(opts), 3, "Expected 3 options to be parsed")
	assert.Equal(i, 2, "Expected non-option args to start at index 2")
	assert.Equal(opts[0], Option{'a', ""})
	assert.Equal(opts[1], Option{'f', ""})
	assert.Equal(opts[2], Option{'o', "output-file"})
}

func TestSeparateArgs(t *testing.T) {
	assert := assert.New(t)
	opts, i, err := Getopts([]string{
		"test_bin", "-a", "-f", "-o", "output-file", "normal arg"}, "afo:")
	assert.Nil(err, "Expected err to be nil")
	assert.Equal(len(opts), 3, "Expected 3 options to be parsed")
	assert.Equal(i, 5, "Expected non-option args to start at index 5")
	assert.Equal(opts[0], Option{'a', ""})
	assert.Equal(opts[1], Option{'f', ""})
	assert.Equal(opts[2], Option{'o', "output-file"})
}

func TestTwoDashes(t *testing.T) {
	assert := assert.New(t)
	opts, i, err := Getopts([]string{
		"test_bin", "-afo", "output-file", "--", "-f", "normal arg"}, "afo:")
	assert.Nil(err, "Expected err to be nil")
	assert.Equal(len(opts), 3, "Expected 3 options to be parsed")
	assert.Equal(i, 4, "Expected non-option args to start at index 4")
	assert.Equal(opts[0], Option{'a', ""})
	assert.Equal(opts[1], Option{'f', ""})
	assert.Equal(opts[2], Option{'o', "output-file"})
}

func TestUnknownOption(t *testing.T) {
	assert := assert.New(t)
	_, _, err := Getopts([]string{"test_bin", "-x"}, "y")
	var errt UnknownOptionError
	assert.IsType(err, errt, "Expected unknown option error")
	assert.Equal(err.Error(), fmt.Sprintf("%s: unknown option -x", os.Args[0]),
		"Expected POSIX-compatible error message")
}

func TestMissingOption(t *testing.T) {
	assert := assert.New(t)
	_, _, err := Getopts([]string{"test_bin", "-x"}, "x:")
	var errt MissingOptionError
	assert.IsType(err, errt, "Expected missing option error")
	assert.Equal(err.Error(), fmt.Sprintf("%s: expected argument for -x",
		os.Args[0]), "Expected POSIX-compatible error message")
}

func TestExpectedMissingOption(t *testing.T) {
	assert := assert.New(t)
	opts, _, err := Getopts([]string{"test_bin", "-x"}, ":x:")
	assert.Nil(err, "Expected err to be nil")
	assert.Equal(len(opts), 1, "Expected 1 option to be parsed")
	assert.Equal(opts[0], Option{':', "x"})
}

func TestNoOption(t *testing.T) {
	assert := assert.New(t)
	_, i, _ := Getopts([]string{"test_bin"}, "")
	assert.Equal(i, 1, "Expected non-option args to start at index 1")
}
