package picker

import (
	"reflect"
	"testing"
)

func TestBMP085_ReadValue(t *testing.T) {
	type fields struct {
		Name    string
		Pressure   float32
		Temperature float32
		Address []byte
	}
	tests := []struct {
		name   string
		fields fields
		want   []float32
	}{
		{
			name: "Proper values",
			fields: fields{
				Pressure: 760,
				Temperature: 33.1,
			},
			want: []float32{760, 33.1},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := BMP085{
				Pressure: tt.fields.Pressure,
				Temperature:  tt.fields.Temperature,
			}
			if got := s.readValues(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PressureSensor.ReadValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBMP085_UpdateValues(t *testing.T) {
	type fields struct {
		Name        string
		Pressure    float32
		Temperature float32
		Address     []byte
	}
	type args struct {
		pressure    float32
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
				Pressure:    760.0,
				Temperature: 28.55,
			},
			args: args{
				pressure:    760.0,
				temperature: 28.55,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := BMP085{
				Name:        tt.fields.Name,
				Pressure:    tt.fields.Pressure,
				Temperature: tt.fields.Temperature,
				Address:     tt.fields.Address,
			}
			s.updateValues([]float32{tt.args.pressure, tt.args.temperature})
			if got := s.readValues(); !reflect.DeepEqual(got, []float32{tt.args.pressure, tt.args.temperature}) {
				t.Errorf("TempSensor.ReadValue() = %v, want %v, %v", got, tt.args.pressure, tt.args.temperature)
			}
		})
	}
}

func TestBMP085_ReadName(t *testing.T) {
	type fields struct {
		Name    string
		Value   float32
		Address []byte
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
			s := BMP085{
				Name:     tt.fields.Name,
				Pressure: tt.fields.Value,
				Address:  tt.fields.Address,
			}
			if got := s.readName(); got != tt.want {
				t.Errorf("PressureSensor.ReadName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBMP085_UpdateName(t *testing.T) {
	type fields struct {
		Name    string
		Value   float32
		Address []byte
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
			s := BMP085{
				Name:     tt.fields.Name,
				Pressure: tt.fields.Value,
				Address:  tt.fields.Address,
			}
			s.updateName(tt.args.name)
			if got := s.readName(); got != tt.args.name {
				t.Errorf("TempSensor.ReadValue() = %v, want %v", got, tt.args.name)
			}
		})
	}
}

func TestBMP085_ReadAddr(t *testing.T) {
	type fields struct {
		Name    string
		Value   float32
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
			s := BMP085{
				Name:     tt.fields.Name,
				Pressure: tt.fields.Value,
				Address:  tt.fields.Address,
			}
			if got := s.readAddr(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PressureSensor.ReadAddr() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBMP085_SetAddr(t *testing.T) {
	type fields struct {
		Name    string
		Value   float32
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
			s := BMP085{
				Name:     tt.fields.Name,
				Pressure: tt.fields.Value,
				Address:  tt.fields.Address,
			}
			s.setAddr(tt.fields.Address)
			if got := s.readAddr(); !reflect.DeepEqual(got, tt.args.addr) {
				t.Errorf("TempSensor.ReadAddr() = %v, want %v", got, tt.args.addr)
			}
		})
	}
}
