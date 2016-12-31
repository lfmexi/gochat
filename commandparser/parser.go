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

func isLogoutCommand(command string) bool {
	validLogout := regexp.MustCompile(`logout`)
	if validLogout.MatchString(command) {
		return true
	}
	return false
}

func isLoginCommand(command string) (bool, string) {
	validLogin := regexp.MustCompile(`^login [a-z]([a-z]|[0-9])+$`)
	if validLogin.MatchString(command) {
		values := strings.Split(command, " ")
		return true, values[1]
	}
	return false, ""
}

func isMessageCommand(command string) (bool, map[string]string) {
	validMessage := regexp.MustCompile(`message [a-z]([a-z]|[0-9])+ .`)
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
func ParseMessage(message string) (CommandType, interface{}) {

	if ok, value := isLoginCommand(message); ok {
		return Login, value
	}
	if ok, value := isMessageCommand(message); ok {
		return Message, value
	}
	if ok := isLogoutCommand(message); ok {
		return Logout, nil
	}
	return "", nil
}
