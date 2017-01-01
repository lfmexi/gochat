package handlers

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

type loggedUser struct {
	username     string
	userWriter   *bufio.Writer
	messageQueue chan []byte
}

type loggedOutUser struct {
	userWriter *bufio.Writer
	clientAddr net.Addr
}

type incomingClient struct {
	address      net.Addr
	clientWriter *bufio.Writer
	username     string
}

type message struct {
	sender       net.Addr
	senderWriter *bufio.Writer
	destination  string
	rawMessage   string
}

type chatBroker struct {
	// Channels of the broker
	incomingClients chan *incomingClient
	incomingUsers   chan *loggedUser
	dyingClients    chan net.Addr
	dyingUsers      chan *loggedOutUser
	messagePool     chan message

	// Internal storages of the broker
	clients map[net.Addr]string
	users   map[string]*loggedUser
}

func sendAlreadyConnMsg(w *bufio.Writer) {
	alreadyMessage := "You are already online\n"
	if _, err := w.Write([]byte(alreadyMessage)); err != nil {
		log.Println(err)
		return
	}
	w.Flush()
}

func notConnectedMsg(w *bufio.Writer) {
	uoffline := "You are offline my friend, please login and try again\n"
	if _, err := w.Write([]byte(uoffline)); err != nil {
		log.Println(err)
		return
	}
	w.Flush()
}

func loggedInMsg(w *bufio.Writer, u string) {
	welcome := fmt.Sprintf("Welcome user %s\n", u)
	if _, err := w.Write([]byte(welcome)); err != nil {
		log.Println(err)
		return
	}
	w.Flush()
}

func logOutMsg(w *bufio.Writer, u string) {
	goodbyeMessage := fmt.Sprintf("Goodbye %s\n", u)
	if _, err := w.Write([]byte(goodbyeMessage)); err != nil {
		log.Println(err)
		return
	}
	w.Flush()
}

func (c *chatBroker) sendMessageToDest(m *message) {
	if dest := c.users[m.destination]; dest != nil {
		msg := fmt.Sprintf("%s: %s", c.clients[m.sender], m.rawMessage)
		dest.messageQueue <- []byte(msg)
		return
	}
	useroffline := fmt.Sprintf("User %s is offline, try again later\n", m.destination)
	if _, err := m.senderWriter.Write([]byte(useroffline)); err != nil {
		log.Println(err)
		return
	}
	m.senderWriter.Flush()
}

// StartBroker creates a go routine for listening all the handler channels in c.
func (c *chatBroker) StartBroker() {
	go func() {
		for {
			select {
			case client := <-c.incomingClients:
				if regClient := c.clients[client.address]; regClient == "" {
					c.clients[client.address] = client.username
					break
				}
				sendAlreadyConnMsg(client.clientWriter)

			case u := <-c.incomingUsers:
				if loggedUser := c.users[u.username]; loggedUser == nil {
					c.users[u.username] = u
					loggedInMsg(u.userWriter, u.username)
					break
				}
				sendAlreadyConnMsg(u.userWriter)
			case dyingClient := <-c.dyingClients:
				username := c.clients[dyingClient]
				if user := c.users[username]; user != nil {
					delete(c.users, username)
					close(user.messageQueue)
				}
				delete(c.clients, dyingClient)
			case dyingUser := <-c.dyingUsers:
				username := c.clients[dyingUser.clientAddr]
				if user := c.users[username]; user != nil {
					delete(c.users, username)
					close(user.messageQueue)
					logOutMsg(dyingUser.userWriter, username)
				}
			case message := <-c.messagePool:
				sender := message.sender
				if senderUser := c.clients[sender]; senderUser != "" {
					c.sendMessageToDest(&message)
					break
				}
				notConnectedMsg(message.senderWriter)
			}
		}
	}()
}

func newChatBroker() *chatBroker {
	return &chatBroker{
		make(chan *incomingClient),
		make(chan *loggedUser),
		make(chan net.Addr),
		make(chan *loggedOutUser),
		make(chan message),
		make(map[net.Addr]string),
		make(map[string]*loggedUser),
	}
}
