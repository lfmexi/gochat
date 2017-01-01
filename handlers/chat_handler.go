package handlers

import (
	"fmt"
	"log"

	"net"

	"bufio"

	"github.com/lfmexi/gochat/commandparser"
)

type userHandler struct {
	username     string
	messageQueue chan []byte
}

type incomingClient struct {
	address  net.Addr
	username string
}

// ChatHandler implements the server.Handler interface, also has
// methods for message broking.
type ChatHandler struct {
	incomingUsers   chan *userHandler
	users           map[string]*userHandler
	dyingUsers      chan string
	incomingClients chan *incomingClient
	clients         map[net.Addr]string
	dyingClients    chan net.Addr
	messagePool     chan map[string]string
}

// StartBroker creates a go routine for listening all the handler channels in c.
func (c *ChatHandler) StartBroker() {
	go func() {
		for {
			select {
			case client := <-c.incomingClients:
				c.clients[client.address] = client.username
			case u := <-c.incomingUsers:
				c.users[u.username] = u
			case dyingClient := <-c.dyingClients:
				delete(c.clients, dyingClient)
			case dyingUser := <-c.dyingUsers:
				user := c.users[dyingUser]
				delete(c.users, dyingUser)
				close(user.messageQueue)
			case message := <-c.messagePool:
				dest := message["user"]
				rawmessage := []byte(message["message"])
				if user := c.users[dest]; user != nil {
					user.messageQueue <- rawmessage
				}
			}
		}
	}()
}

func (c *ChatHandler) logoutClient(clientAddr net.Addr) {
	user := c.clients[clientAddr]
	if user != "" {
		log.Printf("User %s disconnected from %s\n", user, clientAddr)
		c.dyingUsers <- user
		c.dyingClients <- clientAddr
	}
}

func (c *ChatHandler) logoutUser(conn net.Conn, w *bufio.Writer) error {
	goodbyeMessage := fmt.Sprintf("Goodbye %s\n", c.clients[conn.RemoteAddr()])
	c.logoutClient(conn.RemoteAddr())
	if _, err := w.Write([]byte(goodbyeMessage)); err != nil {
		return err
	}
	return w.Flush()
}

func (c *ChatHandler) loginUser(user string, conn net.Conn, writer *bufio.Writer) error {

	onlineuser := c.users[user]

	if onlineuser != nil {
		log.Printf("User %s already logged in", user)
		alreadyMessage := "You are already online\n"
		if _, err := writer.Write([]byte(alreadyMessage)); err != nil {
			return err
		}

		return writer.Flush()
	}

	messageChannel := make(chan []byte)

	c.incomingUsers <- &userHandler{user, messageChannel}
	c.incomingClients <- &incomingClient{conn.RemoteAddr(), user}

	go func() {
		for m := range messageChannel {
			writer.Write(m)
			writer.Flush()
		}
	}()

	log.Printf("User %s logged in from %s", user, conn.RemoteAddr())

	welcome := fmt.Sprintf("Welcome user %s\n", user)
	if _, err := writer.Write([]byte(welcome)); err != nil {
		return err
	}

	return writer.Flush()
}

func (c *ChatHandler) sendMessage(messageBody map[string]string, conn net.Conn, w *bufio.Writer) error {
	if sender := c.clients[conn.RemoteAddr()]; sender == "" {
		uoffline := "You are offline my friend, please login and try again\n"
		if _, err := w.Write([]byte(uoffline)); err != nil {
			return err
		}
		return w.Flush()
	}
	if dest := c.users[messageBody["user"]]; dest == nil {
		useroffline := fmt.Sprintf("User %s is offline, try again later\n", messageBody["user"])
		if _, err := w.Write([]byte(useroffline)); err != nil {
			return err
		}
		return w.Flush()
	}
	messageBody["message"] = fmt.Sprintf("%s: %s", c.clients[conn.RemoteAddr()], messageBody["message"])
	c.messagePool <- messageBody
	return nil
}

// ServeTCP is the implementation of server.Handler.ServeTCP.
// Returns an error if there is something wrong with the incoming connection net.Conn.
func (c *ChatHandler) ServeTCP(conn net.Conn) error {
	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	cmd, _, err := reader.ReadLine()

	if err != nil {
		if err.Error() == "EOF" {
			c.logoutClient(conn.RemoteAddr())
		}
		return err
	}

	if command, value, err := commandparser.ParseMessage(string(cmd)); err == nil {
		switch command {
		case commandparser.Login:
			return c.loginUser(value.(string), conn, writer)
		case commandparser.Message:
			return c.sendMessage(value.(map[string]string), conn, writer)
		case commandparser.Logout:
			return c.logoutUser(conn, writer)
		}
	}
	return nil
}

// NewChatHandler creates a ChatHandler value.
// Returns the pointer of that value.
func NewChatHandler() *ChatHandler {
	return &ChatHandler{
		make(chan *userHandler),
		make(map[string]*userHandler),
		make(chan string),
		make(chan *incomingClient),
		make(map[net.Addr]string),
		make(chan net.Addr),
		make(chan map[string]string),
	}
}
