# gochat
Simple and concurrent tcp chat server

## Requirements

* Go 1.7+

## Install it

For the moment it is not published yet, so try doing a go get to:

```bash
$ go get github.com/lfmexi/gochat
$ go install github.com/lfmexi/gochat
```
## Run it

If you have installed it correctly, just run:

```bash
$ gochat
```
It will start a tcp server on your machine, exposing the port 8888 by default.

You can change the port to be exposed with the environment variable `GOCHAT_PORT`.

## Begin to chat

Connect to the server using telnet or another tcp client (I will create the client eventually):

```bash
$ telnet localhost 8888
```

### Login with a username

In order to send messages, you will need to login into the server:

```bash
login lfmexi
# It will answer with
Welcome user lfmexi
```

### Send messages

Once you are logged in, just send messages with the following format:

```bash
message [toUsuer] [yourMessage]
```
Example:

```
message lfmexi Hi this is a message for me 
```
It will send the message to the corresponding user

```
lfmexi: Hi this is a message for me 
```

### Logout

For login out, just send the command:

```bash
logout 
# it will answer with
Goodbye lfmexi
```
## Author

* Luis Fernando Morales ([lfmexi](https://github.com/lfmexi))