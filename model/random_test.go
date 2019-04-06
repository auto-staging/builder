package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetRandomValueForBranch(t *testing.T) {
	branch := "testBranch"
	expectedRandom := "4a0905a814"

	assert.Equal(t, expectedRandom, getRandomValueForBranch(branch))
}
