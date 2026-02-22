package main

import (
	"fmt"
	"testing"

	"github.com/beevik/guid"
)

func makeSequentialKeys(total int) []string {
	keys := make([]string, total)
	for i := 0; i < total; i++ {
		keys[i] = fmt.Sprintf("key-%d", i)
	}
	return keys
}

func buildHashTable(keys []string, tableLength uint64) *HashTable[int] {
	table := New[int](tableLength)
	for i, key := range keys {
		table.Insert(key, i)
	}
	return table
}

func buildGoMap(keys []string, mapCapacity int) map[string]int {
	m := make(map[string]int, mapCapacity)
	for i, key := range keys {
		m[key] = i
	}
	return m
}



func findCollidingKeys(targetCount int, tableLength uint64) []string {
	buckets := make(map[uint64][]string)
	maxCandidates := 2_000_000
	for i := 0; i < maxCandidates; i++ {
		key := fmt.Sprintf("probe-key-%d", i)
		idx := NewKey(key).hash % tableLength
		bucket := append(buckets[idx], key)
		if len(bucket) >= targetCount {
			return bucket
		}
		buckets[idx] = bucket
	}

	panic("could not find enough colliding keys for probing benchmark")
}

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
	hashTable := New[int](uint64(targetLength))
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
	hashTable := New[int](10)
	key := "foo-1"
	hashTable.Insert(key, 500)
	if hashTable.activeSlotCounter != 1 {
		t.Errorf("HashTable activeSlotCountter = %d, but it should equal 1", hashTable.activeSlotCounter)
	}

}

func TestSearch(t *testing.T) {

	hashTable := New[int](10)
	key := "foo-1"
	value, err := hashTable.Search(key)
	if err == nil || value != 0 {
		t.Errorf(`Search(%s) = %v, want 0 with error, error: %v`, key, value, err)
	}
	targetValue := 500

	hashTable.Insert(key, targetValue)
	valueB, errB := hashTable.Search(key)

	if errB != nil || valueB != targetValue {
		t.Errorf(`Search(%s) = %v, want %v, error: %v`, key, valueB, targetValue, errB)
	}
}

func TestDelete(t *testing.T) {
	hashTable := New[int](10)
	key := "foo-1"
	err := hashTable.Delete(key)
	if err == nil {
		t.Errorf(`Delete(%s) = nil, expected: Delete(%s) = error,  error: %v`, key, key, err)
	}
}

func TestResizeUp(t *testing.T) {
	var desiredLength uint64 = 40
	var expectedLength uint64 = 193
	hashTable := New[int](desiredLength)
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
	hashTable := New[int](desiredLength)
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
	hashTable := New[int](length)
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
	hashTable := New[int](length)
	keyA := "foo-1"   // This would hash to 3
	keyB := "foo-111" // This would hash to 3
	hashTable.Insert(keyA, 200)

	hashTable.Insert(keyB, 400)
	value, err := hashTable.Search(keyB)
	if err != nil || value != 400 {
		t.Errorf(`Search(%s) = %v, want %d, error: %v`, keyB, value, 400, err)
	}
}

func BenchmarkSearchExistingKey(b *testing.B) {
	totalItems := 1_000_000
	keys := makeSequentialKeys(totalItems)
	targetKey := keys[totalItems/2]
	table := buildHashTable(keys, 2_000_000)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		value, err := table.Search(targetKey)
		if err != nil || value != totalItems/2 {
			b.Fatalf(`Search(%s) expected %d, got value=%v error=%v`, targetKey, totalItems/2, value, err)
		}
	}
}

func BenchmarkSearchNonExistingKey(b *testing.B) {
	totalItems := 1_000_000
	keys := makeSequentialKeys(totalItems)
	table := buildHashTable(keys, 2_000_000)
	missingKey := "key-missing"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		value, err := table.Search(missingKey)
		if err == nil || value != 0 {
			b.Fatalf(`Search(%s) expected not found, got value=%v error=%v`, missingKey, value, err)
		}
	}
}

func BenchmarkGoMapSearchExistingKey(b *testing.B) {
	totalItems := 1_000_000
	keys := makeSequentialKeys(totalItems)
	targetKey := keys[totalItems/2]
	m := buildGoMap(keys, totalItems*2)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		value, ok := m[targetKey]
		if !ok || value != totalItems/2 {
			b.Fatalf("expected key %s to exist with value %d, got value=%d exists=%v", targetKey, totalItems/2, value, ok)
		}
	}
}

func BenchmarkGoMapSearchNonExistingKey(b *testing.B) {
	totalItems := 1_000_000
	keys := makeSequentialKeys(totalItems)
	m := buildGoMap(keys, totalItems*2)
	missingKey := "key-missing"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, ok := m[missingKey]
		if ok {
			b.Fatalf("expected key %s to be missing", missingKey)
		}
	}
}

func BenchmarkDeleteExistingKey(b *testing.B) {
	totalItems := 200_000
	keys := makeSequentialKeys(totalItems)
	targetIndex := totalItems / 2
	targetKey := keys[targetIndex]

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		table := buildHashTable(keys, 400_000)
		b.StartTimer()
		err := table.Delete(targetKey)
		if err != nil {
			b.Fatalf("Delete(%s) expected nil error, got %v", targetKey, err)
		}
	}
}

func BenchmarkGoMapDeleteExistingKey(b *testing.B) {
	totalItems := 200_000
	keys := makeSequentialKeys(totalItems)
	targetIndex := totalItems / 2
	targetKey := keys[targetIndex]

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		m := buildGoMap(keys, totalItems*2)
		b.StartTimer()
		delete(m, targetKey)
		if _, ok := m[targetKey]; ok {
			b.Fatalf("expected key %s to be deleted", targetKey)
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
	totalKeys := 1_000_000
	keys := makeSequentialKeys(totalKeys)

	var tableLength uint64 = 2_000_000

	for i := 0; i < b.N; i++ {
		b.StopTimer()
		table := New[int](tableLength)
		b.StartTimer()
		for i, key := range keys {
			table.Insert(key, i)
		}

	}

}

func BenchmarkInsertWithResize(b *testing.B) {
	totalKeys := 1_000_000
	keys := makeSequentialKeys(totalKeys)

	// Starts much smaller to force multiple resize-up operations.
	var tableLength uint64 = 100_000

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		table := New[int](tableLength)
		b.StartTimer()
		for i, key := range keys {
			table.Insert(key, i)
		}
	}
}

func BenchmarkGoMapInsertNoResize(b *testing.B) {
	totalKeys := 1_000_000
	keys := makeSequentialKeys(totalKeys)
	mapCapacity := 2_000_000

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		m := make(map[string]int, mapCapacity)
		b.StartTimer()
		for i, key := range keys {
			m[key] = i
		}
	}
}

func BenchmarkGoMapInsertWithResize(b *testing.B) {
	totalKeys := 1_000_000
	keys := makeSequentialKeys(totalKeys)
	mapCapacity := 100_000

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		m := make(map[string]int, mapCapacity)
		b.StartTimer()
		for i, key := range keys {
			m[key] = i
		}
	}
}

func BenchmarkProbingHeavyCollisionSearch(b *testing.B) {
	var tableLength uint64 = 389
	keys := findCollidingKeys(200, tableLength)
	table := buildHashTable(keys, tableLength)
	target := keys[len(keys)-1]

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		v, err := table.Search(target)
		if err != nil || v != len(keys)-1 {
			b.Fatalf("Search(%s) expected %d, got value=%d error=%v", target, len(keys)-1, v, err)
		}
	}
}

func BenchmarkGoMapProbingHeavyCollisionSearch(b *testing.B) {
	var tableLength uint64 = 389
	keys := findCollidingKeys(200, tableLength)
	m := buildGoMap(keys, len(keys)*2)
	target := keys[len(keys)-1]

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		v, ok := m[target]
		if !ok || v != len(keys)-1 {
			b.Fatalf("expected key %s to exist with value %d, got value=%d exists=%v", target, len(keys)-1, v, ok)
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
