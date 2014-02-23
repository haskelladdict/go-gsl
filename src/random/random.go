// Copyright 2014 Markus Dittrich. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// random wraps gsl random number generation routines
package random

// #cgo CFLAGS: -std=c99 -O2
// #cgo pkg-config: gsl
// #include <gsl/gsl_rng.h>
// #include "random_wrap.h"
import "C"

import (
  "fmt"
  "runtime"
  "unsafe"
)

// RngState stores the random number generator state
type RngState struct {
  state *C.gsl_rng
}

// RngType stores the type of rng method used
type RngType struct {
  rng *C.gsl_rng_type
}

// StatePointer encapsulates a raw pointer to the underlying
// rng state within gsl
type StatePointer unsafe.Pointer

// list of defined random number generators. See gsl documentation
// for more detailed info on each of these.
var (
  Mt19937   = RngType{C.gsl_rng_mt19937}
  Ranlxs0   = RngType{C.gsl_rng_ranlxs0}
  Ranlxs1   = RngType{C.gsl_rng_ranlxs1}
  Ranlxs2   = RngType{C.gsl_rng_ranlxs2}
  Ranlxd1   = RngType{C.gsl_rng_ranlxd1}
  Ranlxd2   = RngType{C.gsl_rng_ranlxd2}
  Ranlux    = RngType{C.gsl_rng_ranlux}
  Ranlux389 = RngType{C.gsl_rng_ranlux389}
  Cmrg      = RngType{C.gsl_rng_cmrg}
  Mrg       = RngType{C.gsl_rng_mrg}
  Taus      = RngType{C.gsl_rng_taus}
  Taus2     = RngType{C.gsl_rng_taus2}
  Gfsr4     = RngType{C.gsl_rng_gfsr4}

  // Unix type rngs - these are not high quality so beware
  Rand          = RngType{C.gsl_rng_rand}
  RandomBSD     = RngType{C.gsl_rng_random_bsd}
  RandomBSD_8   = RngType{C.gsl_rng_random8_bsd}
  RandomBSD_32  = RngType{C.gsl_rng_random32_bsd}
  RandomBSD_64  = RngType{C.gsl_rng_random64_bsd}
  RandomBSD_28  = RngType{C.gsl_rng_random128_bsd}
  RandomBSD_256 = RngType{C.gsl_rng_random256_bsd}

  RandomLibc5     = RngType{C.gsl_rng_random_libc5}
  RandomLibc5_8   = RngType{C.gsl_rng_random8_libc5}
  RandomLibc5_32  = RngType{C.gsl_rng_random32_libc5}
  RandomLibc5_64  = RngType{C.gsl_rng_random64_libc5}
  RandomLibc5_128 = RngType{C.gsl_rng_random128_libc5}
  RandomLibc5_256 = RngType{C.gsl_rng_random256_libc5}

  RandomGlibc2     = RngType{C.gsl_rng_random_glibc2}
  RandomGlibc2_8   = RngType{C.gsl_rng_random8_glibc2}
  RandomGlibc2_32  = RngType{C.gsl_rng_random32_glibc2}
  RandomGlibc2_64  = RngType{C.gsl_rng_random64_glibc2}
  RandomGlibc2_128 = RngType{C.gsl_rng_random128_glibc2}
  RandomGlibc2_256 = RngType{C.gsl_rng_random256_glibc2}

  Rand48 = RngType{C.gsl_rng_rand48}

  // compatibility rng types - typically low quality
  Ranf         = RngType{C.gsl_rng_ranf}
  Ranmar       = RngType{C.gsl_rng_ranmar}
  R250         = RngType{C.gsl_rng_r250}
  Tt800        = RngType{C.gsl_rng_tt800}
  Vax          = RngType{C.gsl_rng_vax}
  Transputer   = RngType{C.gsl_rng_transputer}
  Randu        = RngType{C.gsl_rng_randu}
  Minstd       = RngType{C.gsl_rng_minstd}
  Uni          = RngType{C.gsl_rng_uni}
  Uni32        = RngType{C.gsl_rng_uni32}
  Slatec       = RngType{C.gsl_rng_slatec}
  Zuf          = RngType{C.gsl_rng_zuf}
  Knuthran2    = RngType{C.gsl_rng_knuthran2}
  Knuthran2002 = RngType{C.gsl_rng_knuthran2002}
  Knuthran     = RngType{C.gsl_rng_knuthran}
  Borosh13     = RngType{C.gsl_rng_borosh13}
  Fishman18    = RngType{C.gsl_rng_fishman18}
  Fishman20    = RngType{C.gsl_rng_fishman20}
  Lecuyer21    = RngType{C.gsl_rng_lecuyer21}
  Waterman14   = RngType{C.gsl_rng_waterman14}
  Fishman2x    = RngType{C.gsl_rng_fishman2x}
  Coveyou      = RngType{C.gsl_rng_coveyou}
)

// RNG initialization

// Default returns the default random number generator
func Default() RngType {
  return RngType{C.gsl_rng_default}
}

// DefaultSeed returns the default rng seed
func DefaultSeed() uint64 {
  return uint64(C.gsl_rng_default_seed)
}

// Alloc creates a new random number generator and returs
// it as a RngState object.
func Alloc(rngType RngType) RngState {
  state := RngState{C.gsl_rng_alloc(rngType.rng)}

  // make sure we get rid of any memory associated with the
  // rng within gsl
  runtime.SetFinalizer(&state,
    func(rng *RngState) { C.gsl_rng_free(rng.state) })
  return state
}

// Set initializes (or ‘seeds’) the random number generator. If the
// generator is seeded with the same value of seed on two different runs,
// the same stream of random numbers will be generated by successive calls
// to the routines below. If different values of seed ≥ 1 are supplied,
// then the generated streams of random numbers should be completely
// different. If the seed seed is zero then the standard seed from the
// original implementation is used instead. For example, the original
// Fortran source code for the ranlux generator used a seed of 314159265,
// and so choosing seed equal to zero reproduces this when using
// gsl_rng_ranlux.
//
// When using multiple seeds with the same generator, choose seed values
// greater than zero to avoid collisions with the default setting.
// Note that the most generators only accept 32-bit seeds, with higher
// values being reduced modulo 2^32 . For generators with smaller ranges
// the maximum seed value will typically be lower.
func (s RngState) Set(seed uint64) {
  C.gsl_rng_set(s.state, C.ulong(seed))
}

// EnvSetup reads the environment variables GSL_RNG_TYPE and GSL_RNG_SEED
// and uses their values to set the corresponding library variables
// gsl_rng_default and gsl_rng_default_seed returned by Default() and
// DefaultSeed().
func EnvSetup() RngType {
  return RngType{C.gsl_rng_env_setup()}
}

// RNG sampling functions

// Get returns a random integer from the generator s. The minimum and
// maximum values depend on the algorithm used, but all integers in the
// range [min,max] are equally likely. The values of min and max can be
// determined using the auxiliary functions Max and Min.
func (s RngState) Get() uint64 {
  return uint64(C.gsl_rng_get(s.state))
}

// GetSlice is a convenience function returning a slice of random
// uint64 each between min and max of the selected random number
// generator.
func (s RngState) GetSlice(length int) []uint64 {
  slice := make([]uint64, length)
  for i := 0; i < length; i++ {
    slice[i] = s.Get()
  }
  return slice
}

// Uniform returns a double precision floating point number
// uniformly distributed in the range [0,1). The range includes 0.0
// but excludes 1.0.
func (s RngState) Uniform() float64 {
  return float64(C.gsl_rng_uniform(s.state))
}

// UnformSlice is a convenience function returning a slice of length N
// of uniform random floats in [0,1).
func (s RngState) UniformSlice(length int) []float64 {
  slice := make([]float64, length)
  for i := 0; i < length; i++ {
    slice[i] = s.Uniform()
  }
  return slice
}

// UniformPos function returns a positive double precision floating point
// number uniformly distributed in the range (0,1), excluding both 0.0 and
// 1.0. The number is obtained by sampling the generator with the algorithm
// of Uniform until a non-zero value is obtained. You can use this function
// if you need to avoid a singularity at 0.0.
func (s RngState) UniformPos() float64 {
  return float64(C.gsl_rng_uniform_pos(s.state))
}

// UniformInt returns a random integer from 0 to n − 1 inclusive by scaling
// down and/or discarding samples from the generator r. All integers in the
// range [0, n − 1] are produced with equal probability. For generators with
// a non-zero minimum value an offset is applied so that zero is returned
// with the correct probability. Note that this function is designed for
// sampling from ranges smaller than the range of the underlying generator.
// The parameter n must be less than or equal to the range of the generator r.// If n is larger than the range of the generator then the function
// calls the error handler with an error code of GSL_EINVAL and returns zero.
// In particular, this function is not intended for generating the full range
// of unsigned integer values [0, 2 32 − 1]. Instead choose a generator with
// the maximal integer range and zero minimum value, such as gsl_rng_ranlxd1,
// gsl_rng_mt19937 or gsl_rng_taus, and sample it directly using gsl_rng_get.
// The range of each can be found with the help of auxiliary sections.
func (s RngState) UniformInt(limit uint64) uint64 {
  return uint64(C.gsl_rng_uniform_int(s.state, C.ulong(limit)))
}

// UnformIntSlice is a convenience function returning a slice of length N
// of uniform random integers in [0, n - 1].
func (s RngState) UniformIntSlice(limit uint64, length int) []uint64 {
  slice := make([]uint64, length)
  for i := 0; i < length; i++ {
    slice[i] = s.UniformInt(limit)
  }
  return slice
}

// RNG auxiliary functions

// Name returns the name of the random number generator or
// a rng type
func (s RngState) Name() string {
  return C.GoString(C.gsl_rng_name(s.state))
}

func (t RngType) Name() string {
  return C.GoString(t.rng.name)
}

// String provides a printable string representation for
// an RngState and type
func (s RngState) String() string {
  return s.Name()
}

func (t RngType) String() string {
  return C.GoString(t.rng.name)
}

// Max returns the largest value that the rng underlying RngState
// can handle
func (s RngState) Max() uint64 {
  return uint64(C.gsl_rng_max(s.state))
}

// Min returns the largest value that the rng underlying RngState
// can handle
func (s RngState) Min() uint64 {
  return uint64(C.gsl_rng_min(s.state))
}

// State returns a pointer to the underlying rng state from gsl
func (s RngState) State() StatePointer {
  return StatePointer(C.gsl_rng_state(s.state))
}

// Size returns the size of the rng state.
func (s RngState) Size() uint64 {
  return uint64(C.gsl_rng_size(s.state))
}

// TypesSetup returns a map with available rng type names as keys and
// RngType as values
func TypesSetup() map[string]RngType {
  length := uint64(C.rng_types_length())
  var theCArray *(*C.gsl_rng_type) = C.gsl_rng_types_setup()
  type_slice := (*[1 << 30](*C.gsl_rng_type))(unsafe.Pointer(theCArray))[:length:length]

  rng_type_list := make(map[string]RngType)
  for _, v := range type_slice {
    name := C.GoString(v.name)
    rng_type_list[name] = RngType{v}
  }
  return rng_type_list
}

// Copying, cloning, writing and reading rng state

// Memcpy copies the random number generator src into the pre-existing
// generator dest, making dest into an exact copy of src. The two generators
// must be of the same type.
// NOTE: currently this ignores the return type of gsl_rng_memcpy since
// I don't know what it does (the manual is quiet on that)
func (s RngState) Memcpy(dest RngState) {
  C.gsl_rng_memcpy(dest.state, s.state)
}

// Clone returns a newly created generator which is an exact copy
// of the generator r.
func (s RngState) Clone() RngState {
  return RngState{C.gsl_rng_clone(s.state)}
}

// Fwrite writes the random number state of the random number generator s
// to the given file in binary format. Data is written in the
// native binary format and may not be portable between different
// architectures. Returns an error if there was a problem writing.
func (s RngState) Fwrite(s_filename string) error {
  filename := C.CString(s_filename)
  defer C.free(unsafe.Pointer(filename))

  status := int(C.rng_fwrite(filename, s.state))
  if status != 0 {
    return fmt.Errorf("Failed to write rng state to file.")
  }
  return nil
}

// Fread reads the random number state into the random number generator s
// from the given file name in binary format. The random number generator s
// must be preinitialized with the correct random number generator type
// since type information is not saved. The data is assumed to have been
// written in the native binary format on the same architecture. Returns
// an error if reading fails.
func (s RngState) Fread(s_filename string) (RngState, error) {
  filename := C.CString(s_filename)
  defer C.free(unsafe.Pointer(filename))

  status := int(C.rng_fread(filename, s.state))
  if status != 0 {
    return s, fmt.Errorf("Failed to read rng state from file.")
  }
  return s, nil
}
