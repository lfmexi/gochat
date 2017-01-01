package commandparser

import (
	"fmt"
	"regexp"
	"strings"
)

// Login is the CommandType of the login commands.
// Logout is the CommandType of the logout commands.
// Message is the CommandType of the message commands.
const (
	Login   = "login"
	Logout  = "logout"
	Message = "message"
)

// CommandType could be one of the defined const strings.
type CommandType string

// CommandParseError is the custom error for these actions
type CommandParseError struct {
	in  string
	err error
}

func (c *CommandParseError) Error() string {
	return fmt.Sprintf("Parse error with the input \"%s\"", c.in)
}

func isLogoutCommand(command string) bool {
	validLogout := regexp.MustCompile(`logout`)
	if validLogout.MatchString(command) {
		return true
	}
	return false
}

func isLoginCommand(command string) (bool, string) {
	validLogin := regexp.MustCompile(`^login [a-z]([a-zA-Z]|[1-9])+$`)
	if validLogin.MatchString(command) {
		values := strings.Split(command, " ")
		return true, values[1]
	}
	return false, ""
}

func isMessageCommand(command string) (bool, map[string]string) {
	validMessage := regexp.MustCompile(`message [a-z]([a-zA-Z]|[0-9])+ .`)
	if validMessage.MatchString(command) {
		values := strings.Split(command, " ")
		var resultMap = make(map[string]string)
		resultMap["user"] = values[1]
		message := ""
		for i := 2; i < len(values); i++ {
			if i == 2 {
				message = values[i]
			} else {
				message = fmt.Sprintf("%s %s", message, values[i])
			}
		}
		resultMap["message"] = fmt.Sprintf("%s\n", message)
		return true, resultMap
	}
	return false, nil
}

// ParseMessage identifies the CommandType of the received message.
// Returns the CommandType and the corresponding value in an interface{}.
func ParseMessage(message string) (CommandType, interface{}, error) {

	if ok, value := isLoginCommand(message); ok {
		return Login, value, nil
	}
	if ok, value := isMessageCommand(message); ok {
		return Message, value, nil
	}
	if ok := isLogoutCommand(message); ok {
		return Logout, nil, nil
	}
	err := fmt.Errorf("The message given is invalid")
	return "", nil, &CommandParseError{message, err}
}
