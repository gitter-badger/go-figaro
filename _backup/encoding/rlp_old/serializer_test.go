package rlp_old

import (
	"reflect"
	"testing"
)

func TestSerializeUint(t *testing.T) {
	type args struct {
		d uint
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		{"255", args{255}, []byte{0xff}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SerializeUint(tt.args.d)
			if (err != nil) != tt.wantErr {
				t.Errorf("SerializeUint() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SerializeUint() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSerializeUint8(t *testing.T) {
	type args struct {
		d uint8
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		{"255", args{255}, []byte{0xff}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SerializeUint8(tt.args.d)
			if (err != nil) != tt.wantErr {
				t.Errorf("SerializeUint8() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SerializeUint8() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSerializeUint16(t *testing.T) {
	type args struct {
		d uint16
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		{"255", args{255}, []byte{0xff}, false},
		{"257", args{257}, []byte{0x01, 0x01}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SerializeUint16(tt.args.d)
			if (err != nil) != tt.wantErr {
				t.Errorf("SerializeUint16() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SerializeUint16() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSerializeUint32(t *testing.T) {
	type args struct {
		d uint32
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		{"255", args{255}, []byte{0xff}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SerializeUint32(tt.args.d)
			if (err != nil) != tt.wantErr {
				t.Errorf("SerializeUint32() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SerializeUint32() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSerializeUint64(t *testing.T) {
	type args struct {
		d uint64
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		{"255", args{255}, []byte{0xff}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SerializeUint64(tt.args.d)
			if (err != nil) != tt.wantErr {
				t.Errorf("SerializeUint64() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SerializeUint64() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSerializeInt(t *testing.T) {
	type args struct {
		d int
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		{"255", args{255}, []byte{0xfe, 0x03}, false},
		{"-255", args{-255}, []byte{0xfd, 0x03}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SerializeInt(tt.args.d)
			if (err != nil) != tt.wantErr {
				t.Errorf("SerializeInt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SerializeInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSerializeInt8(t *testing.T) {
	type args struct {
		d int8
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		{"127", args{127}, []byte{0xfe, 0x01}, false},
		{"-128", args{-128}, []byte{0xff, 0x01}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SerializeInt8(tt.args.d)
			if (err != nil) != tt.wantErr {
				t.Errorf("SerializeInt8() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SerializeInt8() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSerializeInt16(t *testing.T) {
	type args struct {
		d int16
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		{"255", args{255}, []byte{0xfe, 0x03}, false},
		{"-255", args{-255}, []byte{0xfd, 0x03}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SerializeInt16(tt.args.d)
			if (err != nil) != tt.wantErr {
				t.Errorf("SerializeInt16() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SerializeInt16() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSerializeInt32(t *testing.T) {
	type args struct {
		d int32
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		{"255", args{255}, []byte{0xfe, 0x03}, false},
		{"-255", args{-255}, []byte{0xfd, 0x03}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SerializeInt32(tt.args.d)
			if (err != nil) != tt.wantErr {
				t.Errorf("SerializeInt32() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SerializeInt32() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSerializeInt64(t *testing.T) {
	type args struct {
		d int64
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		{"255", args{255}, []byte{0xfe, 0x03}, false},
		{"-255", args{-255}, []byte{0xfd, 0x03}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SerializeInt64(tt.args.d)
			if (err != nil) != tt.wantErr {
				t.Errorf("SerializeInt64() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SerializeInt64() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSerializeFloat32(t *testing.T) {
	type args struct {
		d float32
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		{"255.5", args{255.5}, []byte("255.5"), false},
		{"-255.5", args{-255.5}, []byte("-255.5"), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SerializeFloat32(tt.args.d)
			if (err != nil) != tt.wantErr {
				t.Errorf("SerializeFloat32() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SerializeFloat32() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSerializeFloat64(t *testing.T) {
	type args struct {
		d float64
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		{"255.5", args{255.5}, []byte("255.5"), false},
		{"-255.5", args{-255.5}, []byte("-255.5"), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SerializeFloat64(tt.args.d)
			if (err != nil) != tt.wantErr {
				t.Errorf("SerializeFloat64() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SerializeFloat64() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSerializeBytes(t *testing.T) {
	type args struct {
		d []byte
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		{"\x66\x66\x66", args{[]byte{0x66, 0x66, 0x66}}, []byte{0x66, 0x66, 0x66}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SerializeBytes(tt.args.d)
			if (err != nil) != tt.wantErr {
				t.Errorf("SerializeBytes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SerializeBytes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSerializeString(t *testing.T) {
	type args struct {
		d string
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		{"my name is inigo montoya", args{"my name is inigo montoya"}, []byte("my name is inigo montoya"), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SerializeString(tt.args.d)
			if (err != nil) != tt.wantErr {
				t.Errorf("SerializeString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SerializeString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSerializeBool(t *testing.T) {
	type args struct {
		d bool
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		{"true", args{true}, []byte{1}, false},
		{"false", args{false}, []byte{0}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SerializeBool(tt.args.d)
			if (err != nil) != tt.wantErr {
				t.Errorf("SerializeBool() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SerializeBool() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSerializeSlice(t *testing.T) {
	type args struct {
		d interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		{"string slice", args{[]string{"my name is", "inigo montoya"}}, []interface{}{[]byte("my name is"), []byte("inigo montoya")}, false},
		{"uint slice", args{[]uint{255, 1024}}, []interface{}{[]byte{0xff}, []byte{0x04, 0x00}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SerializeSlice(tt.args.d)
			if (err != nil) != tt.wantErr {
				t.Errorf("SerializeSlice() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SerializeSlice() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSerializeMap(t *testing.T) {
	type args struct {
		d interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		{"string uint map", args{map[string]uint{"one": 1, "two": 2}}, []interface{}{[]interface{}{[]byte("one"), []byte{0x01}}, []interface{}{[]byte("two"), []byte{0x02}}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SerializeMap(tt.args.d)
			if (err != nil) != tt.wantErr {
				t.Errorf("SerializeMap() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SerializeMap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSerializeStruct(t *testing.T) {
	type args struct {
		d interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		{"only public fields", args{struct {
			Name string
			Age  uint
		}{"Bob", 150}}, []interface{}{[]interface{}{[]byte("Name"), []byte("Bob")}, []interface{}{[]byte("Age"), []byte{0x96}}}, false},
		{"public and private fields", args{struct {
			Name string
			Age  uint
			nono bool
		}{"Bob", 150, true}}, []interface{}{[]interface{}{[]byte("Name"), []byte("Bob")}, []interface{}{[]byte("Age"), []byte{0x96}}}, false},
		{"tagged skip fields", args{struct {
			Name string
			Age  uint
			Nono bool `rlp:"-"`
		}{"Bob", 150, true}}, []interface{}{[]interface{}{[]byte("Name"), []byte("Bob")}, []interface{}{[]byte("Age"), []byte{0x96}}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SerializeStruct(tt.args.d)
			if (err != nil) != tt.wantErr {
				t.Errorf("SerializeStruct() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SerializeStruct() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_serialize(t *testing.T) {
	type args struct {
		data interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		// TODO: add test cases for custom serialization and text marshaler
		{"uint", args{uint(255)}, []byte{0xff}, false},
		{"uint8", args{uint8(255)}, []byte{0xff}, false},
		{"uint16", args{uint16(255)}, []byte{0xff}, false},
		{"uint32", args{uint32(255)}, []byte{0xff}, false},
		{"uint64", args{uint64(255)}, []byte{0xff}, false},
		{"int", args{int(255)}, []byte{0xfe, 0x03}, false},
		{"int8", args{int8(127)}, []byte{0xfe, 0x01}, false},
		{"int16", args{int16(255)}, []byte{0xfe, 0x03}, false},
		{"int32", args{int32(255)}, []byte{0xfe, 0x03}, false},
		{"int64", args{int64(255)}, []byte{0xfe, 0x03}, false},
		{"float32", args{float32(255.5)}, []byte("255.5"), false},
		{"float64", args{float64(255.5)}, []byte("255.5"), false},
		{"[]byte", args{[]byte{0x66, 0x66, 0x66}}, []byte{0x66, 0x66, 0x66}, false},
		{"string", args{"my name is inigo montoya"}, []byte("my name is inigo montoya"), false},
		{"true", args{true}, []byte{1}, false},
		{"slice", args{[]string{"my name is", "inigo montoya"}}, []interface{}{[]byte("my name is"), []byte("inigo montoya")}, false},
		{"map", args{map[string]uint{"one": 1, "two": 2}}, []interface{}{[]interface{}{[]byte("one"), []byte{0x01}}, []interface{}{[]byte("two"), []byte{0x02}}}, false},
		{"struct", args{struct {
			Name string
			Age  uint
		}{"Bob", 150}}, []interface{}{[]interface{}{[]byte("Name"), []byte("Bob")}, []interface{}{[]byte("Age"), []byte{0x96}}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := serialize(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("serialize() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("serialize() = %v, want %v", got, tt.want)
			}
		})
	}
}
