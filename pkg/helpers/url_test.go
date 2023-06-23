package helpers

import "testing"

func TestIsURL(t *testing.T) {
	type args struct {
		input string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Valid uri",
			args: args{input: "http://google.com"},
			want: true,
		},
		{
			name: "Valid uri",
			args: args{input: "http://ozon.ru"},
			want: true,
		},
		{
			name: "Non valid uri",
			args: args{input: "http:/google.com"},
			want: false,
		},
		{
			name: "Non valid uri",
			args: args{input: "google.com"},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsURL(tt.args.input); got != tt.want {
				t.Errorf("IsURL() = %v, want %v", got, tt.want)
			}
		})
	}
}
