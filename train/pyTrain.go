package train

import (
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func StartPYTraining(clientIPs []string, hostIP string) {
	var trainPhrase string

	tmp :=exec.Command("/bin/bash","-c", "pwd;")
	tmpPwd, _ := tmp.CombinedOutput()

	pwd := strings.Replace(string(tmpPwd),"\n","", -1) + "/train/"

	trainPhrase += "python3 " + pwd + "distribute.py --ps_hosts=" + hostIP + ":10563 --worker_hosts="

	trainPhrase += clientIPs[0] + ":2225"

	for i := 1; i < len(clientIPs); i++ {
		trainPhrase += "," + clientIPs[i] + ":2225"
	}

	trainPhrase += " --job_name=ps --task_index=0"

	cmd := exec.Command("/bin/bash","-c", "source " + pwd + "venv/bin/activate; " + trainPhrase)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	log.Println(cmd.Start())
}

func EnterPYTraining(clientIPs []string, hostIP string, index int) {
	var trainPhrase string

	tmp :=exec.Command("/bin/bash","-c", "pwd;")
	tmpPwd, _ := tmp.CombinedOutput()

	pwd := strings.Replace(string(tmpPwd),"\n","", -1) + "/train/"

	trainPhrase += "python3 " + pwd +"distribute.py --ps_hosts=" + hostIP + ":10563 --worker_hosts="

	trainPhrase += clientIPs[0] + ":2225"

	for i := 1; i < len(clientIPs); i++ {
		trainPhrase += "," + clientIPs[i] + ":2225"
	}

	trainPhrase += " --job_name=worker --task_index=" + strconv.Itoa(index)

	//log.Println(trainPhrase)

	cmd := exec.Command("/bin/bash","-c", "source " + pwd + "venv/bin/activate; " + trainPhrase)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	log.Println(cmd.Start())
}
