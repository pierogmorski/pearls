// A vector package.

package vector

import (
	"errors"
	"fmt"
	"math/bits"
)

// A bit vector represented as a slice of positive integer values.
type Vector []uint

// New will take as input the number of bits needed in the vector, and will return a zero value slice
// of uint's to house the bits.
func New(numBits uint) *Vector {
	sliceSize := numBits / bits.UintSize
	if numBits%bits.UintSize != 0 {
		sliceSize += 1
	}

	vec := make(Vector, sliceSize)
	return &vec
}

// MaxBits will return the maximum number of bits representable in the vector.
func (v *Vector) MaxBits() uint {
	return uint(len(*v) * bits.UintSize)
}

// SetBit will take a bit number, and will try to set that bit within the vector.  SetBit will
// return a boolean and an error indicating if the bit was successfully set, if the bit was already
// set (an error condition) or if the bit falls outside the size of the vector (an error condition).
func (v *Vector) SetBit(bit uint) (bool, error) {
	isSet, err := v.TestBit(bit)
	if err != nil {
		return false, errors.New(fmt.Sprintf("Failed to set %v'th bit: %v", bit, err))
	}

	if isSet {
		return isSet, errors.New(fmt.Sprintf("%v'th bit already set", bit))
	}

	vecIndex := bit / bits.UintSize
	bitToSet := bit - (vecIndex * bits.UintSize)
	(*v)[vecIndex] |= 1 << bitToSet
	return true, nil
}

// TestBit will take a bit number, and will test if that bit is set within the vector.  TestBit will
// return a boolean and an error indicating whether or not the bit was set, or if the bit falls
// outside the size of the vector (an error condition).
func (v *Vector) TestBit(bit uint) (bool, error) {
	vecIndex := bit / bits.UintSize
	if vecIndex > uint(len(*v))-1 {
		return false, errors.New(fmt.Sprintf("bit %v falls outside vector range", bit))
	}

	bitToCheck := bit - (vecIndex * bits.UintSize)
	bitVal := (*v)[vecIndex] & (1 << bitToCheck)
	return bitVal > 0, nil
}
