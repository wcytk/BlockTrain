package network

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/wcytk/trainThroughBlockchain/train"
	"github.com/wcytk/trainThroughBlockchain/transaction"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
)

type Server struct {
}

var clientIPs []string

var clientAccounts []string

var hostIP string

var isTraining = false

var fromAddr string

var passphrase string

func NewServer() *Server {
	server := &Server{}
	server.setRoute()
	return server
}

func (server *Server) Start() {
	if err := http.ListenAndServe(":12345", nil); err != nil {
		fmt.Println(err)
		return
	}
}

func (server *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World!"))
}

func (server *Server) setRoute() {
	http.HandleFunc("/", server.ServeHTTP)
	http.HandleFunc("/req", server.req)
	http.HandleFunc("/addIP", server.addIP)
	http.HandleFunc("/prepareTraining", server.prepareTraining)
	http.HandleFunc("/startTraining", server.startTraining)
	http.HandleFunc("/stopTraining", server.stopTraining)
	http.HandleFunc("/enterTraining", server.enterTraining)
}

func (server *Server) req(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("You have successfully sent your request!\n"))
	w.Write([]byte("Your ip is " + ClientIP(r) + "\n"))
	if ClientPublicIP(r) == "" {
		w.Write([]byte("You are communicating in local network!"))
	} else {
		w.Write([]byte("Your public ip is " + ClientPublicIP(r)))
	}
}

func (server *Server) addIP(w http.ResponseWriter, r *http.Request) {
	if isTraining == false {
		w.Write([]byte("The host haven't start a training!"))
	} else {
		_ = r.ParseForm()
		if r.Form["account"][0] != "" {
			if ClientPublicIP(r) == "" {
				w.Write([]byte("You are communicating in local network!\n"))
				w.Write([]byte("Your local ip is " + ClientIP(r) + "\n"))
				// 对clientIP进行查找，查看客户是否已经加入训练
				// Search in clientIP to determine whether client has entered the training
				hasEntered := sort.Search(len(clientIPs), func(hasEntered int) bool {
					return clientIPs[hasEntered] >= ClientIP(r)
				})

				if hasEntered < len(clientIPs) && ClientIP(r) == clientIPs[hasEntered] {
					w.Write([]byte("You have already entered the training!\n"))
				} else {
					clientIPs = append(clientIPs, ClientIP(r))
					clientAccounts = append(clientAccounts, r.Form["account"][0])
				}
			} else {
				w.Write([]byte("Your public ip is " + ClientPublicIP(r) + "\n"))

				hasEntered := sort.Search(len(clientIPs), func(hasEntered int) bool {
					return clientIPs[hasEntered] >= ClientPublicIP(r)
				})

				if hasEntered < len(clientIPs) && ClientPublicIP(r) == clientIPs[hasEntered] {
					w.Write([]byte("You have already entered the training!\n"))
				} else {
					clientIPs = append(clientIPs, ClientPublicIP(r))
				}
			}
		}
		// 对客户的所有的IP进行排序
		// Sort clientIP
		sort.Strings(clientIPs)
		w.Write([]byte("The IPs now in this training\n"))
		for i := 0; i < len(clientIPs); i++ {
			w.Write([]byte(clientIPs[i] + "\n"))
		}
	}
}

func (server *Server) prepareTraining(w http.ResponseWriter, r *http.Request) {
	host, _ := GetLocalPublicIp()
	if host == "" {
		w.Write([]byte("You are not connected to network!"))
	} else {
		w.Write([]byte("Your IP address is " + host + "\n"))
		w.Write([]byte("Now a training is hosting in this IP"))
		hostIP = host
		isTraining = true

		fmt.Print("Input your account: ")
		_, _ = fmt.Scanln(&fromAddr)

		fmt.Print("Input your passphrase: ")
		_, _ = fmt.Scanln(&passphrase)

	}
}

func (server *Server) startTraining(w http.ResponseWriter, r *http.Request) {
	if len(clientIPs) <= 0 {
		w.Write([]byte("No one has entered your training yet!"))
	} else {
		train.StartPYTraining(clientIPs, hostIP)
		postTrainingRequest(clientIPs, hostIP)
		for i := 0; i < len(clientIPs); i++ {
			transaction.StartTransaction(fromAddr, clientAccounts[i], passphrase)
		}
	}
}

func (server *Server) stopTraining(w http.ResponseWriter, r *http.Request) {
	if isTraining == false {
		w.Write([]byte("You haven't start a training"))
	} else {
		isTraining = false
		w.Write([]byte("You have stopped a training"))
	}
}

func postTrainingRequest(clientIPs []string, hostIP string) {
	mapMassage := make(map[string]interface{})
	mapMassage["clientIPs"] = clientIPs
	mapMassage["hostIP"] = hostIP


	for i := 0; i < len(clientIPs); i++ {
		ip := clientIPs[i]
		url := "http://" + ip + ":12345/enterTraining"

		mapMassage["postIP"] = clientIPs[i]
		mapMassage["index"] = i

		message, err := json.Marshal(mapMassage)

		if err != nil {
			log.Fatal("Map 2 Json error ", err)
		}

		req, err := http.NewRequest("POST", url, bytes.NewBuffer(message))
		if err != nil {
			log.Fatal("NewRequest error ", err)
		}
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Fatal("Client sent error ", err)
		}
		body, _ := ioutil.ReadAll(resp.Body)
		log.Println("Response body" + string(body) + "\n")
	}
}

func (server *Server) enterTraining(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatal("Parse form error ", err)
	}

	// 初始化请求变量的结构
	formData := make(map[string]interface{})

	// 调用go自带的json库对body进行解析
	err = json.NewDecoder(r.Body).Decode(&formData)
	if err != nil {
		log.Fatal("Json decode error ", err)
	}

	// 将interface类型转换为[]string类型
	// Convert interface type to []string type

	// 首先将interface类型转化为[]interface类型，因为在加入进json中的时候是以[]string形式，所以现在是[]interface
	// Firstly, convert interface type to []interface type
	tmp := formData["clientIPs"].([]interface{})
	hostClientIPs := make([]string, len(tmp))
	for i, v := range tmp {
		hostClientIPs[i] = v.(string)
	}

	remoteHostIP := formData["hostIP"].(string)

	clientIP := formData["postIP"].(string)

	index := int(formData["index"].(float64))

	hasEntered := sort.Search(len(hostClientIPs), func(hasEntered int) bool {
		return hostClientIPs[hasEntered] >= clientIP
	})

	sort.Strings(hostClientIPs)

	if len(hostClientIPs) <= 0 {
		w.Write([]byte("No one has entered the training yet!"))
	} else if hasEntered < len(hostClientIPs) && clientIP == hostClientIPs[hasEntered] {
		train.EnterPYTraining(hostClientIPs, remoteHostIP, index)
	} else {
		w.Write([]byte("You haven't entered the training yet!"))
	}
}
