package main

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
	"strings"
)

type Message struct {
	Command   string
	Arguments []string
}

func main() {
	conn, err := net.Dial("tcp", "localhost:4000")
	if err != nil {
		fmt.Println(err)
		return
	}

	joinMessage := Message{Command: "join", Arguments: []string{"NN", "superdupersecret"}}

	fmt.Fprint(conn, joinMessage.ToProtocolString())

	var decider Decider
	decider = LongPathDecider{}

	world := NewWorld(0, 0, 0)

	reader := bufio.NewReader(conn)
	for {
		line, err := reader.ReadString('\n')
		message := ReadMessage(line)
		if err != nil {
			fmt.Println(err)
			return
		}

		switch message.Command {
		case "motd":
			fmt.Println("# Message of the Day")
			fmt.Println(message.Arguments[0])
			fmt.Println()
		case "game":
			w := OptimisticParseInt(message.Arguments[0])
			h := OptimisticParseInt(message.Arguments[1])
			pid := OptimisticParseInt(message.Arguments[2])
			world = NewWorld(w, h, pid)
		case "player":
			pid := OptimisticParseInt(message.Arguments[0])
			world.RegisterPlayer(pid, message.Arguments[1])
		case "pos":
			pid := OptimisticParseInt(message.Arguments[0])
			x := OptimisticParseInt(message.Arguments[1])
			y := OptimisticParseInt(message.Arguments[2])
			world.UpdatePosition(pid, x, y)
		case "die":
			for _, pidString := range message.Arguments {
				pid := OptimisticParseInt(pidString)
				world.KillPlayer(pid)
			}
		case "win":
			fmt.Println()
			fmt.Println("🎉🎉🎉 We won! 🎉🎉🎉")
			fmt.Println()
		case "lose":
			fmt.Println()
			fmt.Println("💀 We lost! 💀")
			fmt.Println()
		case "tick":
			world.PrettyPrint()
			move := decider.DecideMove(world)
			moveMessage := Message{Command: "move", Arguments: []string{move}}
			fmt.Fprint(conn, moveMessage.ToProtocolString())
		case "error":
			fmt.Println("An error occured: ", message.Arguments[0])
			return
		default:
			fmt.Println("Unknown message: ", message)
		}
	}
}

func ReadMessage(s string) Message {
	s = strings.TrimSuffix(s, "\n")
	parts := strings.Split(s, "|")
	return Message{
		Command:   parts[0],
		Arguments: parts[1:],
	}
}

func (m Message) ToProtocolString() string {
	parts := []string{m.Command}
	parts = append(parts, m.Arguments...)
	return fmt.Sprintln(strings.Join(parts, "|"))
}

func OptimisticParseInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		fmt.Println("Could not parse integer: ", s)
		return -1
	}

	return i
}
