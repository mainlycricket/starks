package byte_utils_test

import (
	"bytes"
	"testing"

	byte_utils "github.com/mainlycricket/starks/pkg/byte-utils"
)

func TestParseInt32Bytes(t *testing.T) {
	tests := []struct {
		name    string
		d       []byte
		want    int32
		wantErr bool
	}{
		{
			name:    "simple",
			d:       []byte{0, 0, 0, 23},
			want:    23,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := byte_utils.ParseInt32Bytes(tt.d)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("ParseInt32Bytes() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("ParseInt32Bytes() succeeded unexpectedly")
			}
			if got != tt.want {
				t.Errorf("ParseInt32Bytes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseInt16Bytes(t *testing.T) {
	tests := []struct {
		name    string
		d       []byte
		want    int16
		wantErr bool
	}{
		{
			name:    "simple",
			d:       []byte{1, 0},
			want:    256,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := byte_utils.ParseInt16Bytes(tt.d)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("ParseInt16Bytes() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("ParseInt16Bytes() succeeded unexpectedly")
			}
			if got != tt.want {
				t.Errorf("ParseInt16Bytes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetIntBytesInt16(t *testing.T) {
	got, err := byte_utils.GetIntBytes(int16(258))
	if err != nil {
		t.Fatal(err)
	}

	want := []byte{0x01, 0x02}

	if !bytes.Equal(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestGetIntBytesInt32(t *testing.T) {
	got, err := byte_utils.GetIntBytes(int32(258))
	if err != nil {
		t.Fatal(err)
	}

	want := []byte{0x00, 0x00, 0x01, 0x02}

	if !bytes.Equal(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}
