package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

// Create a new proxy
func newProxy(targetHost string) *httputil.ReverseProxy {
	url, err := url.Parse(targetHost)
	if err != nil {
		return nil
	}
	// Create a proxy to convert the internal port of K8s to external IPs
	return httputil.NewSingleHostReverseProxy(url)
}

// Establish the proxy server
func establishProxy(){
	// Internal Ports defined by YAML
	JupyterUrl := os.Getenv("JUPYTER_NOTEBOOK") + ":8888"
	SonarQubeUrl  := os.Getenv("SONARQUBE") + ":9000"
	SparkUrl  := os.Getenv("SPARK") + ":8080"
	HadoopUrl := os.Getenv("HADOOP") + ":50070"

	go proxyServer(HadoopUrl,":6766")
	go proxyServer(SparkUrl,":6866")
	go proxyServer(JupyterUrl,":6966")
	go proxyServer(SonarQubeUrl,":6070")

}

// Start service
func proxyServer(proxyUrl string, addr string) {
	_ = http.ListenAndServe(addr, newProxy(proxyUrl))
}

// Send back URL to clients
func sendRequest(conn net.Conn, url string) {
	conn.Write([]byte(url))
	fmt.Printf("Send back URL: %s\n", url)
}

func main() {

	// Establish the socket
	addr,_:=net.ResolveTCPAddr("tcp4", ":6666")

	netListen, err := net.ListenTCP("tcp", addr)
	CheckError(err)
	defer netListen.Close()

	// Establish the proxy server
	establishProxy()
	for {
		conn, err := netListen.Accept()
		if err != nil {
			continue
		}

		// Handlle new connection
		handleConnection(conn)
	}
}

// Hard code of expose GCP internal IPs to external URLs
func handleConnection(conn net.Conn) {
	buffer := make([]byte, 2048)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			return
		}
		receiveData := string(buffer[:n])
		switch receiveData {
		// Format: "http:// <GCP_IPS> : <PORTS>"
		case "1":
			sendRequest(conn,"http://35.202.46.187:6766")
		case "2":
			sendRequest(conn,"http://35.202.46.187:6866")
		case "3":
			sendRequest(conn,"http://35.202.46.187:6966")
		case "4":
			sendRequest(conn,"http://35.202.46.187:6070")
		}
	}
}

func Log(v ...interface{}) {
	log.Println(v...)
}

func CheckError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}