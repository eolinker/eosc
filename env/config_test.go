package env

import "testing"

func Test_formatPath(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
		{
			name: "abs dir",
			args: args{path: "/abs_dir/"},
			want: "/abs_dir",
		},
		{
			name: "abs file",
			args: args{path: "/abs_dir"},
			want: "/abs_dir",
		},
		{
			name: "abs file 2",
			args: args{path: "/abs_dir/file.yml"},
			want: "/abs_dir/file.yml",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := formatPath(tt.args.path); got != tt.want {
				t.Errorf("formatPath() = %v, want %v", got, tt.want)
			}
		})
	}
}
