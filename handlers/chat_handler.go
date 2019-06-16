package handlers

import (
	"log"

	"net"

	"bufio"

	"github.com/lfmexi/gochat/commandparser"
)

// ChatHandler implements the server.Handler interface, also has
// methods for message broking.
type ChatHandler struct {
	broker *chatBroker
}

func (c *ChatHandler) logoutClient(clientAddr net.Addr) {
	log.Printf("User disconnection from %s\n", clientAddr)
	c.broker.dyingClients <- clientAddr
}

func (c *ChatHandler) logoutUser(conn net.Conn, w *bufio.Writer) {
	c.broker.dyingUsers <- &loggedOutUser{w, conn.RemoteAddr()}
	c.logoutClient(conn.RemoteAddr())
	return
}

func (c *ChatHandler) loginUser(user string, conn net.Conn, w *bufio.Writer) {
	messageChannel := make(chan []byte)

	c.broker.incomingUsers <- &loggedUser{user, w, messageChannel}
	c.broker.incomingClients <- &incomingClient{conn.RemoteAddr(), w, user}

	go func() {
		for m := range messageChannel {
			w.Write(m)
			w.Flush()
		}
	}()

	log.Printf("Loggin in user %s in from %s", user, conn.RemoteAddr())
}

func (c *ChatHandler) sendMessage(messageBody map[string]string, conn net.Conn, w *bufio.Writer) error {
	msg := message{conn.RemoteAddr(), w, messageBody["user"], messageBody["message"]}
	c.broker.messagePool <- msg
	return nil
}

// ServeTCP is the implementation of server.Handler.ServeTCP.
// Returns an error if there is something wrong with the incoming connection net.Conn.
func (c *ChatHandler) ServeTCP(conn net.Conn) error {
	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	cmd, _, err := reader.ReadLine()
	log.Printf("Getting message from %s", conn.RemoteAddr())
	
	if err != nil {
		if err.Error() == "EOF" {
			c.logoutClient(conn.RemoteAddr())
		}
		return err
	}

	if command, value, err := commandparser.ParseMessage(string(cmd)); err == nil {
		switch command {
		case commandparser.Login:
			c.loginUser(value.(string), conn, writer)
		case commandparser.Message:
			c.sendMessage(value.(map[string]string), conn, writer)
		case commandparser.Logout:
			c.logoutUser(conn, writer)
		}
		return nil
	}
	return nil
}

// NewChatHandler creates a ChatHandler value.
// Returns the pointer of that value.
func NewChatHandler() *ChatHandler {
	handler := &ChatHandler{
		newChatBroker(),
	}
	handler.broker.StartBroker()
	return handler
}
