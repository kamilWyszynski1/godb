package main

import (
	"bufio"
	"fmt"
	"godb"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

const commandChar = '.'

func main() {
	SetupCloseHandler()

	reader := bufio.NewReader(os.Stdin)
	t, err := godb.OpenDB("db.godb")
	if err != nil {
		panic(err)
	}
	for {
		fmt.Print("godb-> ")
		text, err := reader.ReadString('\n')
		text = strings.Replace(text, "\n", "", -1)
		if err != nil {
			panic(fmt.Sprintf("failed to read line, %s", err))
		}

		if text[0] == commandChar {
			status := godb.Parse(text)

			switch status {
			case godb.SuccessfulExit:
				os.Exit(0)
			case godb.UnrecognizedCommand:
				fmt.Println("unrecognized command")
				continue
			}
		}

		statement := godb.PrepareStatement(text)

		if statement.Status == godb.PrepareUnrecognizedStatement {
			fmt.Println("unrecognized statement")
		} else if statement.Status == godb.PrepareSuccess {
			switch statement.Type {
			case godb.StatementSelect, godb.StatementInsert:
				godb.ExecuteStatement(statement, t)
			}
		} else if statement.Status == godb.PrepareSyntaxError {
			fmt.Println("Syntax error. Could not parse statement.")

		}
	}
}

// SetupCloseHandler creates a 'listener' on a new goroutine which will notify the
// program if it receives an interrupt from the OS. We then handle this by calling
// our clean up procedure and exiting the program.
func SetupCloseHandler() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		os.Exit(0)
	}()
}
