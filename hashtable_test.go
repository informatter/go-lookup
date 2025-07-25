package main

import (
	"fmt"
	"github.com/beevik/guid"
	"testing"
)

func TestNodeKeyMaxCharactersExceeded(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("expected panic, but none occurred")
		}
	}()

	NewKey("747447474788323824328947329847328974329874328974328974329874238974")
}

func TestCreateHashTable(t *testing.T) {
	targetLength := 10
	actualLength := 17
	hashTable := New(uint64(targetLength))
	if hashTable.length != uint64(actualLength) {
		t.Errorf("HashTable length = %d, but should be: %d", hashTable.length, actualLength)
	}

}

func TestGetPrimeNextSizeUp(t *testing.T) {
	candidate := uint64(23)
	expected := uint64(23)
	prime := getPrime(candidate, true)

	if prime != expected {
		t.Errorf("expected computed prime to be: %d, got: %d", expected, prime)
	}

	candidate = uint64(46)
	expected = uint64(53)
	prime = getPrime(candidate, true)
	if prime != expected {
		t.Errorf("expected computed prime to be: %d, got: %d", expected, prime)
	}

	candidate = uint64(1610612741 * 2)
	prime = getPrime(candidate, true)
	expected = uint64(3221225533)
	if prime != expected {
		t.Errorf("expected computed prime to be: %d, got: %d", expected, prime)
	}

}

func TestInsert(t *testing.T) {
	hashTable := New(10)
	key := "foo-1"
	hashTable.Insert(key, 500)
	if hashTable.activeSlotCounter != 1 {
		t.Errorf("HashTable activeSlotCountter = %d, but it should equal 1", hashTable.activeSlotCounter)
	}

}

func TestSearch(t *testing.T) {

	hashTable := New(10)
	key := "foo-1"
	value, err := hashTable.Search(key)
	if value != nil {
		t.Errorf(`Search(%s) = %v, want nil, error: %v`, key, value, err)
	}
	targetValue := 500

	hashTable.Insert(key, targetValue)
	valueB, errB := hashTable.Search(key)

	if valueB == nil {
		t.Errorf(`Search(%s) = nil, want %v, error: %v`, key, targetValue, errB)
	}
}

func TestDelete(t *testing.T) {
	hashTable := New(10)
	key := "foo-1"
	err := hashTable.Delete(key)
	if err == nil {
		t.Errorf(`Delete(%s) = nil, expected: Delete(%s) = error,  error: %v`, key, key, err)
	}
}

func TestResizeUp(t *testing.T) {
	var desiredLength uint64 = 40
	var expectedLength uint64 = 193
	hashTable := New(desiredLength)
	totalItems := 35
	for i := 0; i < totalItems; i++ {
		hashTable.Insert(fmt.Sprintf("foo-%d", i), i*2)
	}
	if hashTable.length != expectedLength {
		t.Errorf("HashTable length = %d, but expected: %d", hashTable.length, expectedLength)
	}

	if hashTable.activeSlotCounter != uint64(totalItems) {
		t.Errorf("HashTable activeSlotCounter =  %d, but expected: %d", hashTable.activeSlotCounter, totalItems)
	}
	if hashTable.occupiedSlotCounter != uint64(totalItems) {
		t.Errorf("HashTable occupiedSlotCounter =  %d, but expected: %d", hashTable.occupiedSlotCounter, totalItems)
	}
}

func TestResizeDown(t *testing.T) {
	var desiredLength uint64 = 40
	hashTable := New(desiredLength)
	var expectedLength uint64 = 23

	for i := 0; i < 3; i++ {
		hashTable.Insert(fmt.Sprintf("foo-%d", i), i*2)
	}

	hashTable.Delete("foo-0")

	if hashTable.length != expectedLength {

		t.Errorf("HashTable length = %d, but expected: %d", hashTable.length, expectedLength)
	}

	expectedActiveSlotCounter := 2
	expectedOccupiedSlotCounter := 2
	if hashTable.activeSlotCounter != uint64(expectedActiveSlotCounter) {
		t.Errorf("HashTable activeSlotCounter =  %d, but expected: %d", hashTable.activeSlotCounter, expectedActiveSlotCounter)
	}
	if hashTable.occupiedSlotCounter != 2 {
		t.Errorf("HashTable occupiedSlotCounter =  %d, but expected: %d", hashTable.occupiedSlotCounter, expectedOccupiedSlotCounter)
	}
}

func TestUpdateValue(t *testing.T) {
	var length uint64 = 10
	hashTable := New(length)
	key := "foo-1"
	hashTable.Insert(key, 100)
	hashTable.Insert(key, 200)
	value, err := hashTable.Search(key)
	if err != nil {
		t.Errorf("Expected 200, got error: %v", err)
	}
	if value != 200 {
		t.Errorf("Value should equal 200, got: %d", value)
	}
}

func TestProbingWhenInserting(t *testing.T) {
	var length uint64 = 5
	hashTable := New(length)
	keyA := "foo-1"   // This would hash to 3
	keyB := "foo-111" // This would hash to 3
	hashTable.Insert(keyA, 200)

	hashTable.Insert(keyB, 400)
	value, err := hashTable.Search(keyB)
	if value == nil {
		t.Errorf(`Search(%s) = nil, want %d, error: %v`, keyB, value, err)
	}
}

func BenchmarkSearchExistingKey(b *testing.B) {
	var length uint64 = 2000000
	key := "foo-3300"
	totalItems := 1000000
	hashTable := New(length)
	for i := 0; i < totalItems; i++ {
		guid := guid.New()
		hashTable.Insert(guid.String(), i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		value, err := hashTable.Search(key)
		if value == nil {
			b.Errorf(`Search(%s) = %v, want 3300, error: %v`, key, value, err)
		}
	}
}

func BenchmarkNonExistingKey(b *testing.B) {
	var length uint64 = 2000000
	key := "foo-%300"
	totalItems := 1000000
	hashTable := New(length)
	for i := 0; i < totalItems; i++ {
		guid := guid.New()
		hashTable.Insert(guid.String(), i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		value, err := hashTable.Search(key)
		if value != nil {
			b.Errorf(`Search(%s) = %v, want nil, error: %v`, key, value, err)
		}
	}
}

func BenchmarkInsertNoResize(b *testing.B) {

	// 255 ms -> with pointers to data structs
	// 222 ms -> with actual data structs

	// 65.980232 MBs -> with pointers to data structs
	// 15.98 Mbs -> with actual data structs

	// with isSoftDeleted still present in the data struct  and computing hashFnv every time
	// BenchmarkInsertNoResize-12           100         271884392 ns/op        10718200 B/op    1089742 allocs/op

	// without isSoftDeleted in the data struct and computing hashFnv every time
	// BenchmarkInsertNoResize-12           100         267893551 ns/op        10718226 B/op    1089742 allocs/op

	// // without isSoftDeleted in the data struct and computing fnv hash only once by using nodeKey
	// BenchmarkInsertNoResize-12           100         258080503 ns/op        10718210 B/op    1089742 allocs/op
	totalKeys := 1000_000
	keys := make([]string, totalKeys)

	for i := 0; i < totalKeys; i++ {
		guid := guid.New()
		keys[i] = guid.String()
	}

	var tableLength uint64 = 2000_000

	for i := 0; i < b.N; i++ {
		b.StopTimer()
		table := New(tableLength)
		b.StartTimer()
		for i, key := range keys {
			table.Insert(key, i)
		}

	}

}

// using key: fmt.Sprintf("foudhiuwediwuendiw1uend834u2390u3029u402--4-4-423e23eo-%d", i*1000000000000)
// Before modifying fnv hash
// BenchmarkInsertNoResize-12           100         300636513 ns/op         9040243 B/op    1019756 allocs/op

// After modifying fnv hash
//BenchmarkInsertNoResize-12           100         298962644 ns/op         9040275 B/op    1019756 allocs/op

func BenchmarkFnvHash1a(b *testing.B) {
	for _, totalKeys := range []int{1000, 50_000, 1000_000} {
		b.Run(fmt.Sprintf("keys = %d", totalKeys), func(b *testing.B) {
			keys := make([]string, totalKeys)
			for i := range keys {
				guid := guid.New()
				keys[i] = guid.String()
			}
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				for _, key := range keys {
					fnvHash(key)
				}
			}
		})
	}
}
func BenchmarkLibFnvHash1a(b *testing.B) {

	for _, totalKeys := range []int{1000, 50_000, 1000_000} {
		b.Run(fmt.Sprintf("keys = %d", totalKeys), func(b *testing.B) {
			keys := make([]string, totalKeys)
			for i := range keys {
				guid := guid.New()
				keys[i] = guid.String()
			}
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				for _, key := range keys {
					fnvHashLib(key)
				}
			}
		})
	}
}
