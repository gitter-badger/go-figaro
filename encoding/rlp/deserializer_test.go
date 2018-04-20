package rlp

import (
	"reflect"
	"testing"
)

func TestDeserializeUint(t *testing.T) {
	h := uint(0)
	type args struct {
		dest *uint
		d    []byte
	}
	tests := []struct {
		name    string
		args    args
		want    uint
		wantErr bool
	}{
		{"255", args{&h, []byte{0xff}}, 255, false},
		{"0", args{&h, []byte{}}, 0, false},
		{"invalid padding", args{&h, []byte{0x00, 0xff}}, 0, true},
		{"invalid data", args{&h, []byte("cats and dogs living together")}, 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := DeserializeUint(tt.args.dest, tt.args.d); (err != nil) != tt.wantErr {
				t.Errorf("DeserializeUint() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(*tt.args.dest, tt.want) {
				t.Errorf("DeserializeUint() = %v, want %v", *tt.args.dest, tt.want)
			}
		})
	}
}

func TestDeserializeUint8(t *testing.T) {
	h := uint8(0)
	type args struct {
		dest *uint8
		d    []byte
	}
	tests := []struct {
		name    string
		args    args
		want    uint8
		wantErr bool
	}{
		{"255", args{&h, []byte{0xff}}, 255, false},
		{"0", args{&h, []byte{}}, 0, false},
		{"invalid padding", args{&h, []byte{0x00}}, 0, true},
		{"invalid data", args{&h, []byte("Human sacrifice, dogs and cats living together... mass hysteria!")}, 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := DeserializeUint8(tt.args.dest, tt.args.d); (err != nil) != tt.wantErr {
				t.Errorf("DeserializeUint8() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(*tt.args.dest, tt.want) {
				t.Errorf("DeserializeUint8() = %v, want %v", *tt.args.dest, tt.want)
			}
		})
	}
}

func TestDeserializeUint16(t *testing.T) {
	h := uint16(0)
	type args struct {
		dest *uint16
		d    []byte
	}
	tests := []struct {
		name    string
		args    args
		want    uint16
		wantErr bool
	}{
		{"255", args{&h, []byte{0xff}}, 255, false},
		{"0", args{&h, []byte{}}, 0, false},
		{"invalid padding", args{&h, []byte{0x00, 0xff}}, 0, true},
		{"invalid data", args{&h, []byte("cats and dogs living together")}, 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := DeserializeUint16(tt.args.dest, tt.args.d); (err != nil) != tt.wantErr {
				t.Errorf("DeserializeUint16() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(*tt.args.dest, tt.want) {
				t.Errorf("DeserializeUint16() = %v, want %v", *tt.args.dest, tt.want)
			}
		})
	}
}

func TestDeserializeUint32(t *testing.T) {
	h := uint32(0)
	type args struct {
		dest *uint32
		d    []byte
	}
	tests := []struct {
		name    string
		args    args
		want    uint32
		wantErr bool
	}{
		{"255", args{&h, []byte{0xff}}, 255, false},
		{"0", args{&h, []byte{}}, 0, false},
		{"invalid padding", args{&h, []byte{0x00, 0xff}}, 0, true},
		{"invalid data", args{&h, []byte("cats and dogs living together")}, 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := DeserializeUint32(tt.args.dest, tt.args.d); (err != nil) != tt.wantErr {
				t.Errorf("DeserializeUint32() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(*tt.args.dest, tt.want) {
				t.Errorf("DeserializeUint32() = %v, want %v", *tt.args.dest, tt.want)
			}
		})
	}
}

func TestDeserializeUint64(t *testing.T) {
	h := uint64(0)
	type args struct {
		dest *uint64
		d    []byte
	}
	tests := []struct {
		name    string
		args    args
		want    uint64
		wantErr bool
	}{
		{"255", args{&h, []byte{0xff}}, 255, false},
		{"0", args{&h, []byte{}}, 0, false},
		{"invalid padding", args{&h, []byte{0x00, 0xff}}, 0, true},
		{"invalid data", args{&h, []byte("cats and dogs living together")}, 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := DeserializeUint64(tt.args.dest, tt.args.d); (err != nil) != tt.wantErr {
				t.Errorf("DeserializeUint64() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(*tt.args.dest, tt.want) {
				t.Errorf("DeserializeUint64() = %v, want %v", *tt.args.dest, tt.want)
			}
		})
	}
}

func TestDeserializeInt(t *testing.T) {
	h := int(0)
	type args struct {
		dest *int
		d    []byte
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{"255", args{&h, []byte{0xfe, 0x03}}, 255, false},
		{"-255", args{&h, []byte{0xfd, 0x03}}, -255, false},
		{"zero", args{&h, []byte{}}, 0, false},
		{"invalid data", args{&h, []byte("🙁🙁🙁🙁🙁🙁🙁🙁🙁🙁🙁🙁🙁🙁🙁🙁🙁🙁🙁🙁🙁🙁🙁🙁🙁🙁🙁🙁🙁🙁🙁🙁🙁🙁🙁🙁🙁🙁🙁🙁🙁🙁🙁")}, 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := DeserializeInt(tt.args.dest, tt.args.d); (err != nil) != tt.wantErr {
				t.Errorf("DeserializeInt() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(*tt.args.dest, tt.want) {
				t.Errorf("DeserializeInt() = %v, want %v", *tt.args.dest, tt.want)
			}
		})
	}
}

func TestDeserializeInt8(t *testing.T) {
	h := int8(0)
	type args struct {
		dest *int8
		d    []byte
	}
	tests := []struct {
		name    string
		args    args
		want    int8
		wantErr bool
	}{
		{"127", args{&h, []byte{0xfe, 0x01}}, 127, false},
		{"-128", args{&h, []byte{0xff, 0x01}}, -128, false},
		{"zero", args{&h, []byte{}}, 0, false},
		{"invalid data", args{&h, []byte("🙁🙁🙁🙁🙁🙁🙁🙁🙁🙁🙁🙁🙁🙁🙁🙁🙁🙁🙁🙁🙁🙁🙁🙁🙁🙁🙁🙁🙁🙁🙁🙁🙁🙁🙁🙁🙁🙁🙁🙁🙁🙁🙁")}, 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := DeserializeInt8(tt.args.dest, tt.args.d); (err != nil) != tt.wantErr {
				t.Errorf("DeserializeInt8() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(*tt.args.dest, tt.want) {
				t.Errorf("DeserializeInt8() = %v, want %v", *tt.args.dest, tt.want)
			}
		})
	}
}

func TestDeserializeInt16(t *testing.T) {
	h := int16(0)
	type args struct {
		dest *int16
		d    []byte
	}
	tests := []struct {
		name    string
		args    args
		want    int16
		wantErr bool
	}{
		{"255", args{&h, []byte{0xfe, 0x03}}, 255, false},
		{"-255", args{&h, []byte{0xfd, 0x03}}, -255, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := DeserializeInt16(tt.args.dest, tt.args.d); (err != nil) != tt.wantErr {
				t.Errorf("DeserializeInt16() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(*tt.args.dest, tt.want) {
				t.Errorf("DeserializeInt16() = %v, want %v", *tt.args.dest, tt.want)
			}
		})
	}
}

func TestDeserializeInt32(t *testing.T) {
	h := int32(0)
	type args struct {
		dest *int32
		d    []byte
	}
	tests := []struct {
		name    string
		args    args
		want    int32
		wantErr bool
	}{
		{"255", args{&h, []byte{0xfe, 0x03}}, 255, false},
		{"-255", args{&h, []byte{0xfd, 0x03}}, -255, false},
		{"zero", args{&h, []byte{}}, 0, false},
		{"invalid data", args{&h, []byte("🙁🙁🙁🙁🙁🙁🙁🙁🙁🙁🙁🙁🙁🙁🙁🙁🙁🙁🙁🙁🙁🙁🙁🙁🙁🙁🙁🙁🙁🙁🙁🙁🙁🙁🙁🙁🙁🙁🙁🙁🙁🙁🙁")}, 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := DeserializeInt32(tt.args.dest, tt.args.d); (err != nil) != tt.wantErr {
				t.Errorf("DeserializeInt32() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(*tt.args.dest, tt.want) {
				t.Errorf("DeserializeInt32() = %v, want %v", *tt.args.dest, tt.want)
			}
		})
	}
}

func TestDeserializeInt64(t *testing.T) {
	h := int64(0)
	type args struct {
		dest *int64
		d    []byte
	}
	tests := []struct {
		name    string
		args    args
		want    int64
		wantErr bool
	}{
		{"255", args{&h, []byte{0xfe, 0x03}}, 255, false},
		{"-255", args{&h, []byte{0xfd, 0x03}}, -255, false},
		{"zero", args{&h, []byte{}}, 0, false},
		{"invalid data", args{&h, []byte("🙁🙁🙁🙁🙁🙁🙁🙁🙁🙁🙁🙁🙁🙁🙁🙁🙁🙁🙁🙁🙁🙁🙁🙁🙁🙁🙁🙁🙁🙁🙁🙁🙁🙁🙁🙁🙁🙁🙁🙁🙁🙁🙁")}, 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := DeserializeInt64(tt.args.dest, tt.args.d); (err != nil) != tt.wantErr {
				t.Errorf("DeserializeInt64() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(*tt.args.dest, tt.want) {
				t.Errorf("DeserializeInt64() = %v, want %v", *tt.args.dest, tt.want)
			}
		})
	}
}

func TestDeserializeFloat32(t *testing.T) {
	h := float32(0)
	type args struct {
		dest *float32
		d    []byte
	}
	tests := []struct {
		name    string
		args    args
		want    float32
		wantErr bool
	}{
		{"255.5", args{&h, []byte("255.5")}, 255.5, false},
		{"-255.5", args{&h, []byte("-255.5")}, -255.5, false},
		{"invalid data", args{&h, []byte("bob")}, 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := DeserializeFloat32(tt.args.dest, tt.args.d); (err != nil) != tt.wantErr {
				t.Errorf("DeserializeFloat32() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(*tt.args.dest, tt.want) {
				t.Errorf("DeserializeFloat32() = %v, want %v", *tt.args.dest, tt.want)
			}
		})
	}
}

func TestDeserializeFloat64(t *testing.T) {
	h := float64(0)
	type args struct {
		dest *float64
		d    []byte
	}
	tests := []struct {
		name    string
		args    args
		want    float64
		wantErr bool
	}{
		{"255.5", args{&h, []byte("255.5")}, 255.5, false},
		{"-255.5", args{&h, []byte("-255.5")}, -255.5, false},
		{"invalid data", args{&h, []byte("bob")}, 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := DeserializeFloat64(tt.args.dest, tt.args.d); (err != nil) != tt.wantErr {
				t.Errorf("DeserializeFloat64() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(*tt.args.dest, tt.want) {
				t.Errorf("DeserializeFloat64() = %v, want %v", *tt.args.dest, tt.want)
			}
		})
	}
}

func TestDeserializeBytes(t *testing.T) {
	h := []byte{}
	type args struct {
		dest *[]byte
		d    []byte
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{"\x66\x66\x66", args{&h, []byte{0x66, 0x66, 0x66}}, []byte{0x66, 0x66, 0x66}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := DeserializeBytes(tt.args.dest, tt.args.d); (err != nil) != tt.wantErr {
				t.Errorf("DeserializeBytes() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(*tt.args.dest, tt.want) {
				t.Errorf("DeserializeBytes() = %v, want %v", *tt.args.dest, tt.want)
			}
		})
	}
}

func TestDeserializeString(t *testing.T) {
	h := ""
	type args struct {
		dest *string
		d    []byte
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"my name is inigo montoya", args{&h, []byte("my name is inigo montoya")}, "my name is inigo montoya", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := DeserializeString(tt.args.dest, tt.args.d); (err != nil) != tt.wantErr {
				t.Errorf("DeserializeString() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(*tt.args.dest, tt.want) {
				t.Errorf("DeserializeString() = %v, want %v", *tt.args.dest, tt.want)
			}
		})
	}
}

func TestDeserializeBool(t *testing.T) {
	h := false
	type args struct {
		dest *bool
		d    []byte
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{"true", args{&h, []byte{1}}, true, false},
		{"false", args{&h, []byte{0}}, false, false},
		{"invalid", args{&h, []byte{0xff}}, false, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := DeserializeBool(tt.args.dest, tt.args.d); (err != nil) != tt.wantErr {
				t.Errorf("DeserializeBool() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(*tt.args.dest, tt.want) {
				t.Errorf("DeserializeBool() = %v, want %v", *tt.args.dest, tt.want)
			}
		})
	}
}

func TestDeserializeSlice(t *testing.T) {
	h := []string{}
	h2 := []uint{}
	h3 := []*uint{}
	hi1 := uint(255)
	hi2 := uint(1024)
	type args struct {
		dest interface{}
		data interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		{"string slice", args{&h, []interface{}{[]byte("my name is"), []byte("inigo montoya")}}, []string{"my name is", "inigo montoya"}, false},
		{"uint slice", args{&h2, []interface{}{[]byte{0xff}, []byte{0x04, 0x00}}}, []uint{255, 1024}, false},
		{"uint ptr slice", args{&h3, []interface{}{[]byte{0xff}, []byte{0x04, 0x00}}}, []*uint{&hi1, &hi2}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := DeserializeSlice(tt.args.dest, tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("DeserializeSlice() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(reflect.Indirect(reflect.ValueOf(tt.args.dest)).Interface(), tt.want) {
				t.Errorf("DeserializeSlice() = %v, want %v", reflect.Indirect(reflect.ValueOf(tt.args.dest)), tt.want)
			}
		})
	}
}

func TestDeserializeMap(t *testing.T) {
	h := make(map[string]uint)
	h2 := make(map[string]*uint)
	hi1 := uint(1)
	hi2 := uint(2)
	type args struct {
		dest interface{}
		data interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		{"string uint map", args{&h, []interface{}{[]interface{}{[]byte("one"), []byte{0x01}}, []interface{}{[]byte("two"), []byte{0x02}}}}, map[string]uint{"one": 1, "two": 2}, false},
		{"string uint map", args{&h2, []interface{}{[]interface{}{[]byte("one"), []byte{0x01}}, []interface{}{[]byte("two"), []byte{0x02}}}}, map[string]*uint{"one": &hi1, "two": &hi2}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := DeserializeMap(tt.args.dest, tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("DeserializeMap() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(reflect.Indirect(reflect.ValueOf(tt.args.dest)).Interface(), tt.want) {
				t.Errorf("DeserializeMap() = %v, want %v", reflect.Indirect(reflect.ValueOf(tt.args.dest)), tt.want)
			}
		})
	}
}

func TestDeserializeStruct(t *testing.T) {
	type h struct {
		Name string
		Age  uint
	}
	type h2 struct {
		Name string
		Age  uint
		nono bool
	}
	type h3 struct {
		Name string
		Age  uint
		Nono bool `rlp:"-"`
	}
	type h4 struct {
		Name   string
		Age    uint
		Yesyes *bool
	}
	hi := true
	type args struct {
		dest interface{}
		data interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		{"only public fields", args{&h{}, []interface{}{[]interface{}{[]byte("Name"), []byte("Bob")}, []interface{}{[]byte("Age"), []byte{0x96}}}}, h{"Bob", 150}, false},
		{"only public fields", args{&h4{}, []interface{}{[]interface{}{[]byte("Name"), []byte("Bob")}, []interface{}{[]byte("Age"), []byte{0x96}}, []interface{}{[]byte("Yesyes"), []byte{1}}}}, h4{"Bob", 150, &hi}, false},
		{"public and private fields", args{&h2{}, []interface{}{[]interface{}{[]byte("Name"), []byte("Bob")}, []interface{}{[]byte("Age"), []byte{0x96}}}}, h2{"Bob", 150, false}, false},
		{"tagged skip fields", args{&h3{}, []interface{}{[]interface{}{[]byte("Name"), []byte("Bob")}, []interface{}{[]byte("Age"), []byte{0x96}}}}, h3{"Bob", 150, false}, false},
		{"tagged skip fields invalid", args{&h3{}, []interface{}{[]interface{}{[]byte("Name"), []byte("Bob")}, []interface{}{[]byte("Age"), []byte{0x96}}, []interface{}{[]byte("Nono"), []byte{0}}}}, h3{"", 0, false}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := DeserializeStruct(tt.args.dest, tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("DeserializeStruct() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(reflect.Indirect(reflect.ValueOf(tt.args.dest)).Interface(), tt.want) {
				t.Errorf("DeserializeStruct() = %v, want %v", reflect.Indirect(reflect.ValueOf(tt.args.dest)), tt.want)
			}
		})
	}
}

func Test_deserialize(t *testing.T) {
	h := uint(0)
	h2 := uint8(0)
	h3 := uint16(0)
	h4 := uint32(0)
	h5 := uint64(0)
	h6 := int(0)
	h7 := int8(0)
	h8 := int16(0)
	h9 := int32(0)
	h10 := int64(0)
	h11 := float32(0)
	h12 := float64(0)
	h13 := []byte{}
	h14 := ""
	h15 := []string{}
	h16 := make(map[string]uint)
	type h17 struct {
		Name string
		Age  uint
	}
	type args struct {
		dest interface{}
		data interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		// TODO: Add test cases for custom desearization and text unmarshaler
		{"uint", args{&h, []byte{0xff}}, 255, false},
		{"uint8", args{&h2, []byte{0xff}}, 255, false},
		{"uint16", args{&h3, []byte{0xff}}, 255, false},
		{"uint32", args{&h4, []byte{0xff}}, 255, false},
		{"uint64", args{&h5, []byte{0xff}}, 255, false},
		{"int", args{&h6, []byte{0xfe, 0x01}}, 255, false},
		{"int8", args{&h7, []byte{0xfe, 0x01}}, 255, false},
		{"int16", args{&h8, []byte{0xfe, 0x01}}, 255, false},
		{"int32", args{&h9, []byte{0xfe, 0x01}}, 255, false},
		{"int64", args{&h10, []byte{0xfe, 0x01}}, 255, false},
		{"float32", args{&h11, []byte("255.5")}, 255.5, false},
		{"float64", args{&h12, []byte("255.5")}, 255.5, false},
		{"[]byte", args{&h13, []byte{0x66, 0x66, 0x66}}, []byte{0x66, 0x66, 0x66}, false},
		{"string", args{&h14, []byte("my name is inigo montoya")}, "my name is inigo montoya", false},
		{"slice", args{&h15, []interface{}{[]byte("my name is"), []byte("inigo montoya")}}, []string{"my name is", "inigo montoya"}, false},
		{"map", args{&h16, []interface{}{[]interface{}{[]byte("one"), []byte{0x01}}, []interface{}{[]byte("two"), []byte{0x02}}}}, map[string]uint{"one": 1, "two": 2}, false},
		{"struct", args{&h17{}, []interface{}{[]interface{}{[]byte("Name"), []byte("Bob")}, []interface{}{[]byte("Age"), []byte{0x96}}}}, h17{"Bob", 150}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := deserialize(tt.args.dest, tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("deserialize() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
