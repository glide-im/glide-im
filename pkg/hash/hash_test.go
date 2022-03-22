package hash

import "testing"

func TestHash(t *testing.T) {
	type args struct {
		data []byte
		seed uint32
	}
	tests := []struct {
		name string
		args args
		want uint32
	}{
		{
			name: "",
			args: struct {
				data []byte
				seed uint32
			}{
				data: []byte("ABCD"),
				seed: 1,
			},
			want: 3027734286,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Hash(tt.args.data, tt.args.seed); got != tt.want {
				t.Errorf("Hash() = %v, want %v", got, tt.want)
			}
		})
	}
}
