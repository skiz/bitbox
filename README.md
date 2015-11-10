# BitBox

BitBox is a dynamically sized bit array container for Go.

## Overview

* Dynamically resizes your bit container based on usage
* Optimized for speed and memory consumption
* Built in support for AND, OR, XOR bitwise operations
* Supports preallocating memory to avoid reallocations

## Uses

* Flexible Flags
* Reference Tables
* Dirty Cache Markers
* Component Entity Systems
* [Many More...](https://en.wikipedia.org/wiki/Bit_array)

## Install

```
go get -u github.com/skiz/bitbox
```

## Documentation

[![GoDoc](https://godoc.org/github.com/skiz/bitbox?status.svg)](https://godoc.org/github.com/skiz/bitbox)

`go doc` format documentation for this project can be viewed online without
installing the package by using the GoDoc page at:
http://godoc.org/github.com/skiz/bitbox

## How It Works

BitBox keeps track of a byte slice and dynamically reallocates the
underlying array based on which bits you set.

When you Set() a bit, the byte slice is checked to see if the
location of the byte is already available, and if not resizes
the byte array to support setting the bit location.

All bit manipulations are based on their 0 indexed position within
the complete byte array (bit 0 is the first bit).

When querying a single bit, the result is a boolean value.

## Example

Here is a simple sample of the bit based operations and the
underlying byte slice bit values.

```
		b := &bitbox.Bitbox{}
		b.Set(7)            // [0000 0001]
		b.Toggle(3)         // [0001 0001]
		b.Unset(7)          // [0001 0000]
		b.Set(12)           // [0001 0000] [0000 1000]
		b.And([]int{3, 12}) // true
		b.Xor([]int{3, 5})  // true
		b.Or([]int{3, 99})  // true
		b.Toggle(3)         // [0000 0000] [0000 1000]
		b.Resize(20)        // [0001 0000] [0000 1000] [0000 0000]
		b.Resize(7)         // [0001 0000]
		b.Size()            // 8
		b.Get(3)            // true
		b.GetByte(0)        // 0x10
		b.Clear()           // [0000 0000]
		b.Get(4500)         // false
```

# Tests & Benchmarks
This library is fully tested and currently in use by other applications.  There
are several benchmarks available which showcases the raw speed of this library.
All benchmarks are done against 10k+ bits

```
BenchmarkToggle-4           	200000000	         6.32 ns/op
BenchmarkGet-4              	500000000	         3.91 ns/op
BenchmarkSet-4              	200000000	         6.61 ns/op
BenchmarkUnset-4            	500000000	         3.89 ns/op
BenchmarkClear-4            	50000000	         38.7 ns/op
BenchmarkTwoAnd-4           	100000000	         13.0 ns/op
BenchmarkThreeAnd-4         	100000000	         17.7 ns/op
BenchmarkFourAndWorstCase-4 	100000000	         22.4 ns/op
BenchmarkFourAndBestCase-4  	200000000	         8.80 ns/op
BenchmarkTwoOrFirst-4       	200000000	         8.28 ns/op
BenchmarkTwoOrSecond-4      	100000000	         12.7 ns/op
BenchmarkThreeOrNone-4      	100000000	         17.9 ns/op
BenchmarkXorTwo-4           	100000000	         12.3 ns/op
BenchmarkXorTwoNone-4       	100000000	         13.1 ns/op
BenchmarkXorThreeWorstCase-4	100000000	         17.3 ns/op
BenchmarkXorThreeBestCase-4 	200000000	         8.48 ns/op

```
# Future Improvements

Track minimum bit position to remove lower byte requirements when all set bits
are in a high range.  For example if your only bit is 10,000,000 we should only
need to allocate a single byte, calculating the offset when computing the value.

An extention of this thought is to map chunks, but it may be a performance tradeoff.

# License

Released under the [BSD License](https://github.com/skiz/bitbox/blob/master/LICENSE).
