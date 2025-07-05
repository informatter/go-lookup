package main

import (
	"fmt"
	"testing"
)

func TestCreateHashTable(t *testing.T) {
	targetLength := 10
	actualLength := 17
	hashTable := New(uint64(targetLength))
	if hashTable.length != uint64(actualLength) {
		t.Errorf("HashTable length = %d, but should be: %d", hashTable.length, actualLength)
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
	expectedOccupiedSlotCounter :=2
	if hashTable.activeSlotCounter != uint64(expectedActiveSlotCounter) {
		t.Errorf("HashTable activeSlotCounter =  %d, but expected: %d", hashTable.activeSlotCounter, expectedActiveSlotCounter)
	}
	if hashTable.occupiedSlotCounter != 2 {
		t.Errorf("HashTable occupiedSlotCounter =  %d, but expected: %d", hashTable.occupiedSlotCounter, expectedOccupiedSlotCounter)
	}
}

func TestUpdateValue( t *testing.T){
	var length uint64 = 10
	hashTable := New(length)
	key := "foo-1"
	hashTable.Insert(key,100)
	hashTable.Insert(key, 200)
	value, err := hashTable.Search(key)
	if err != nil{
		t.Errorf("Expected 200, got error: %v", err)
	}
	if value != 200{
		t.Errorf("Value should equal 200, got: %d",value)
	}
}

func TestProbingWhenInserting( t *testing.T){
	var length uint64 = 5
	hashTable := New(length)
	keyA:="foo-1" // This would hash to 3
	keyB :="foo-111" // This would hash to 3
	hashTable.Insert(keyA,200)

	hashTable.Insert(keyB,400)
	value, err := hashTable.Search(keyB)
	if value == nil {
		t.Errorf(`Search(%s) = nil, want %d, error: %v`, keyB, value, err)
	}
}

func BenchmarkSearchExistingKey(b *testing.B) {
	var length uint64 = 2000000
	key := "foo-3300"
	totalItems:= 1000000
	hashTable := New(length)
	for i := 0; i < totalItems; i++ {
		hashTable.Insert(fmt.Sprintf("foo-%d", i), i)
	}
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        value, err := hashTable.Search(key)
		if value == nil{
			b.Errorf(`Search(%s) = %v, want 3300, error: %v`, key, value, err)
		}
    }
}

func BenchmarkNonExistingKey(b *testing.B) {
	var length uint64 = 2000000
	key := "foo-%300"
	totalItems:= 1000000
	hashTable := New(length)
	for i := 0; i < totalItems; i++ {
		hashTable.Insert(fmt.Sprintf("foo-%d", i), i)
	}
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        value, err := hashTable.Search(key)
		if value != nil{
			b.Errorf(`Search(%s) = %v, want nil, error: %v`, key, value, err)
		}
    }
}

func BenchmarkInsertNoResize(b *testing.B) {

	// 255 ms -> with pointers to data structs
	// 222 ms -> with actual data structs

	// 65.980232 MBs -> with pointers to data structs
	// 15.98 Mbs -> with actual data structs


	totalKeys := 1000_000
	keys := make([]string, totalKeys)
	
	for i := range totalKeys{
		keys[i] = fmt.Sprintf("foo-%d",i)
	}

	var tableLength uint64 = 2000_000

	for i := 0; i < b.N; i++ {
		b.StopTimer()
		table := New(tableLength)
		b.StartTimer()
		for i,key:=range keys{
			table.Insert(key,i)
		}

		b.Logf("Total collisions: %d",table.collistionCount)

	}

}