package handlers

import (
	"bufio"
	"bytes"
	"testing"
)

func Test_sendAlreadyConnMsg(t *testing.T) {
	type args struct {
		w *bufio.Writer
	}

	var initialBuffer bytes.Buffer

	tests := []struct {
		name   string
		args   args
		wanted string
	}{
		{"an Already connected user", args{bufio.NewWriter(&initialBuffer)}, "You are already online\n"},
	}
	for _, tt := range tests {
		sendAlreadyConnMsg(tt.args.w)
		got := initialBuffer.String()
		if got != tt.wanted {
			t.Errorf("sendAlreadyConnMsg expected %v in the buffer and got %v", tt.wanted, got)
		}

		initialBuffer.Reset()
	}
}

func Test_notConnectedMsg(t *testing.T) {
	type args struct {
		w *bufio.Writer
	}

	var initialBuffer bytes.Buffer

	tests := []struct {
		name   string
		args   args
		wanted string
	}{
		{"not connected user", args{bufio.NewWriter(&initialBuffer)}, "You are offline my friend, please login and try again\n"},
	}
	for _, tt := range tests {
		notConnectedMsg(tt.args.w)

		got := initialBuffer.String()

		if got != tt.wanted {
			t.Errorf("notConnectedMsg expected %v in the buffer and got %v", tt.wanted, got)
		}

		initialBuffer.Reset()
	}
}

func Test_loggedInMsg(t *testing.T) {
	type args struct {
		w *bufio.Writer
		u string
	}

	var initialBuffer bytes.Buffer

	tests := []struct {
		name   string
		args   args
		wanted string
	}{
		{"first user", args{bufio.NewWriter(&initialBuffer), "luis"}, "Welcome user luis\n"},
		{"second user", args{bufio.NewWriter(&initialBuffer), "john"}, "Welcome user john\n"},
		{"third user", args{bufio.NewWriter(&initialBuffer), "doe"}, "Welcome user doe\n"},
	}

	for _, tt := range tests {
		loggedInMsg(tt.args.w, tt.args.u)

		got := initialBuffer.String()

		if got != tt.wanted {
			t.Errorf("loggedInMsg expected %v in the buffer and got %v", tt.wanted, got)
		}

		initialBuffer.Reset()
	}
}

func Test_logOutMsg(t *testing.T) {
	type args struct {
		w *bufio.Writer
		u string
	}
	var initialBuffer bytes.Buffer

	tests := []struct {
		name   string
		args   args
		wanted string
	}{
		{"first user", args{bufio.NewWriter(&initialBuffer), "luis"}, "Goodbye luis\n"},
		{"second user", args{bufio.NewWriter(&initialBuffer), "john"}, "Goodbye john\n"},
		{"third user", args{bufio.NewWriter(&initialBuffer), "doe"}, "Goodbye doe\n"},
	}

	for _, tt := range tests {
		logOutMsg(tt.args.w, tt.args.u)

		got := initialBuffer.String()

		if got != tt.wanted {
			t.Errorf("logOutMsg expected %v in the buffer and got %v", tt.wanted, got)
		}

		initialBuffer.Reset()
	}
}
