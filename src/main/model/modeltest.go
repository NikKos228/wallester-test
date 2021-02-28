package models

import (
	"testing"
)

func Test_parseDate(t *testing.T) {
	type args struct {
		date string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "2020-10-05T15:04:05Z",
			args: args{"2020-10-05T15:04:05Z"},
			want: "05.10.2020",
		},
		{
			name: "1996-04-06T00:00:00Z",
			args: args{"1996-04-06T00:00:00Z"},
			want: "06.04.1996",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := birthdate(tt.args.date); got != tt.want {
				t.Errorf("parseDate() = %v, want %v", got, tt.want)
			}
		})
	}
}
