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

type nodeKey struct {
	value string
	hash  uint64
}

func NewKey(key string) nodeKey {

	if len(key) > 36 {
		panic("A Key can't be longer than 36 characters!")
	}

	return nodeKey{
		value: key,
		hash:  fnvHash(key),
	}
}

type data[V any] struct {
	key      nodeKey
	value    V
	occupied bool
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

type HashTable[V any] struct {
	length               uint64
	slots                []data[V]
	activeSlotCounter    uint64
	occupiedSlotCounter  uint64
	debugCollistionCount uint64
}

func New[V any](length uint64) *HashTable[V] {

	primeLength := pickLargestLength(length)
	return &HashTable[V]{
		length:               primeLength,
		slots:                make([]data[V], primeLength),
		activeSlotCounter:    0,
		occupiedSlotCounter:  0,
		debugCollistionCount: 0,
	}
}

func (h *HashTable[V]) computeNextSizeDown() uint64 {

	candidate := h.length / 2
	return getPrime(candidate, false)
}

func (h *HashTable[V]) computeNextSizeUp() uint64 {
	if h.length*2 >= maxUint64 {
		panic("The hash table cant be resized again because it will overflow uint64!")
	}
	candidate := h.length * 2

	return getPrime(candidate, true)
}

func (h *HashTable[V]) doubleHashing(key nodeKey, collisionCount uint64) uint64 {
	hashKey := key.hash
	hash1 := hashKey % h.length
	hash2 := 1 + (hashKey % (h.length - 1))

	floatCollisionCount := float64(collisionCount)
	tetrahedralFloat := (math.Pow(floatCollisionCount, 3) - floatCollisionCount) / 6
	return ((hash1 + collisionCount*hash2) + uint64(tetrahedralFloat)) % h.length
	//return (hash1 + collisionCount*hash2) % h.length
}

func (h *HashTable[V]) computeLoadFactor() float32 {

	return float32(h.occupiedSlotCounter) / float32(h.length)
}

func (h *HashTable[V]) resize(newSize uint64) {

	h.length = newSize
	h.activeSlotCounter = 0
	h.occupiedSlotCounter = 0
	newSlots := make([]data[V], newSize)

	for i := range len(h.slots) {
		item := h.slots[i]

		if !item.occupied {
			continue
		}

		h.insert(newSlots, item.key, item.value)
	}
	h.slots = newSlots

}
func (h *HashTable[V]) insertItem(slots []data[V], index uint64, key nodeKey, value V) {
	slots[index] = data[V]{
		key:      key,
		value:    value,
		occupied: true,
	}
	h.activeSlotCounter++
	h.occupiedSlotCounter++
}

func (*HashTable[V]) updateValue(slots []data[V], index uint64, key nodeKey, value V) bool {
	if slots[index].occupied && slots[index].key.value == key.value {
		slots[index].value = value
		return true
	}
	return false
}

func (h *HashTable[V]) insert(slots []data[V], key nodeKey, value V) {

	var collisionCount uint64 = 0
	homeLocation := h.doubleHashing(key, collisionCount)

	if !slots[homeLocation].occupied {
		h.insertItem(slots, homeLocation, key, value)
		return
	}

	if h.updateValue(slots, homeLocation, key, value) {
		return
	}

	// Start Probing
	for {
		collisionCount++
		h.debugCollistionCount++
		deltaLocation := h.doubleHashing(key, collisionCount)

		if deltaLocation == homeLocation {
			break
		}

		if h.updateValue(slots, deltaLocation, key, value) {
			return
		}

		if !slots[deltaLocation].occupied {
			h.insertItem(slots, deltaLocation, key, value)
			return
		}
	}

}

func (h *HashTable[V]) Insert(key string, value V) {

	loadFactor := h.computeLoadFactor()
	if loadFactor >= risizeUpThreshold {

		newLength := h.computeNextSizeUp()
		h.resize(newLength)
	}
	k := NewKey(key)

	h.insert(h.slots, k, value)
}

func (h *HashTable[V]) Search(key string) (V, error) {
	var collisionCount uint64 = 0
	var zero V

	k := NewKey(key)

	homeLocation := h.doubleHashing(k, collisionCount)
	item := h.slots[homeLocation]
	if !item.occupied {
		return zero, errors.New(keyNotFoundErrorMsg)
	}
	if item.key.value == key {
		return item.value, nil
	}

	// Probe!
	for {
		collisionCount++
		deltaLocation := h.doubleHashing(k, collisionCount)
		if deltaLocation == homeLocation {
			break
		}
		item := h.slots[deltaLocation]
		if !item.occupied {
			return zero, errors.New(keyNotFoundErrorMsg)
		}

		if item.key.value == key {
			return item.value, nil
		}
	}

	return zero, errors.New(keyNotFoundErrorMsg)
}

func (h *HashTable[V]) deleteItem(item *data[V]) {
	var zero V
	item.value = zero
	item.occupied = false
	h.activeSlotCounter--
	h.occupiedSlotCounter--
	loadFactor := h.computeLoadFactor()
	if loadFactor <= resizeDownThreshold {
		newLength := h.computeNextSizeDown()
		h.resize(newLength)
	}
}

func (h *HashTable[V]) Delete(key string) error {

	// TODO: Test deletion when probing

	k := NewKey(key)

	var collisionCount uint64 = 0
	homeLocation := h.doubleHashing(k, collisionCount)
	item := &h.slots[homeLocation]

	if !item.occupied {
		return errors.New(keyNotFoundErrorMsg)
	}
	if item.key.value == key {

		h.deleteItem(item)
		return nil
	}

	// Probe!

	for {
		collisionCount++
		deltaLocation := h.doubleHashing(k, collisionCount)
		if deltaLocation == homeLocation {
			break
		}
		item := &h.slots[deltaLocation]
		if !item.occupied {
			return errors.New(keyNotFoundErrorMsg)
		}
		if item.key.value == key {
			h.deleteItem(item)
			return nil
		}
	}
	return errors.New(keyNotFoundErrorMsg)
}
