package struct2

import (
	"reflect"
	"testing"
)

func TestJSONInlineTag(t *testing.T) {
	type TestCase struct {
		name    string
		setup   func() (interface{}, map[string]interface{})
		verify  func(t *testing.T, result interface{})
		wantErr bool
	}

	tests := []TestCase{
		{
			name: "basic embedded struct",
			setup: func() (interface{}, map[string]interface{}) {
				type Embedded struct {
					Field1 string `json:"field1"`
				}
				type TestStruct struct {
					Embedded `json:",inline"`
					Field2   string `json:"field2"`
				}
				return &TestStruct{}, map[string]interface{}{
					"field1": "value1",
					"field2": "value2",
				}
			},
			verify: func(t *testing.T, result interface{}) {
				v := reflect.ValueOf(result).Elem()
				field1 := v.FieldByName("Field1").String()
				field2 := v.FieldByName("Field2").String()

				if field1 != "value1" || field2 != "value2" {
					t.Errorf("Expected {Field1:value1 Field2:value2}, got {Field1:%s Field2:%s}", field1, field2)
				}
			},
			wantErr: false,
		},
		{
			name: "nested struct with inline",
			setup: func() (interface{}, map[string]interface{}) {
				type Nested struct {
					NestedField string `json:"nested_field"`
				}
				type Embedded struct {
					Field1 string `json:"field1"`
					Nested Nested `json:"nested"`
				}
				type TestStruct struct {
					Embedded `json:",inline"`
					Field2   string `json:"field2"`
				}
				return &TestStruct{}, map[string]interface{}{
					"field1": "value1",
					"field2": "value2",
					"nested": map[string]interface{}{
						"nested_field": "nested_value",
					},
				}
			},
			verify: func(t *testing.T, result interface{}) {
				v := reflect.ValueOf(result).Elem()
				field1 := v.FieldByName("Field1").String()
				nestedField := v.FieldByName("Nested").FieldByName("NestedField").String()
				field2 := v.FieldByName("Field2").String()

				if field1 != "value1" || nestedField != "nested_value" || field2 != "value2" {
					t.Errorf("Expected {Field1:value1 Nested:{NestedField:nested_value} Field2:value2}, got {Field1:%s Nested:{NestedField:%s} Field2:%s}",
						field1, nestedField, field2)
				}
			},
			wantErr: false,
		},
		{
			name: "multiple embedded structs",
			setup: func() (interface{}, map[string]interface{}) {
				type Embedded1 struct {
					Field1 string `json:"field1"`
				}
				type Embedded2 struct {
					Field2 string `json:"field2"`
				}
				type TestStruct struct {
					Embedded1 `json:",inline"`
					Embedded2 `json:",inline"`
					Field3    string `json:"field3"`
				}
				return &TestStruct{}, map[string]interface{}{
					"field1": "value1",
					"field2": "value2",
					"field3": "value3",
				}
			},
			verify: func(t *testing.T, result interface{}) {
				v := reflect.ValueOf(result).Elem()
				field1 := v.FieldByName("Field1").String()
				field2 := v.FieldByName("Field2").String()
				field3 := v.FieldByName("Field3").String()

				if field1 != "value1" || field2 != "value2" || field3 != "value3" {
					t.Errorf("Expected {Field1:value1 Field2:value2 Field3:value3}, got {Field1:%s Field2:%s Field3:%s}",
						field1, field2, field3)
				}
			},
			wantErr: false,
		},
		{
			name: "transitive inline structs",
			setup: func() (interface{}, map[string]interface{}) {
				type Level3 struct {
					Field3 string `json:"field3"`
				}

				type Level2 struct {
					Level3 `json:",inline"`
					Field2 string `json:"field2"`
				}

				type Level1 struct {
					Level2 `json:",inline"`
					Field1 string `json:"field1"`
				}

				type TestStruct struct {
					Level1 `json:",inline"`
					Field0 string `json:"field0"`
				}

				return &TestStruct{}, map[string]interface{}{
					"field0": "value0",
					"field1": "value1",
					"field2": "value2",
					"field3": "value3",
				}
			},
			verify: func(t *testing.T, result interface{}) {
				v := reflect.ValueOf(result).Elem()
				field0 := v.FieldByName("Field0").String()
				field1 := v.FieldByName("Field1").String()
				field2 := v.FieldByName("Field2").String()
				field3 := v.FieldByName("Field3").String()

				expected := field0 == "value0" && field1 == "value1" &&
					field2 == "value2" && field3 == "value3"

				if !expected {
					t.Errorf("Expected {Field0:value0 Field1:value1 Field2:value2 Field3:value3}, got {Field0:%s Field1:%s Field2:%s Field3:%s}",
						field0, field1, field2, field3)
				}
			},
			wantErr: false,
		},
		{
			name: "transitive_composite_structs",
			setup: func() (interface{}, map[string]interface{}) {
				type Level3 struct {
					Field3 string `json:"field3"`
				}

				type Level2 struct {
					Nested Level3 `json:"nested,inline"`
					Field2 string `json:"field2"`
				}

				type Level1 struct {
					Nested Level2 `json:"nested,inline"`
					Field1 string `json:"field1"`
				}

				type TestStruct struct {
					Nested Level1 `json:"nested,inline"`
					Field0 string `json:"field0"`
				}

				return &TestStruct{}, map[string]interface{}{
					"field0": "value0",
					"nested": map[string]interface{}{
						"field1": "value1",
						"nested": map[string]interface{}{
							"field2": "value2",
							"nested": map[string]interface{}{
								"field3": "value3",
							},
						},
					},
				}
			},
			verify: func(t *testing.T, result interface{}) {
				v := reflect.ValueOf(result).Elem()
				field0 := v.FieldByName("Field0").String()

				nested1 := v.FieldByName("Nested")
				field1 := nested1.FieldByName("Field1").String()
				nested2 := nested1.FieldByName("Nested")
				field2 := nested2.FieldByName("Field2").String()
				nested3 := nested2.FieldByName("Nested")
				field3 := nested3.FieldByName("Field3").String()

				if field0 != "value0" || field1 != "value1" || field2 != "value2" || field3 != "value3" {
					t.Errorf("Expected {Field0:value0 Nested:{Field1:value1 Nested:{Field2:value2 Nested:{Field3:value3}}}}, got %+v", result)
				}
			},
			wantErr: false,
		},
		{
			name: "pointer fields",
			setup: func() (interface{}, map[string]interface{}) {
				type Level2 struct {
					Field2 string `json:"field2"`
				}
				type Level1 struct {
					Nested *Level2 `json:"nested,inline"`
					Field1 string  `json:"field1"`
				}
				type TestStruct struct {
					Nested *Level1 `json:"nested,inline"`
					Field0 string  `json:"field0"`
				}
				return &TestStruct{}, map[string]interface{}{
					"field0": "value0",
					"field1": "value1",
					"field2": "value2",
				}
			},
			verify: func(t *testing.T, result interface{}) {
				v := reflect.ValueOf(result).Elem()
				field0 := v.FieldByName("Field0").String()
				nested1 := v.FieldByName("Nested").Elem()
				field1 := nested1.FieldByName("Field1").String()
				nested2 := nested1.FieldByName("Nested").Elem()
				field2 := nested2.FieldByName("Field2").String()

				if field0 != "value0" || field1 != "value1" || field2 != "value2" {
					t.Errorf("Expected {Field0:value0 Nested:{Field1:value1 Nested:{Field2:value2}}}, got {Field0:%s Nested:{Field1:%s Nested:{Field2:%s}}}",
						field0, field1, field2)
				}
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			decoder := &Decoder{
				TagName: "json",
			}

			result, input := tt.setup()
			err := decoder.Decode(input, result)

			if (err != nil) != tt.wantErr {
				t.Fatalf("Decode() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr {
				tt.verify(t, result)
			}
		})
	}
}

func Test_interface2StructValue(t *testing.T) {
	type args struct {
		s interface{}
	}
	tests := []struct {
		name      string
		args      args
		wantPanic bool
	}{
		{
			name: "panic not struct value",
			args: args{
				s: "not struct",
			},
			wantPanic: true,
		},
		{
			name: "struct value",
			args: args{
				s: struct {
					name int
				}{
					name: 123,
				},
			},
			wantPanic: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); (r != nil) != tt.wantPanic {
					t.Errorf("panic %t, wantPanic %t", (r != nil), tt.wantPanic)
				}
			}()

			value2StructValue(reflect.ValueOf(tt.args.s))
		})
	}
}
