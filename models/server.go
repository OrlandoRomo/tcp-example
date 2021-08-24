package models

import (
	"errors"
	"fmt"
	"log"
	"net"
	"strings"
)

type Server struct {
	Rooms    map[string]*Room
	Commands chan Command
}

func NewServer() *Server {
	return &Server{
		Rooms:    make(map[string]*Room),
		Commands: make(chan Command),
	}
}

func (s *Server) Run() {
	for cmd := range s.Commands {
		switch cmd.Id {
		case CMD_NICK:
			s.Nick(cmd.Client, cmd.Args)
		case CMD_JOIN:
			s.Join(cmd.Client, cmd.Args)
		case CMD_ROOMS:
			s.ListRooms(cmd.Client, cmd.Args)
		case CMD_MSG:
			s.Msg(cmd.Client, cmd.Args)
		case CMD_QUIT:
			s.Quit(cmd.Client)
		}
	}
}

func (s *Server) NewClient(conn net.Conn) {
	log.Printf("new client connected: %s", conn.RemoteAddr().String())
	c := &Client{
		Conn:     conn,
		Nick:     "alv",
		Commands: s.Commands,
	}

	c.ReadInput()
}

func (s *Server) Nick(c *Client, args []string) {
	c.Nick = args[1]
	c.Msg(fmt.Sprintf("Welcome %s", c.Nick))
}

func (s *Server) Join(c *Client, args []string) {
	room := args[1]
	r, ok := s.Rooms[room]
	if !ok {
		r = &Room{
			Name:    room,
			Members: make(map[net.Addr]*Client),
		}
		s.Rooms[room] = r
	}
	r.Members[c.Conn.RemoteAddr()] = c

	s.QuitCurrentRoom(c)

	c.Room = r
	r.Broadcast(c, fmt.Sprintf("%s has joined the room", c.Nick))
	c.Msg(fmt.Sprintf("Welcome prro %s", c.Nick))
}
func (s *Server) ListRooms(c *Client, args []string) {
	var rooms []string
	for name := range s.Rooms {
		rooms = append(rooms, name)
	}
	c.Msg(fmt.Sprintf("List of rooms %s", strings.Join(rooms, ", ")))
}

func (s *Server) Msg(c *Client, args []string) {
	if c.Room == nil {
		c.Err(errors.New("You have to join a room first"))
		return
	}
	c.Room.Broadcast(c, c.Nick+strings.Join(args[1:len(args)], " "))
}

func (s *Server) Quit(c *Client) {
	log.Printf("client has disconnected %s", c.Conn.RemoteAddr().String())

	s.QuitCurrentRoom(c)

	c.Conn.Close()

}

func (s *Server) QuitCurrentRoom(c *Client) {
	if c.Room != nil {
		oldRoom := s.Rooms[c.Room.Name]
		delete(s.Rooms[c.Room.Name].Members, c.Conn.RemoteAddr())
		oldRoom.Broadcast(c, fmt.Sprintf("%s has left the room", c.Nick))

	}
}
