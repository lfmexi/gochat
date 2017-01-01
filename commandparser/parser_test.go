package commandparser

import "testing"
import "fmt"

func TestParseLogin(t *testing.T) {
	message := "login user"
	if ct, v, err := ParseMessage(message); err == nil {
		if ct != Login {
			t.Fail()
		}
		if v.(string) != "user" {
			t.Fail()
		}
		return
	}
	t.Fail()
}

func TestParseMessage(t *testing.T) {
	message := "message testuser this is a test message"
	if ct, v, err := ParseMessage(message); err == nil {
		if ct != Message {
			t.Fail()
		}

		messageMap := v.(map[string]string)
		if messageMap["user"] != "testuser" {
			fmt.Printf("Unexpected value: %s\n", messageMap["user"])
			t.Fail()
		}
		if messageMap["message"] != "this is a test message\n" {
			fmt.Printf("Unexpected message value: %s\n", messageMap["message"])
			t.Fail()
		}
		return
	}
	t.Fail()
}

func TestLogout(t *testing.T) {
	message := "logout"
	if ct, v, err := ParseMessage(message); err == nil {
		if ct != Logout {
			t.Fail()
		}
		if v != nil {
			t.Fail()
		}
		return
	}
	t.Fail()
}

func TestInvalidMessage(t *testing.T) {

	message := "This is an invalid message"
	_, _, err := ParseMessage(message)

	if err == nil {
		t.Fail()
	}

	if err.Error() != "Parse error with the input \"This is an invalid message\"" {
		t.Fail()
	}

	return
}
