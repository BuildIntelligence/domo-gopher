package domo

import (
	"reflect"
	"testing"
	"time"
)

type DomoSample struct {
	Foo string `domo:""`
	Bar int `domo:"bar"`
	Baz float64 `domo:"baz,Baz"`
	IgnoreFooBar string `domo:"-"`
	BazBar int
	OptionalBar *int `domo:"obar,omitempty"`
	OptionalBaz *float32 `domo:"obaz, omitempty"`
}
type DomoNestedSample struct {
	Blah float64 `domo:"DECIMAL"`
	FirstBlahDay time.Time `domo:"firstBlahDay,DATE"`
	FirstBlahTime time.Time `domo:"firstBlahTime"`
	Sample DomoSample `domo:"-"` // ignore making a field called "Sample". Still gets DomoSample Fields
}
type DomoEmbeddedSample struct {
	Blah float64 `domo:"DECIMAL"`
	FirstBlahDay time.Time `domo:"firstBlahDay,DATE"`
	FirstBlahTime time.Time `domo:"firstBlahTime"`
	DomoSample // Anonymous Embedded struct, will get DomoSample fieldInfo
}
func TestGenerateDataSetSchema(t *testing.T) {
	expectedColTypes := []string{"STRING","LONG","DOUBLE","LONG","LONG","DOUBLE"}
	expectedColNames := []string{"Foo","bar","baz","BazBar","obar","obaz"}
	domoSampleSchema := GenerateDataSetSchema(reflect.TypeOf(DomoSample{}))
	if len(domoSampleSchema.Columns) != 6 {
		t.Fatalf("Expected 6 columns but got %d\n%v", len(domoSampleSchema.Columns), domoSampleSchema.Columns)
	}
	for i, col := range domoSampleSchema.Columns {
		if col.ColumnType != expectedColTypes[i]	{
			t.Fatalf("Expected column %d (%s) to be type %s but got %s", i, col.Name, expectedColTypes[i], col.ColumnType)
		}
		if col.Name != expectedColNames[i]	{
			t.Fatalf("Expected column %d to be named %s but got %s", i, col.Name, expectedColNames[i])
		}
	}
}

func TestGenerateDataSetSchema_NestedStruct(t *testing.T) {
	expectedColTypes := []string{"DECIMAL","DATE","DATETIME","STRING","LONG","DOUBLE","LONG","LONG","DOUBLE"}
	expectedColNames := []string{"Blah","firstBlahDay","firstBlahTime","Foo","bar","baz","BazBar","obar","obaz"}
	domoSampleSchema := GenerateDataSetSchema(reflect.TypeOf(DomoNestedSample{}))
	if len(domoSampleSchema.Columns) != 9 {
		t.Fatalf("Expected 9 columns but got %d\n%v", len(domoSampleSchema.Columns), domoSampleSchema.Columns)
	}
	for i, col := range domoSampleSchema.Columns {
		if col.ColumnType != expectedColTypes[i]	{
			t.Fatalf("Expected column %d (%s) to be type %s but got %s", i, col.Name, expectedColTypes[i], col.ColumnType)
		}
		if col.Name != expectedColNames[i]	{
			t.Fatalf("Expected column %d to be named %s but got %s", i, col.Name, expectedColNames[i])
		}
	}
}

func TestGenerateDataSetSchema_EmbedStruct(t *testing.T) {
	expectedColTypes := []string{"DECIMAL","DATE","DATETIME","STRING","LONG","DOUBLE","LONG","LONG","DOUBLE"}
	expectedColNames := []string{"Blah","firstBlahDay","firstBlahTime","Foo","bar","baz","BazBar","obar","obaz"}
	domoSampleSchema := GenerateDataSetSchema(reflect.TypeOf(DomoEmbeddedSample{}))
	if len(domoSampleSchema.Columns) != 9 {
		t.Fatalf("Expected 9 columns but got %d\n%v", len(domoSampleSchema.Columns), domoSampleSchema.Columns)
	}
	for i, col := range domoSampleSchema.Columns {
		if col.ColumnType != expectedColTypes[i]	{
			t.Fatalf("Expected column %d (%s) to be type %s but got %s", i, col.Name, expectedColTypes[i], col.ColumnType)
		}
		if col.Name != expectedColNames[i]	{
			t.Fatalf("Expected column %d to be named %s but got %s", i, col.Name, expectedColNames[i])
		}
	}
}
