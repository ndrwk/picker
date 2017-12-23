package picker

import (
	"reflect"
	"testing"
)

func TestDHT22_ReadValue(t *testing.T) {
	type fields struct {
		Name        string
		Moisture    float32
		Temperature float32
		Address     []byte
	}
	tests := []struct {
		name   string
		fields fields
		want   []float32
	}{
		{
			name: "Proper values",
			fields: fields{
				Moisture:    59.3,
				Temperature: 25.5,
			},
			want: []float32{59.3, 25.5},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := DHT22{
				Name:        tt.fields.Name,
				Moisture:    tt.fields.Moisture,
				Temperature: tt.fields.Temperature,
				Address:     tt.fields.Address,
			}
			if got := s.readValues(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PressureSensor.ReadValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDHT22_UpdateValues(t *testing.T) {
	type fields struct {
		Name        string
		Moisture    float32
		Temperature float32
		Address     []byte
	}
	type args struct {
		mousture    float32
		temperature float32
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "Proper values",
			fields: fields{
				Moisture:    43.4,
				Temperature: 28.55,
			},
			args: args{
				mousture:    43.4,
				temperature: 28.55,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := DHT22{
				Moisture:    tt.fields.Moisture,
				Temperature: tt.fields.Temperature,
			}
			s.updateValues([]float32{tt.fields.Moisture, tt.fields.Temperature})
			if got := s.readValues(); !reflect.DeepEqual(got, []float32{tt.args.mousture, tt.args.temperature}) {
				t.Errorf("TempSensor.ReadValue() = %v, want %v, %v", got, tt.args.mousture, tt.args.temperature)
			}
		})
	}
}

func TestDHT22_ReadName(t *testing.T) {
	type fields struct {
		Name string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "Proper values",
			fields: fields{
				Name: "kjhvgcfxd",
			},
			want: "kjhvgcfxd",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := DHT22{
				Name: tt.fields.Name,
			}
			if got := s.readName(); got != tt.want {
				t.Errorf("PressureSensor.ReadName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDHT22_UpdateName(t *testing.T) {
	type fields struct {
		Name string
	}
	type args struct {
		name string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "Proper values",
			fields: fields{
				Name: "khvcrfxlbiyuc",
			},
			args: args{
				name: "khvcrfxlbiyuc",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := DHT22{
				Name: tt.fields.Name,
			}
			s.updateName(tt.args.name)
			if got := s.readName(); got != tt.args.name {
				t.Errorf("TempSensor.ReadValue() = %v, want %v", got, tt.args.name)
			}
		})
	}
}

func TestDHT22_ReadAddr(t *testing.T) {
	type fields struct {
		Address []byte
	}
	tests := []struct {
		name   string
		fields fields
		want   []byte
	}{
		{
			name: "Proper values",
			fields: fields{
				Address: []byte{127},
			},
			want: []byte{127},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := DHT22{
				Address: tt.fields.Address,
			}
			if got := s.readAddr(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PressureSensor.ReadAddr() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDHT22_SetAddr(t *testing.T) {
	type fields struct {
		Address []byte
	}
	type args struct {
		addr []byte
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "Proper values",
			fields: fields{
				Address: []byte{125, 248, 7, 2, 44, 127},
			},
			args: args{
				addr: []byte{125, 248, 7, 2, 44, 127},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := DHT22{
				Address: tt.fields.Address,
			}
			s.setAddr(tt.fields.Address)
			if got := s.readAddr(); !reflect.DeepEqual(got, tt.args.addr) {
				t.Errorf("TempSensor.ReadAddr() = %v, want %v", got, tt.args.addr)
			}
		})
	}
}
