package config

import "testing"

func TestGetTargetURL(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "It should return the correct hardcoded URL",
			want: "https://www.google.com",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetTargetURL(); got != tt.want {
				t.Errorf("GetTargetURL() = %v, want %v", got, tt.want)
			}
		})
	}
}
