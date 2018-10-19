package main

import (
	"encoding/binary"
	"fmt"
	"os"
	"time"
)

const (
	size = 1000000
	min  = 0
	max  = 9999999
)

// randNano returns an array of 'not-so-random' integers obtained from the nanosecond portion
// of a monotonic clock reading.  The values are not particularly random in that they are always
// increasing until they roll over.
func randNano() [size]int {
	var rans [size]int

	for i := 0; i < size; {
		s := time.Now()
		nano := s.Nanosecond()
		if nano < min || nano > max {
			continue
		}
		rans[i] = nano
		i++
	}
	return rans
}

// randNanoHash returns an array of 'somewhat-random' integers by applying a simple hashing function
// to the nanosecond portion of a monotonic clock reading.  Suffers from the same problem as randNano.
func randNanoHash() [size]int {
	var rans [size]int

	for i := 0; i < size; {
		s := time.Now()
		nano := s.Nanosecond()

		// Hash.
		for i := 0; i < 12; i++ {
			nano = ((nano>>8)^nano)*0x6b + i
		}

		if nano < min || nano > max {
			continue
		}
		rans[i] = nano
		i++
	}
	return rans
}

func randTsNanoHash() [size]uint {
	var rans [size]uint
	ts := [size]time.Time{}

	for i := range ts {
		ts[i] = time.Now()
	}

	j, k := 0, size-1
	for i := 0; i < size; i += 2 {
		nano1 := uint(ts[j].Nanosecond())
		nano2 := uint(ts[k].Nanosecond())
		for l := 0; l < 12; l++ {
			nano1 = ((nano1>>8)^nano1)*0x6b + uint(l)
			nano2 = ((nano2>>8)^nano2)*0x6b + uint(l)
		}
		nano1 >>= 41
		nano2 >>= 41
		rans[i] = nano1
		// Don't overflow rans.
		if i+1 < size {
			rans[i+1] = nano2
		}
		j++
		k--
	}

	return rans
}

// randDev returns an array of random uint64 integers obtained from /dev/random.
func randDev() [size]uint64 {
	var urans [size]uint64

	f, err := os.Open("/dev/random")
	if err != nil {
		return urans
	}
	defer f.Close()

	b := make([]byte, 8)
	for i := 0; i < size; {
		_, _ = f.Read(b)
		val, _ := binary.Uvarint(b)
		if val < min || val > max {
			continue
		}
		urans[i] = val
		i++
	}
	return urans
}

func main() {
	fmt.Println("randomness from nanoseconds")
	start := time.Now()
	rans := randNano()
	stop := time.Now()
	//fmt.Println(rans)
	fmt.Printf("took %v to obtain %v values\n", stop.Sub(start), len(rans))

	fmt.Println("randomness from hashed nanoseconds")
	start = time.Now()
	rans = randNano()
	stop = time.Now()
	//fmt.Println(rans)
	fmt.Printf("took %v to obtain %v values\n", stop.Sub(start), len(rans))

	fmt.Println("randomness from /dev/random")
	start = time.Now()
	u64rans := randDev()
	stop = time.Now()
	//fmt.Println(u64rans)
	fmt.Printf("took %v to obtain %v values\n", stop.Sub(start), len(u64rans))
	for _, val := range u64rans {
		fmt.Println(val)
	}

	fmt.Println("randomness from randomized timestamp + hash")
	start = time.Now()
	urans := randTsNanoHash()
	stop = time.Now()
	//fmt.Println(urans)
	fmt.Printf("took %v to obtain %v values\n", stop.Sub(start), len(urans))
}
