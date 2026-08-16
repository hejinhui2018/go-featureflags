package getopt

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestInt(t *testing.T) {
	assert := assert.New(t)

	p := NewFlagSet("", 0)
	var k int
	p.IntVar(&k, "k", 16, "set k")
	i := p.Int64("i", -1, "set i")
	j := p.Uint("j", 64, "set j")

	err := p.parse([]string{"bin", "-i", "32", "normal arg"})
	assert.Nil(err, "Expected err to be nil")
	assert.Equal(3, p.optindex, "Expected to only parse two arguments")
	assert.Equal(int64(32), *i, "Expected -i argument to equal 32")
	assert.Equal(uint(64), *j, "Expected -j argument to equal 64, since unset")
	assert.Equal(16, k, "Expected -k argument to equal 16, since unset")
}

func TestBool(t *testing.T) {
	assert := assert.New(t)

	p := NewFlagSet("", 0)
	var a bool
	p.BoolVar(&a, "a", false, "set a")
	b := p.Bool("b", false, "set b")

	err := p.parse([]string{"bin", "-a", "normal arg"})
	assert.Nil(err, "Expected err to be nil")
	assert.Equal(2, p.optindex, "Expected to only parse two arguments")
	assert.Equal(true, a, "Expected -a argument to be set")
	assert.Equal(false, *b, "Expected -b argument to not be set")
}

func TestString(t *testing.T) {
	assert := assert.New(t)

	p := NewFlagSet("", 0)
	get := p.String("c", "default", "get -c")

	opt := "some options"
	err := p.parse([]string{"bin", "-c", opt, "normal arg"})
	assert.Nil(err, "Expected err to be nil")
	assert.Equal(opt, *get, "Expected argument to be parsed")
}

func TestFloat64(t *testing.T) {
	assert := assert.New(t)

	p := NewFlagSet("", 0)
	f := p.Float64("f", -3.14, "get -f")

	err := p.parse([]string{"bin", "-f", "3.14", "normal arg"})
	assert.Nil(err, "Expected err to be nil")
	assert.Equal(3.14, *f, "Expected -f to equal 3.14")
}

func TestDuration(t *testing.T) {
	assert := assert.New(t)

	p := NewFlagSet("", 0)
	d := p.Duration("d", 0, "get -d")

	err := p.parse([]string{"bin", "-d", "1h3m", "normal arg"})
	assert.Nil(err, "Expected err to be nil")
	assert.Equal(time.Hour+3*time.Minute, *d, "Expected -d to equal 1 hour and 3 minutes")
}
