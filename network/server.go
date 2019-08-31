package network

import (
	"bufio"
	"bytes"
	"crypto/md5"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"github.com/wcytk/BlockTrain/train"
	"github.com/wcytk/BlockTrain/transaction"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
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

const (
	maxUploadSize = 2 * 1024 * 2014 // 2MB
)

var rootPath = ""

var uploadPath = ""

func (server *Server) Start() {
	if err := http.ListenAndServe(":12345", nil); err != nil {
		fmt.Println(err)
		return
	}
}

func (server *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Welcome to BlockTrain!"))
}

func (server *Server) setRoute() {
	tmp := exec.Command("/bin/bash", "-c", "pwd;")
	tmpPwd, _ := tmp.CombinedOutput()
	rootPath = strings.Replace(string(tmpPwd), "\n", "", -1)
	uploadPath = strings.Replace(string(tmpPwd), "\n", "", -1) + "/upload"
	http.HandleFunc("/", server.ServeHTTP)
	http.HandleFunc("/addIP", server.addIP)
	http.HandleFunc("/getFiles", server.getFiles)
	http.HandleFunc("/addToIPFS", server.addToIPFS)
	http.HandleFunc("/prepareTraining", server.prepareTraining)
	http.HandleFunc("/startTraining", server.startTraining)
	http.HandleFunc("/stopTraining", server.stopTraining)
	http.HandleFunc("/enterTraining", server.enterTraining)
	http.HandleFunc("/upload", uploadFileHandler())
	fs := http.FileServer(http.Dir(uploadPath))
	http.Handle("/files/", http.StripPrefix("/files", fs))
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

func (server *Server) getFiles(w http.ResponseWriter, r *http.Request) {
	filePath := uploadPath + "/file.txt"
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}
	defer file.Close()

	br := bufio.NewReader(file)
	fileInfo := make(map[string]string)
	files := make(map[int]map[string]string)
	i := 0
	for {
		a, c := br.ReadString('\n')

		if len(a) != 0 {
			fileInfo["fileHash"] = strings.Split(a, " ")[0]
			fileInfo["fileName"] = strings.Split(a, " ")[1]
			files[i] = fileInfo
			i++
		}
		if c != nil {
			if err == io.EOF {
				break
			}
			fmt.Println(err)
		}
	}
	data, err := json.Marshal(files)
	if err != nil {
		log.Println(err)
	}
	w.Write(data)
}

func (server *Server) addToIPFS(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
		return
	}
	if r.Form["fileName"][0] != "" {
		filePath := uploadPath + "/" + r.Form["fileName"][0]
		toIpfs := "ipfs add " + filePath + " | awk '{print $2 \" \" $3}' >> " + uploadPath + "/file.txt"

		cmd := exec.Command("/bin/bash", "-c", toIpfs)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		log.Println(cmd.Start())

		updateModel(r.Form["fileName"][0])
	}
}

func updateModel(fileName string) {
	filePath := uploadPath + "/" + fileName
	update := "cp " + filePath + " " + rootPath + "/train/distribute.py"

	cmd := exec.Command("/bin/bash", "-c", update)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	log.Println(cmd.Start())
}

func updateModelFromIPFS(fileHash string) {
	if isTraining == true {
		log.Println("You are currently training a model!")
	} else {
		download := "ipfs get " + fileHash
		cmd := exec.Command("/bin/bash", "-c", download)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		log.Println(cmd.Start())

		filePath := uploadPath + "/" + fileHash
		update := "cp " + filePath + " " + rootPath + "/train/distribute.py"

		cmd = exec.Command("/bin/bash", "-c", update)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		log.Println(cmd.Start())
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

	fileHash := formData["fileHash"].(string)

	hasEntered := sort.Search(len(hostClientIPs), func(hasEntered int) bool {
		return hostClientIPs[hasEntered] >= clientIP
	})

	updateModelFromIPFS(fileHash)

	sort.Strings(hostClientIPs)

	if len(hostClientIPs) <= 0 {
		w.Write([]byte("No one has entered the training yet!"))
	} else if hasEntered < len(hostClientIPs) && clientIP == hostClientIPs[hasEntered] {
		train.EnterPYTraining(hostClientIPs, remoteHostIP, index)
	} else {
		w.Write([]byte("You haven't entered the training yet!"))
	}

}

func uploadFileHandler() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// validate file size
		r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)
		if err := r.ParseMultipartForm(maxUploadSize); err != nil {
			renderError(w, "FILE_TOO_BIG", http.StatusBadRequest)
			return
		}

		// parse and validate file and post parameters
		file, handler, err := r.FormFile("file")
		if err != nil {
			renderError(w, "INVALID_FILE", http.StatusBadRequest)
			return
		}
		defer file.Close()
		fileBytes, err := ioutil.ReadAll(file)
		if err != nil {
			renderError(w, "INVALID_FILE", http.StatusBadRequest)
			return
		}

		// check file type, detectcontenttype only needs the first 512 bytes
		filetype := handler.Filename
		switch strings.Split(filetype, ".")[1] {
		case "py":
			break
		default:
			renderError(w, "INVALID_FILE_TYPE", http.StatusBadRequest)
			return
		}

		hex := md5.New()
		md5Name := hex.Sum([]byte(strings.Split(filetype, ".")[0]))
		fileName := fmt.Sprintf("%x", md5Name)
		randString := randToken(5)
		fileName = randString + fileName

		fileType := strings.Split(filetype, ".")[1]
		//fileEndings, err := mime.ExtensionsByType(fileType)
		if err != nil {
			renderError(w, err.Error(), http.StatusInternalServerError)
			return
		}
		newPath := filepath.Join(uploadPath, fileName+"."+fileType)
		fmt.Printf("FileType: %s, File: %s\n", fileType, newPath)

		// write file
		newFile, err := os.Create(newPath)
		if err != nil {
			renderError(w, "CANT_WRITE_FILE", http.StatusInternalServerError)
			return
		}
		defer newFile.Close() // idempotent, okay to call twice
		if _, err := newFile.Write(fileBytes); err != nil || newFile.Close() != nil {
			renderError(w, "CANT_WRITE_FILE", http.StatusInternalServerError)
			return
		}

		w.Write([]byte("SUCCESS"))

		toIpfs := "ipfs add " + uploadPath + "/" + fileName + "." + fileType + " | awk '{print $2 \" \" $3}' >> " + uploadPath + "/file.txt"

		cmd := exec.Command("/bin/bash", "-c", toIpfs)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		log.Println(cmd.Start())

		updateModel(fileName + "." + fileType)
	})
}

func renderError(w http.ResponseWriter, message string, statusCode int) {
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte(message))
}

func randToken(len int) string {
	b := make([]byte, len)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}
