package exercise10

import "testing"

func TestNewList_EmptySlice(t *testing.T) {
	s := []int{}
	list := NewList(s)
	if list != nil {
		t.Errorf("NewList with empty slice should return nil, got %+v", list)
	}
}

func TestNewList_SingleElement(t *testing.T) {
	s := []string{"hello"}
	list := NewList(s)
	if list == nil {
		t.Fatal("NewList with single element should not return nil")
	}
	if list.val != "hello" {
		t.Errorf("Expected val to be 'hello', got '%v'", list.val)
	}
	if list.next != nil {
		t.Errorf("Expected next to be nil, got %+v", list.next)
	}
}

func TestNewList_MultipleElements(t *testing.T) {
	s := []int{1, 2, 3}
	list := NewList(s)
	if list == nil {
		t.Fatal("NewList with multiple elements should not return nil")
	}

	if list.val != 1 {
		t.Errorf("Expected first val to be 1, got %v", list.val)
	}
	if list.next == nil {
		t.Fatal("Expected first next to not be nil")
	}

	if list.next.val != 2 {
		t.Errorf("Expected second val to be 2, got %v", list.next.val)
	}
	if list.next.next == nil {
		t.Fatal("Expected second next to not be nil")
	}

	if list.next.next.val != 3 {
		t.Errorf("Expected third val to be 3, got %v", list.next.next.val)
	}
	if list.next.next.next != nil {
		t.Errorf("Expected third next to be nil, got %+v", list.next.next.next)
	}
}

func TestNewList_DifferentTypes(t *testing.T) {
	stringSlice := []string{"a", "b"}
	stringList := NewList(stringSlice)
	if stringList == nil || stringList.val != "a" || stringList.next == nil || stringList.next.val != "b" || stringList.next.next != nil {
		t.Errorf("NewList with string slice failed: %+v", stringList)
	}

	floatSlice := []float64{1.1, 2.2}
	floatList := NewList(floatSlice)
	if floatList == nil || floatList.val != 1.1 || floatList.next == nil || floatList.next.val != 2.2 || floatList.next.next != nil {
		t.Errorf("NewList with float slice failed: %+v", floatList)
	}

	boolSlice := []bool{true, false}
	boolList := NewList(boolSlice)
	if boolList == nil || boolList.val != true || boolList.next == nil || boolList.next.val != false || boolList.next.next != nil {
		t.Errorf("NewList with bool slice failed: %+v", boolList)
	}
}

func TestNewList_PointerTypes(t *testing.T) {
	type CustomType struct {
		Value int
	}
	ptr1 := &CustomType{Value: 10}
	ptr2 := &CustomType{Value: 20}
	ptrSlice := []*CustomType{ptr1, ptr2}
	ptrList := NewList(ptrSlice)

	if ptrList == nil || ptrList.val != ptr1 || ptrList.next == nil || ptrList.next.val != ptr2 || ptrList.next.next != nil {
		t.Errorf("NewList with pointer slice failed: %+v", ptrList)
	}
	if ptrList.val.Value != 10 || ptrList.next.val.Value != 20 {
		t.Errorf("NewList with pointer slice has incorrect values: %+v", ptrList)
	}
}

func TestNewList_NilSlice(t *testing.T) {
	var s []int
	list := NewList(s)
	if list != nil {
		t.Errorf("NewList with nil slice should return nil, got %+v", list)
	}
}
