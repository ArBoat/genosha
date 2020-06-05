package testifyExample

import (
  "github.com/stretchr/testify/assert"
  "testing"
)

func TestAdd(t *testing.T) {
  result := Add(1, 2)
  assert.Equal(t, 3, result, "they should be equal")
}

func TestSubtract(t *testing.T) {
  result := Subtract(1, 2)
  assert.Equal(t, -1, result, "they should not be equal")
}

func TestMultiply(t *testing.T) {
  result := Multiply(3, 2)
  assert.Equal(t, 6, result, "they should not be equal")
}

func TestDivision(t *testing.T) {
  result, err := Division(10, 2)
  assert.Nil(t, err, "should be nil")
  assert.Equal(t, 5, result, "they should not be equal")
  _, err = Division(10, 0)
  assert.NotNil(t, err, "should not be nil")
}
