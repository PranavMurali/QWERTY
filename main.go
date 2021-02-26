package main

import (
	"bufio"
	crand "crypto/rand"
	"encoding/binary"
	"errors"
	"fmt"
	"log"
	rand "math/rand"
	"os"
	"os/exec"
	"os/user"
	"strconv"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	user,err:=user.Current()
	tmps:=user.Username
	if err != nil {
		panic(err)
	}
	cmd:=exec.Command("toilet",tmps)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Run()

	for {
		fmt.Print("🔥🐲> ")
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}

		if err = execInput(input); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}
}

type cryptoSource struct{}

func (s cryptoSource) Seed(seed int64) {}

func (s cryptoSource) Int63() int64 {
        return int64(s.Uint64() & ^uint64(1<<63))
}

func (s cryptoSource) Uint64() (v uint64) {
        err := binary.Read(crand.Reader, binary.BigEndian, &v)
        if err != nil {
                log.Fatal(err)
        }
        return v
}

func execInput(input string) error {

	input = strings.TrimSuffix(input, "\n")

	args := strings.Split(input, " ")

	switch args[0] {

	case "cd":

		if len(args) < 2 {
			return errors.New("path required")
		}

		return os.Chdir(args[1])

	case "vscode":

		if len(args) < 2 {
			return errors.New("path required")
		}

		cmd := exec.Command("code", args[1:]...)
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout

		return cmd.Run()

	case "nano":

		if len(args) < 2 {
			return errors.New("path required")
		}

		cmd := exec.Command("nano", args[1:]...)
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout

		return cmd.Run()

	case "touch":

		if len(args) < 2 {
			return errors.New("path required")
		}

		cmd := exec.Command("touch", args[1:]...)
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout

		return cmd.Run()

	case "golang":

		if len(args) < 2 {
			return errors.New("path required")
		}

		cmd := exec.Command("go", args[1:]...)
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout

		return cmd.Run()

	case "userinfo":
		user, err := user.Current()
		if err != nil {
			panic(err)
		}
		fmt.Println("Hi " + user.Name + " (id: " + user.Uid + ")")
		fmt.Println("Username: " + user.Username)
		fmt.Println("Home Dir: " + user.HomeDir)
		fmt.Println("Real User: " + os.Getenv("SUDO_USER"))
		return nil

	case "wther":
		tmp:="http://wttr.in/chennai"
		cmd := exec.Command("curl", tmp)
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout

		return cmd.Run()
	
	case "art":
		var src cryptoSource
		rnd := rand.New(src)
		ndate:=rnd.Intn(30)
		date:=strconv.Itoa(ndate)
		if ndate<10{
			date="0"+date	
		}
		fmt.Print(date)
		tmp:="http://samiare.net/daily/1901"+date+"?width=20"
		cmd := exec.Command("curl", tmp)
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout

		return cmd.Run()

	case "ospref":
		if len(args) < 2 {
			return errors.New("name required")
		}
	
	case "exit":
		os.Exit(0)

	}

	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	return cmd.Run()
}
