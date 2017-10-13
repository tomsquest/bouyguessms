package main

import (
	"flag"
	"fmt"
	"github.com/tomsquest/bouyguessms"
	"os"
	"time"
)

// Filled by goreleaser
var version = "master"
var commit = "HEAD"
var date = "now"

func main() {
	login := flag.String("login", "", "Your Bouygues Telecom `login`")
	pass := flag.String("pass", "", "Your Bouygues Telecom `password`")
	msg := flag.String("msg", "", "The `message`, ex: \"Hello World\", truncated if more than 140 chars")
	to := flag.String("to", "", "The destination `phonenumbers`, with a maximum of 5 numbers, separated by a semicolon and starting by 06 or 07, ex: \"0601010101;0702020202\"")
	flag.Parse()

	if *login == "" || *pass == "" || *msg == "" || *to == "" {
		fmt.Fprintf(os.Stderr, "Error: login, password, message and to are required\n\n")
		flag.Usage()
		fmt.Fprintf(os.Stderr, "Using %s %s %s %s\n\n", os.Args[0], version, commit, date)
		os.Exit(1)
	}

	smsLeft, err := bouyguessms.SendSms(*login, *pass, *msg, *to)
	if err != nil {
		fmt.Printf("Unable to send SMS: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("SMS sent successfully at %s. SMS left: %d.\n", time.Now().Format(time.RFC3339), smsLeft)
}
