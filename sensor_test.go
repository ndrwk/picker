package picker

import (
	"reflect"
	"testing"
)

func TestTempSensor_ReadValue(t *testing.T) {
	type fields struct {
		Name    string
		Value   float32
		Address []byte
	}
	tests := []struct {
		name   string
		fields fields
		want   float32
	}{
		{
			name: "Proper values",
			fields: fields {
				Value: 25.88,
			},
			want: 25.88,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := TempSensor{
				Name:    tt.fields.Name,
				Value:   tt.fields.Value,
				Address: tt.fields.Address,
			}
			if got := s.ReadValue(); got != tt.want {
				t.Errorf("TempSensor.ReadValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTempSensor_UpdateValue(t *testing.T) {
	type fields struct {
		Name    string
		Value   float32
		Address []byte
	}
	type args struct {
		value float32
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := TempSensor{
				Name:    tt.fields.Name,
				Value:   tt.fields.Value,
				Address: tt.fields.Address,
			}
			s.UpdateValue(tt.args.value)
		})
	}
}

func TestTempSensor_ReadName(t *testing.T) {
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
			fields: fields {
				Name: "kjhvgcfxd",
			},
			want: "kjhvgcfxd",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := TempSensor{
				Name:    tt.fields.Name,
				Value:   tt.fields.Value,
				Address: tt.fields.Address,
			}
			if got := s.ReadName(); got != tt.want {
				t.Errorf("TempSensor.ReadName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTempSensor_UpdateName(t *testing.T) {
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
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := TempSensor{
				Name:    tt.fields.Name,
				Value:   tt.fields.Value,
				Address: tt.fields.Address,
			}
			s.UpdateName(tt.args.name)
		})
	}
}

func TestTempSensor_ReadAddr(t *testing.T) {
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
			fields: fields {
				Address: []byte{125, 248, 7, 2, 44, 127},
			},
			want: []byte{125, 248, 7, 2, 44, 127},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := TempSensor{
				Name:    tt.fields.Name,
				Value:   tt.fields.Value,
				Address: tt.fields.Address,
			}
			if got := s.ReadAddr(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TempSensor.ReadAddr() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTempSensor_SetAddr(t *testing.T) {
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
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := TempSensor{
				Name:    tt.fields.Name,
				Value:   tt.fields.Value,
				Address: tt.fields.Address,
			}
			s.SetAddr(tt.args.addr)
		})
	}
}

func TestPressureSensor_ReadValue(t *testing.T) {
	type fields struct {
		Name    string
		Value   float32
		Address []byte
	}
	tests := []struct {
		name   string
		fields fields
		want   float32
	}{
		{
			name: "Proper values",
			fields: fields {
				Value: 760,
			},
			want: 760,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := PressureSensor{
				Name:    tt.fields.Name,
				Value:   tt.fields.Value,
				Address: tt.fields.Address,
			}
			if got := s.ReadValue(); got != tt.want {
				t.Errorf("PressureSensor.ReadValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPressureSensor_UpdateValue(t *testing.T) {
	type fields struct {
		Name    string
		Value   float32
		Address []byte
	}
	type args struct {
		value float32
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := PressureSensor{
				Name:    tt.fields.Name,
				Value:   tt.fields.Value,
				Address: tt.fields.Address,
			}
			s.UpdateValue(tt.args.value)
		})
	}
}

func TestPressureSensor_ReadName(t *testing.T) {
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
			fields: fields {
				Name: "kjhvgcfxd",
			},
			want: "kjhvgcfxd",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := PressureSensor{
				Name:    tt.fields.Name,
				Value:   tt.fields.Value,
				Address: tt.fields.Address,
			}
			if got := s.ReadName(); got != tt.want {
				t.Errorf("PressureSensor.ReadName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPressureSensor_UpdateName(t *testing.T) {
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
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := PressureSensor{
				Name:    tt.fields.Name,
				Value:   tt.fields.Value,
				Address: tt.fields.Address,
			}
			s.UpdateName(tt.args.name)
		})
	}
}

func TestPressureSensor_ReadAddr(t *testing.T) {
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
			fields: fields {
				Address: []byte{127},
			},
			want: []byte{127},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := PressureSensor{
				Name:    tt.fields.Name,
				Value:   tt.fields.Value,
				Address: tt.fields.Address,
			}
			if got := s.ReadAddr(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PressureSensor.ReadAddr() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPressureSensor_SetAddr(t *testing.T) {
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
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := PressureSensor{
				Name:    tt.fields.Name,
				Value:   tt.fields.Value,
				Address: tt.fields.Address,
			}
			s.SetAddr(tt.args.addr)
		})
	}
}
