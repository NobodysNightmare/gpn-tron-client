package main

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
	"strings"
	"sync"
)

// Message encapsulates a message defined by the protocol
// at https://github.com/freehuntx/gpn-tron/blob/master/PROTOCOL.md
type Message struct {
	Command   string
	Arguments []string
}

func main() {
	var wg sync.WaitGroup

	PlayAsync(&wg, "NN-collision", "superdupersecure", CollisionAvoidDecider{}, false)
	PlayAsync(&wg, "NN-path", "superdupersecure", LongPathDecider{}, false)
	PlayAsync(&wg, "NN-score", "superdupersecure", HighScoreDecider{}, true)

	wg.Wait()
}

func PlayAsync(wg *sync.WaitGroup, name, secret string, decider Decider, draw bool) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		Play(name, secret, decider, draw)
	}()
}

// Play connects to a server and uses the given decider to make its moves.
func Play(name, secret string, decider Decider, draw bool) {
	conn, err := net.Dial("tcp", "localhost:4000")
	if err != nil {
		fmt.Println(err)
		return
	}

	joinMessage := Message{Command: "join", Arguments: []string{name, secret}}

	fmt.Fprint(conn, joinMessage.ToProtocolString())

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
			world.UpdatePosition(pid, Position{x, y})
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
			if draw {
				world.PrettyPrint()
			}

			move := decider.DecideMove(world)

			if draw {
				fmt.Println("Will move", move)
			}

			moveMessage := Message{Command: "move", Arguments: []string{string(move)}}
			fmt.Fprint(conn, moveMessage.ToProtocolString())
		case "error":
			fmt.Println("An error occured: ", message.Arguments[0])
			return
		default:
			fmt.Println("Unknown message: ", message)
		}
	}
}

// ReadMessage parses a single line returned by the server and returns
// the parsed Message.
func ReadMessage(s string) Message {
	s = strings.TrimSuffix(s, "\n")
	parts := strings.Split(s, "|")
	return Message{
		Command:   parts[0],
		Arguments: parts[1:],
	}
}

// ToProtocolString formats the Message in a way that it can be sent
// to the server.
func (m Message) ToProtocolString() string {
	parts := []string{m.Command}
	parts = append(parts, m.Arguments...)
	return fmt.Sprintln(strings.Join(parts, "|"))
}

// Parses the string, assuming that it can be parsed as an integer.
// Returns -1 if a parsing error occured or the parsed string was the
// number -1.
func OptimisticParseInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		fmt.Println("Could not parse integer: ", s)
		return -1
	}

	return i
}
