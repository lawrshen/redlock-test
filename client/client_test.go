package client

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNormal(t *testing.T) {
	reset()
	expectResult := int(finalStockValue)
	result := client(normal)
	assert.Equal(t, expectResult, result)
}

func TestCrash(t *testing.T) {
	reset()
	expectResult := int(finalStockValue)
	result := client(crash)
	assert.Equal(t, expectResult, result)
}

func TestClockJump(t *testing.T) {
	reset()
	expectResult := int(1)
	result := client(clockjump)
	assert.Equal(t, expectResult, result)
}