package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/exec"
	"runtime"
	"strings"
)


func sendRequest(conn net.Conn, text string) {
	conn.Write([]byte(text))
	fmt.Printf("send over %s\n",text)
}

func receiveRespond(conn net.Conn) string{
	buffer := make([]byte, 2048)
	n, err := conn.Read(buffer)
	if err != nil {
		fmt.Println(conn.RemoteAddr().String(), " connection error: ", err)
		return ""
	}
	fmt.Println(conn.RemoteAddr().String(), "receive data string:\n", string(buffer[:n]))
	return string(buffer[:n])
}

// open opens the specified URL in the default browser of the user.
func open(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c","start"}
	case "darwin":
		cmd = "open"
	default: // "linux", "freebsd", "openbsd", "netbsd"
		cmd = "xdg-open"
	}
	args = append(args, url)
	return exec.Command(cmd, args...).Start()
}

func terminal(conn net.Conn) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Simple Shell")
	for {
		fmt.Println("Welcome to Big Data Processing Application")
		fmt.Println("Please type the number that corresponds to which application you would like to run:")
		fmt.Println("1. Apache Hadoop")
		fmt.Println("2. Apache Spark")
		fmt.Println("3. Jupyter Notebook")
		fmt.Println("4. SonarQube and SonarScanner")
		fmt.Println("Type the number here >")

		text, _ := reader.ReadString('\n')
		// convert CRLF to LF
		text = strings.Replace(text, "\n", "", -1)
		sendRequest(conn,text)
		website := receiveRespond(conn)
		open(website)
	}
}

func main() {
	server := "127.0.0.1:6666"
	tcpAddr, err := net.ResolveTCPAddr("tcp4", server)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}

	fmt.Println("connect success")
	terminal(conn)

}