package main

import (
	"errors"

	"math"
	// "fmt"
	"hash/fnv"
)

const risizeUpThreshold float32 = 0.60
const resizeDownThreshold float32 = 0.12
const keyNotFoundErrorMsg string = "key not found"

var primes = []uint64{
	17,
	23,
	37,
	53,
	97,
	193,
	389,
	769,
	1543,
	3079,
	6151,
	12289,
	24593,
	49157,
	98317,
	196613,
	393241,
	786433,
	1572869,
	3145739,
	6291469,
	12582917,
	25165843,
	50331653,
	100663319,
	201326611,
	402653189,
	805306457,
	1610612741,
}

const fnvPrime uint64 = 1099511628211
const maxUint64 uint64 = 18446744073709551615

type data struct {
	key           string
	value         any
}

func pickLargestLength(candidate uint64) uint64 {
	for _, v := range primes {
		if v >= candidate {
			return v
		}
	}
	return 0
}

func pickSmallestLength(candidate uint64) uint64 {
	for i := len(primes) - 1; i >= 0; i-- {
		prime := primes[i]
		if prime <= candidate {
			return prime
		}

	}
	return 0
}

func isPrime(candidate uint64) bool {

	limit := uint64(math.Sqrt(float64(candidate)))

	for i := uint64(3); i <= limit; i += 2 {

		// if even return false
		if candidate%i == 0 {
			return false
		}
	}
	return true
}

func computePrimeNumber(candidate uint64) uint64 {
	var start uint64 = 0

	if candidate%2 == 0 {
		start = candidate + uint64(1)
	}

	for i := start; i < maxUint64; i += 2 {

		if isPrime(i) {
			return i
		}
	}

	return uint64(0)

}

func getPrime(candidate uint64, nextSizeUp bool) uint64 {

	var foundPrime uint64 = 0

	if nextSizeUp {
		foundPrime = pickLargestLength(candidate)
	} else {
		foundPrime = pickSmallestLength(candidate)
	}
	if foundPrime != 0 {
		return foundPrime
	}

	// Compute prime number from candidate if not found
	// in pre-computed primes array
	foundPrime = computePrimeNumber(candidate)

	if foundPrime == 0 {
		panic("Prime could not be found!")
	}
	return foundPrime

}

// Custom implementation of the FNV-1a hashing algorithm
func fnvHash(key string) uint64 {
	var hash uint64 = 14695981039346656037
	len := len(key)
	for i := 0; i < len; i++ {
		hash ^= uint64(key[i])
		hash *= fnvPrime

	}

	return hash
}

func fnvHashLib(key string) uint64 {

	hash := fnv.New64a()
	hash.Write([]byte(key))
	return hash.Sum64()
}

type HashTable struct {
	length               uint64
	slots                []data
	activeSlotCounter    uint64
	occupiedSlotCounter  uint64
	debugCollistionCount uint64
}

func New(length uint64) *HashTable {

	primeLength := pickLargestLength(length)
	return &HashTable{
		length:               primeLength,
		slots:                make([]data, primeLength),
		activeSlotCounter:    0,
		occupiedSlotCounter:  0,
		debugCollistionCount: 0,
	}
}

func (h *HashTable) computeNextSizeDown() uint64 {

	candidate := h.length / 2
	return getPrime(candidate, false)
}

func (h *HashTable) computeNextSizeUp() uint64 {
	if h.length*2 >= maxUint64 {
		panic("The hash table cant be resized again because it will overflow uint64!")
	}
	candidate := h.length * 2

	return getPrime(candidate, true)
}

func (h *HashTable) doubleHashing(key string, collisionCount uint64) uint64 {
	hashKey := fnvHash(key)
	hash1 := hashKey % h.length
	hash2 := 1 + (hashKey % (h.length - 1))

	floatCollisionCount := float64(collisionCount)
	tetrahedralFloat := (math.Pow(floatCollisionCount, 3) - floatCollisionCount) / 6
	return ((hash1 + collisionCount*hash2) + uint64(tetrahedralFloat)) % h.length
	//return (hash1 + collisionCount*hash2) % h.length
}

func (h *HashTable) computeLoadFactor() float32 {

	return float32(h.occupiedSlotCounter) / float32(h.length)
}

func (h *HashTable) resize(newSize uint64) {

	h.length = newSize
	h.activeSlotCounter = 0
	h.occupiedSlotCounter = 0
	newSlots := make([]data, newSize)

	for i := range len(h.slots) {
		item := h.slots[i]

		if item.value == nil {
			continue
		}

		h.insert(newSlots, item.key, item.value)
	}
	h.slots = newSlots

}
func (h *HashTable) insertItem(slots []data, index uint64, key string, value any) {
	slots[index] = data{
		key:   key,
		value: value,
	}
	h.activeSlotCounter++
	h.occupiedSlotCounter++
}

func (*HashTable) updateValue(slots []data, index uint64, key string, value any) {
	if slots[index].value != nil && slots[index].key == key {
		slots[index].value = value
	}
}

func (h *HashTable) insert(slots []data, key string, value any) {
	var collisionCount uint64 = 0
	homeLocation := h.doubleHashing(key, collisionCount)

	if slots[homeLocation].value == nil {
		h.insertItem(slots, homeLocation, key, value)
		return
	}

	h.updateValue(slots, homeLocation, key, value)

	// Start Probing
	for {
		collisionCount++
		h.debugCollistionCount++
		deltaLocation := h.doubleHashing(key, collisionCount)

		if deltaLocation == homeLocation {
			break
		}

		h.updateValue(slots, deltaLocation, key, value)

		if slots[deltaLocation].value == nil {
			h.insertItem(slots, deltaLocation, key, value)
			return
		}
	}

}

func (h *HashTable) Insert(key string, value any) {

	loadFactor := h.computeLoadFactor()
	if loadFactor >= risizeUpThreshold {

		newLength := h.computeNextSizeUp()
		h.resize(newLength)
	}
	h.insert(h.slots, key, value)
}

func (h *HashTable) Search(key string) (any, error) {
	var collisionCount uint64 = 0
	homeLocation := h.doubleHashing(key, collisionCount)
	item := h.slots[homeLocation]
	if item.value == nil {
		return nil, errors.New(keyNotFoundErrorMsg)
	}
	if item.key == key && item.value != nil {
		return item.value, nil
	}

	// Probe!
	for {
		collisionCount++
		deltaLocation := h.doubleHashing(key, collisionCount)
		if deltaLocation == homeLocation {
			break
		}
		item := h.slots[deltaLocation]
		if item.value == nil {
			return nil, errors.New(keyNotFoundErrorMsg)
		}

		if item.key == key && item.value !=nil {
			return item.value, nil
		}
	}

	return nil, errors.New(keyNotFoundErrorMsg)
}



func ( h *HashTable) deleteItem( item *data){
	item.value = nil
	h.activeSlotCounter--
	loadFactor := h.computeLoadFactor()
	if loadFactor <= resizeDownThreshold {
		newLength := h.computeNextSizeDown()
		h.resize(newLength)
	}
}

func (h *HashTable) Delete(key string) error {

	// TODO: Test deletion when probing

	var collisionCount uint64 = 0
	homeLocation := h.doubleHashing(key, collisionCount)
	item := &h.slots[homeLocation]

	if item.value == nil {
		return errors.New(keyNotFoundErrorMsg)
	}
	if item.key == key && item.value == nil {
		return errors.New(keyNotFoundErrorMsg)
	}
	if item.key == key && item.value !=nil {

		h.deleteItem(item)
		return nil
	}

	// Probe!

	for {
		collisionCount++
		deltaLocation := h.doubleHashing(key, collisionCount)
		if deltaLocation == homeLocation {
			break
		}
		item := &h.slots[deltaLocation]
		if item.value == nil {
			return errors.New(keyNotFoundErrorMsg)
		}
		if item.key == key && item.value != nil {
			h.deleteItem(item)
			return nil
		}
	}
	return errors.New(keyNotFoundErrorMsg)
}
