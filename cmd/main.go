package main

import (
	"bufio"
	"fmt"
	"godb"
	"os"
	"os/signal"
	"strings"
	"syscall"

	log "github.com/sirupsen/logrus"
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
				if err := t.Close(); err != nil {
					log.Errorf("failed to close table, %s", err)
				}
				return
			case godb.UnrecognizedCommand:
				log.Errorln("unrecognized command")
				continue
			}
		}

		statement := godb.PrepareStatement(text)

		if statement.Status == godb.PrepareUnrecognizedStatement {
			log.Errorln("unrecognized command")
		} else if statement.Status == godb.PrepareSuccess {
			switch statement.Type {
			case godb.StatementSelect, godb.StatementInsert:
				if err := godb.ExecuteStatement(statement, t); err != nil {
					log.Errorf("failed to execute statement, %s", err)
				}
			}
		} else if statement.Status == godb.PrepareSyntaxError {
			log.Errorln("Syntax error. Could not parse statement.")

		}
	}
}

// SetupCloseHandler creates a 'listener' on a new goroutine which will notify the
// program if it receives an interrupt from the OS. We then handle this by calling
// our clean up procedure and exiting the program.
func SetupCloseHandler() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-c
		log.Infof("received interrupt, exiting")
		os.Exit(0)
	}()
}
