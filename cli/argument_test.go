package cli

// TODO: neaten up some of these tests
// ie testList := []struct{} style

import (
	"reflect"
	"testing"
)

func TestArgument_HasValue(t *testing.T) {
	arg := Argument{Name: "test"}
	if arg.HasValue() {
		t.Fatalf("expected false, got true")
	}
	arg.value = "blah"
	if !arg.HasValue() {
		t.Fatalf("expected true, got false")
	}
}

func TestArgument_Value(t *testing.T) {
	arg := Argument{Name: "test"}
	if arg.Value() != nil {
		t.Fatalf("expected nil")
	}
	arg.value = "blah"
	if arg.Value() == nil {
		t.Fatal("expected not nil")
	}
	if !reflect.DeepEqual(arg.Value(), "blah") {
		t.Fatalf("not equal")
	}
}

func TestArgument_Set(t *testing.T) {
	arg := Argument{Name: "test"}
	if err := arg.Set("blah"); err == nil {
		t.Fatalf("expected error")
	}
	arg.TakesValue = true
	if err := arg.Set("blah"); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	arg2 := Argument{Name: "test", TakesValue: true, TakesMultiple: true}
	if err := arg2.Set("blah1"); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if err := arg2.Set("blah2"); err != nil {
		t.Fatalf("expected no error")
	}
	if !reflect.DeepEqual(arg2.StringSlice(), []string{"blah1", "blah2"}) {
		t.Fatalf("values not equal")
	}

	arg3 := Argument{Name: "test", TakesValue: true, TakesMultiple: true}
	if err := arg3.Set("blah1,blah2"); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if !reflect.DeepEqual(arg3.StringSlice(), []string{"blah1", "blah2"}) {
		t.Fatalf("values not equal")
	}

	arg4 := Argument{Name: "test", TakesValue: true, TakesMultiple: true}
	if err := arg4.Set("blah1, blah2"); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if !reflect.DeepEqual(arg4.StringSlice(), []string{"blah1", "blah2"}) {
		t.Fatalf("values not equal")
	}
}

func TestArgument_String(t *testing.T) {
	arg := Argument{Name: "test"}
	if arg.String() != "" {
		t.Fatalf("expected empty")
	}
	arg.value = "blah"
	if arg.String() == "" {
		t.Fatal("expected not empty")
	}
	if !reflect.DeepEqual(arg.String(), "blah") {
		t.Fatalf("not equal")
	}
}

func TestArgument_StringOrDefault(t *testing.T) {
	arg := Argument{Name: "test", DefaultValue: SimpleValue{"default", "blah"}}
	if arg.StringOrDefault() == "" {
		t.Fatalf("expected not empty")
	}
	arg.value = "blah"
	if arg.StringOrDefault() == "" {
		t.Fatal("expected not empty")
	}
	if !reflect.DeepEqual(arg.StringOrDefault(), "blah") {
		t.Fatalf("not equal")
	}
}

func TestArgument_StringSlice(t *testing.T) {
	arg := Argument{Name: "test", TakesValue: true, TakesMultiple: true}

	if err := arg.Set("blah1"); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	val := arg.StringSlice()
	if len(val) != 1 {
		t.Fatalf("expected len 1, got: %d", len(val))
	}
	if !reflect.DeepEqual(val, []string{"blah1"}) {
		t.Fatalf("not equal")
	}

	if err := arg.Set("blah2"); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	val = arg.StringSlice()
	if len(val) != 2 {
		t.Fatalf("expected len 2, got: %d", len(val))
	}
	if !reflect.DeepEqual(val, []string{"blah1", "blah2"}) {
		t.Fatalf("not equal")
	}
}

func TestArgument_StringSliceOrDefault(t *testing.T) {
	arg := Argument{Name: "test", DefaultValue: SimpleValue{[]string{"default"}, "blah"}, TakesValue: true, TakesMultiple: true}

	val := arg.StringSliceOrDefault()
	if !reflect.DeepEqual(val, []string{"default"}) {
		t.Fatalf("not equal")
	}

	if err := arg.Set("notdefault"); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	val = arg.StringSliceOrDefault()
	if !reflect.DeepEqual(val, []string{"notdefault"}) {
		t.Fatalf("not equal")
	}
}

func TestArgument_Bool(t *testing.T) {
	arg := Argument{
		Name:  "test",
		value: true,
	}
	if arg.Bool() != true {
		t.Fatalf("expected true, got false")
	}
}

func TestArgument_BoolOrDefault(t *testing.T) {
	arg := Argument{
		Name: "test",
	}
	if arg.BoolOrDefault() {
		t.Fatalf("expected false, got true")
	}
	arg.DefaultValue = SimpleValue{nil, ""}
	if arg.BoolOrDefault() {
		t.Fatalf("expected false, got true")
	}
	arg.IsPresent = true
	if arg.BoolOrDefault() {
		t.Fatalf("expected false, got true")
	}
	arg.IsPresent = false
	arg.DefaultValue = SimpleValue{true, "trueish"}
	if !arg.BoolOrDefault() {
		t.Fatalf("expected true, got false")
	}
}
