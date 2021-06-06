package main

import (
	"bufio"
	"bytes"
	crand "crypto/rand"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	rand "math/rand"
	"net/http"
	"os"
	"os/signal"
	"os/exec"
	"os/user"
	"regexp"
	"strconv"
	"strings"
    "time"
	"github.com/common-nighthawk/go-figure"
	"github.com/fatih/color"
	"github.com/jedib0t/go-pretty/table"
  "github.com/howeyc/gopass"
)
	

/*
Package dependency
sudo apt install bitwise
sudo npm i -g mdlt
sudo apt install toilet
*/

var history []string
var htime []string
type profile struct {
	Name     string
	colour string
	pref  int
}


func main() {

	reader := bufio.NewReader(os.Stdin)
	user, err := user.Current()
	tmps := user.Username
	if err != nil {
		panic(err)
	}
	//yellow := color.New(color.FgYellow).SprintFunc()
	//red := color.New(color.FgRed).SprintFunc()
	figure.NewFigure("Welcome "+tmps, "basic", true).Scroll(4000, 300, "left")
	cmd2 := exec.Command("clear")
	cmd2.Stderr = os.Stderr
	cmd2.Stdout = os.Stdout
	cmd2.Run()
	temps := [5]string{"QWERTY", "-F", "metal", "-f", "smblock"}
	cmd := exec.Command("toilet", temps[0:]...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Run()
	fmt.Print("\n")
	for {
		path, _ := os.Getwd()
		color.Red(path)
		fmt.Print("📀>> ")
		input, err := reader.ReadString('\n')
		history = append(history, input)
		dt := time.Now()
		dtf := dt.Format("01-02-2006 15:04:05")
		htime = append(htime, dtf)
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

	green := color.New(color.Bold, color.FgGreen).SprintFunc()
	red := color.New(color.Bold, color.FgRed).SprintFunc()
	blue := color.New(color.Bold, color.FgBlue).SprintFunc()
	cyan := color.New(color.Bold, color.FgCyan).SprintFunc()
	input = strings.TrimSuffix(input, "\n")
	args := strings.Split(input, " ")
	killSignal := make(chan os.Signal, 1)
	signal.Notify(killSignal, os.Interrupt)

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
		cmd.Stdin = os.Stdin

		return cmd.Run()

	case "dock-stat":
		cmd := exec.Command("dockly")
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout
		cmd.Stdin=os.Stdin

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

	case "supercow":

		if len(args) < 2 {
			return errors.New("User Name Required")
		}
		fmt.Println("Enter Password")
		pass,_ := gopass.GetPasswd()
		fmt.Println(pass)

	case "userinfo":
		user, err := user.Current()
		if err != nil {
			panic(err)
		}

		fmt.Println(green("Hi ") + user.Name + cyan(" (id: "+user.Uid+")"))
		fmt.Println(green("Username: ") + red(user.Username))
		fmt.Println(green("Home Dir: ") + blue(user.HomeDir))
		fmt.Println(green("Real User: ") + red(os.Getenv("SUDO_USER")))
		return nil

	case "wther":
		tmp := "http://wttr.in/chennai"
		cmd := exec.Command("curl", tmp)
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout
		return cmd.Run()

	case "cpuinfo":
		cmd := exec.Command("./cpufetch")
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout
		return cmd.Run()

	case "covid":
		if len(args) < 2 {
			return errors.New("Country name required")
		}
		response, err := http.Get("https://coronavirus-19-api.herokuapp.com/countries/" + args[1])

		if err != nil {
			fmt.Print(err.Error())
			os.Exit(1)
		}

		responseData, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Fatal(err)
		}
		var r map[string]interface{}
		json.Unmarshal([]byte(responseData), &r)
		if len(args) == 2 {
			fmt.Printf("Country:       " + blue(args[1]+"\n"))
			fmt.Printf("Cases:         %.f\n", r["cases"])
			fmt.Printf("Deaths:        %.f\n", r["deaths"])
			fmt.Printf("Recovered:     %.f\n", r["recovered"])
			fmt.Printf("Active Cases:  %.f\n", r["active"])
			fmt.Printf("Deaths Today:  %.f\n", r["todayDeaths"])
			fmt.Printf("Cases Today:   %.f\n", r["todayCases"])
		} else {
			fmt.Printf("Country:                     " + blue(args[1]+"\n"))
			fmt.Printf("Cases:                       %.f\n", r["cases"])
			fmt.Printf("Deaths:                      %.f\n", r["deaths"])
			fmt.Printf("Recovered:                   %.f\n", r["recovered"])
			fmt.Printf("Active Cases:                %.f\n", r["active"])
			fmt.Printf("Deaths Today:                %.f\n", r["todayDeaths"])
			fmt.Printf("Cases Today:                 %.f\n", r["todayCases"])
			fmt.Printf("Cases per 1 Million:         %.f\n", r["casesPerOneMillion"])
			fmt.Printf("Deaths per 1 Million:        %.f\n", r["deathsPerOneMillion"])
			fmt.Printf("Total Tests:                 %.f\n", r["totalTests"])
			fmt.Printf("Tests per 1 Million:         %.f\n", r["testsPerOneMillion"])
		}
		return nil

	case "art":
		var src cryptoSource
		rnd := rand.New(src)
		ndate := rnd.Intn(30)
		date := strconv.Itoa(ndate)
		if ndate < 10 {
			date = "0" + date
		}
		tmp := "http://samiare.net/daily/1901" + date + "?width=20"
		cmd := exec.Command("curl", tmp)
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout

		return cmd.Run()

	case "ospref":
		if len(args) < 2 {
			return errors.New("name required")
		}

	case "ls":
		var sarr []string
		// var sarr1 []string
		isAlpha := regexp.MustCompile(`^[A-Za-z]+$`).MatchString
		cmd := exec.Command("ls", "-lsh")
		var outb, errb bytes.Buffer
		cmd.Stderr = &errb
		cmd.Stdout = &outb
		cmd.Run()
		outs := outb.String()
		sarr = strings.Split(outs, "\n")
		i := 0
		fmt.Print("\n")

		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)
		t.AppendHeader(table.Row{"File Name", "Size"})
		for i < len(sarr) {
			stemp := strings.Split(strings.TrimLeft(sarr[i], " "), " ")
			if len(stemp) > 2 {
				if isAlpha(stemp[len(stemp)-5]) {
					t.AppendRows([]table.Row{
						{stemp[len(stemp)-1], stemp[len(stemp)-6]},
					})
				} else {
					t.AppendRows([]table.Row{
						{stemp[len(stemp)-1], stemp[len(stemp)-5]},
					})
				}

			}
			i++
			if i == len(sarr)-1 {
				t.Render()
			} else if len(sarr) == 1 {
				t.Render()
			}
		}

		return nil

	case "history":
		i := 0
		for i < len(history) {
			fmt.Printf(strings.TrimSpace(htime[i]))
			fmt.Print("               ")
			fmt.Printf(strings.TrimSpace(history[i]))
			fmt.Printf("\n")
			i = i + 1
		}

		return nil

	case "c++":
		cmd := exec.Command("g++", args[1])
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout
		cmd.Run()
		time.Sleep(10)
		cmd2 := exec.Command("./a.out")
		cmd2.Stderr = os.Stderr
		cmd2.Stdout = os.Stdout
		cmd2.Run()
		return nil
	
	case "pkgman":
		
		cmd := exec.Command("/bin/sh", "-c", "sudo apt-get install mps-youtube", args[1])
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout
		cmd.Stdin = os.Stdin
		return cmd.Run()

	case "youplayer":

		cmd := exec.Command("mpsyt")
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout
		cmd.Stdin = os.Stdin
		return cmd.Run()

	case "init-project":
		if args[1] == "options" {
			fmt.Println("vanilla\nvanilla-ts\nvue\nvue-ts\nreact\nreact-ts\npreact\npreact-ts\nlit-element\nlit-element-ts\nsvelte\nsvelte-ts")
		} else {
			cmd := exec.Command("npm", "init", "@vitejs/app", "myproject", "--template", args[1])
			cmd.Stderr = os.Stderr
			cmd.Stdout = os.Stdout
			cmd.Run()
			fmt.Println("\n\n Project Initialised \n\n")
		}

		return nil

	case "convert":
		cmd := exec.Command("bitwise", args[1])
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout
		cmd.Run()
		return nil

	case "calculate":
		cmd := exec.Command("mdlt", args[1], args[2])
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout
		cmd.Run()
		return nil

	case "python-details":
		cmd := exec.Command("python", "--version")
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout
		cmd.Run()
		fmt.Print("\n")
		cmd2 := exec.Command("python3", "--version")
		cmd2.Stderr = os.Stderr
		cmd2.Stdout = os.Stdout
		cmd2.Run()
		fmt.Print("\n")
		cmd3 := exec.Command("pip", "freeze")
		var outb, errb bytes.Buffer
		cmd3.Stderr = &errb
		cmd3.Stdout = &outb
		cmd3.Run()
		outs := outb.String()
		var sarr []string
		sarr = strings.Split(outs, "\n")
		i := 0
		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)
		t.AppendHeader(table.Row{"Package", "Version"})
		for i < len(sarr)-1 {
			stemp := strings.Split(sarr[i], "==")
			t.AppendRows([]table.Row{
				{stemp[0], stemp[1]},
			})
			i++

			if i == len(sarr)-2 {
				t.Render()
			}
		}

		return nil

	case "exit":
		os.Exit(0)
	
	default:
		cmd := exec.Command(args[0], args[1:]...)
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout
		return cmd.Run()
	}
	<-killSignal
	return nil
	
}

