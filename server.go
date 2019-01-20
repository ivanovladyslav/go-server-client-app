package main

import "fmt"
import "net"
import "log"
import "bufio"
import "strconv"
import "math/rand"
import "unicode/utf8"
import "unicode"
import "flag"

func handleServerConnection(c net.Conn, i int, clCount *int) {
    iteration := 0
    nextKey := ""
    var s Session_protector
        // scan message
    scanner := bufio.NewScanner(c)

    for scanner.Scan() {
        msg := scanner.Text()
        log.Printf("Client %v sends: %v", i, msg)
        fmt.Println("---")
        if iteration == 0 {
          s.__hash = msg
          c.Write([]byte(msg + "\n"))
          iteration = iteration + 1
        } else if iteration == 1 {
          nextKey = s.next_session_key(msg)
          c.Write([]byte(nextKey + "\n"))
          log.Printf("Client %v receives: %v", i, nextKey)
          iteration = iteration + 1
        } else if iteration < 7 {
          nextKey = s.next_session_key(msg)
          c.Write([]byte(nextKey + "\n"))
          log.Printf("Client %v receives: %v", i, nextKey)
          iteration = iteration + 1
        }
        if scanner.Text() == "end" {
          *clCount = *clCount - 1
          fmt.Println("client disconnected")
          c.Close()
        }
    }
}

func main() {
    argsPort := flag.String("port",":9998","a string") //":9999"
    argsMaxClients := flag.Int("n",2,"an int") //"10"
    flag.Parse()
    log.Println("Server launched...")
    // listen on all interfaces
    ln, _ := net.Listen("tcp", ":"+*argsPort)
    i := 0
    for {
        // accept connection on port
        if i < *argsMaxClients {
          c, _ := ln.Accept()
          i++
          log.Printf("Client %v connected...", i)
          fmt.Println("---")
          // handle the connection
          go handleServerConnection(c, i, &i)
        }
    }
}

type Session_protector struct {
	__hash string
}

func (s Session_protector) __calc_hash(session_key string, val int) int {
	var result = 0
	for i := 0; i < utf8.RuneCountInString(session_key); i++ {
		current, err := strconv.Atoi(string(session_key[i]))
		if err != nil {

		} else {
			if val == 1 {
				result = result + current + 41
			} else if val == 2 {
				result = result + current * 72
			} else if val == 3{
				result = result + current ^ 10
			} else {
				result = result + current * 35
			}
		}
	}
	return result
}

func (s Session_protector) next_session_key(session_key string) string {
	var result = 0
	if s.__hash == "" {
		fmt.Println("It's empty")
	} else {
		for _, element := range s.__hash {
			if !unicode.IsDigit(element) {
				fmt.Println("Its not digit!!")
			}
		}
		for i := 0; i < utf8.RuneCountInString(s.__hash); i++ {
			result += s.__calc_hash(session_key, i)
		}
	}
	return (strconv.Itoa(result))
}
