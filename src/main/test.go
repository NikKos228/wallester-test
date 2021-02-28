package main

import (
	_ "WallesterTest/vendor/github.com/jackc/pgx"
	"testing"
)

func Test_isBirthdayValid(t *testing.T) {
	type args struct {
		birthDay string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "In 18-60 age range.",
			args: args{"1987-11-15"},
			want: true,
		},
		{
			name: "Younger than 18",
			args: args{"2003-11-15"},
			want: false,
		},
		{
			name: "Older than 60",
			args: args{"1959-11-15"},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isAgeValid(tt.args.birthDay); got != tt.want {
				t.Errorf("isBirthdayValid() = %v, want %v", got, tt.want)
			}
		})
	}
}
