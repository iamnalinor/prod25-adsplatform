package floatutil

import "testing"

func TestNorm(t *testing.T) {
	type args struct {
		target float64
		other  float64
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{
			name: "both zero",
			args: args{target: 0, other: 0},
			want: 0,
		},
		{
			name: "target smaller",
			args: args{target: 33, other: 100},
			want: 0.33,
		},
		{
			name: "target larger",
			args: args{target: 100, other: 33},
			want: 1,
		},
		{
			name: "negative target smaller",
			args: args{target: -33, other: 100},
			want: -0.33,
		},
		{
			name: "negative target larger",
			args: args{target: -100, other: 33},
			want: -1,
		},
		{
			name: "both negative",
			args: args{target: -50, other: -100},
			want: -0.5,
		},
		{
			name: "target zero",
			args: args{target: 0, other: 75},
			want: 0,
		},
		{
			name: "other zero",
			args: args{target: 75, other: 0},
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Norm(tt.args.target, tt.args.other); got != tt.want {
				t.Errorf("Norm() = %v, want %v", got, tt.want)
			}
		})
	}
}
