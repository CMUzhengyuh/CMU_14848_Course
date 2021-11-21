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

func NewProxy(targetHost string) *httputil.ReverseProxy {
	url, err := url.Parse(targetHost)
	if err != nil {
		return nil
	}

	return httputil.NewSingleHostReverseProxy(url)
}

func proxyServer(proxyUrl string, addr string) {
	//被代理的服务器host和port
	fmt.Println(proxyUrl)
	serveErr := http.ListenAndServe(addr, NewProxy(proxyUrl))
	if serveErr != nil {
		panic(serveErr)
	}
}

func establishProxy(){
	JupyterNotebook := os.Getenv("JUPYTER_NOTEBOOK")
	//JupyterNotebook:="http://localhost"
	JupyterUrl := JupyterNotebook + ":8888"
	//jupyterUrl := "https://www.baidu.com"

	SonarQube := os.Getenv("SONARQUBE")
	//Spark:="http://localhost"
	SonarQubeUrl := SonarQube + ":9000"

	Spark := os.Getenv("SPARK")
	//Spark:="http://localhost"
	SparkUrl := Spark + ":8080"

	Hadoop := os.Getenv("HADOOP")
	//Spark:="http://localhost"
	HadoopUrl := Hadoop + ":50070"

	go proxyServer(HadoopUrl,":6766")
	go proxyServer(SparkUrl,":6866")
	go proxyServer(JupyterUrl,":6966")
	go proxyServer(SonarQubeUrl,":6070")

}

func sendRequest(conn net.Conn, text string) {
	conn.Write([]byte(text))
	fmt.Printf("send over %s\n",text)
}

func main() {

	//建立socket，监听端口
	addr,_:=net.ResolveTCPAddr("tcp4", ":6666")

	netListen, err := net.ListenTCP("tcp", addr)
	CheckError(err)
	defer netListen.Close()
	establishProxy()

	Log("Waiting for clients")
	for {
		conn, err := netListen.Accept()
		if err != nil {
			continue
		}
		Log(conn.RemoteAddr().String(), " tcp connect success")
		handleConnection(conn)
	}
}
//处理连接
func handleConnection(conn net.Conn) {
	currentIP := "35.202.46.187"
	buffer := make([]byte, 2048)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			Log(conn.RemoteAddr().String(), " connection error: ", err)
			return
		}
		Log(conn.RemoteAddr().String(), "receive data string:\n", string(buffer[:n]))
		receiveData := string(buffer[:n])
		switch receiveData {
		case "1":
			sendRequest(conn,"http://34.134.47.101:8888")
		case "2":
			sendRequest(conn,"http://34.135.7.97:9000")
		case "3":
			sendRequest(conn,"http://34.71.164.18:8080")
		case "4":
			sendRequest(conn,"http://"+currentIP+":6070")
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