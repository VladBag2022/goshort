package misc

import (
	"net/url"
	"testing"

	log "github.com/sirupsen/logrus"
)

func BenchmarkSign(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Sign("123", "456")
	}
}

func BenchmarkVerify(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if _, _, err := Verify("123", "345|48cecbcac0ebb1d8bc0c395d5cc742c8f0eaf5e59696a8e0462ffa75990781df"); err != nil {
			log.Error(err)
		}
	}
}

func TestSign(t *testing.T) {
	type args struct {
		key   string
		value string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "positive test",
			args: args{
				key:   "123",
				value: "345",
			},
			want: "345|48cecbcac0ebb1d8bc0c395d5cc742c8f0eaf5e59696a8e0462ffa75990781df",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Sign(tt.args.key, tt.args.value); got != tt.want {
				t.Errorf("Sign() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVerify(t *testing.T) {
	type args struct {
		key string
		msg string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		want1   string
		wantErr bool
	}{
		{
			name: "positive test",
			args: args{
				key: "123",
				msg: "345|48cecbcac0ebb1d8bc0c395d5cc742c8f0eaf5e59696a8e0462ffa75990781df",
			},
			want:    true,
			want1:   "345",
			wantErr: false,
		},
		{
			name: "negative test - no separator",
			args: args{
				key: "123",
				msg: "48cecbcac0ebb1d8bc0c395d5cc742c8f0eaf5e59696a8e0462ffa75990781df",
			},
			want:    false,
			want1:   "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := Verify(tt.args.key, tt.args.msg)
			if (err != nil) != tt.wantErr {
				t.Errorf("Verify() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Verify() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("Verify() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestShorten(t *testing.T) {
	testURL, err := url.Parse("https://google.com")
	if err != nil {
		t.Error(err)
		return
	}

	type args struct {
		in0 *url.URL
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{
			name: "positive test",
			args: args{
				in0: testURL,
			},
			want:    10,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Shorten(tt.args.in0)
			if (err != nil) != tt.wantErr {
				t.Errorf("Shorten() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != tt.want {
				t.Errorf("Shorten() got = %v, want %v", len(got), tt.want)
			}
		})
	}
}
