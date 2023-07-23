package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"text/template"
)

const (
	templateFile = "etcd.bt.template"
	outFile      = "etcd.bt"

	arm64ReqReg = "r3"
	amd64ReqReg = "ax"
)

var targets = []string{
	"(*kvServer).Range",
	"(*kvServer).Put",
}

type templateParams struct {
	EtcdBinaryPath string
	Addrs          map[string]string
	ReqReg         string
}

func main() {
	if len(os.Args) != 2 {
		log.Fatalf("usage: %s $PATH_TO_ETCD", os.Args[0])
	}
	etcdBinaryPath := os.Args[1]

	cmd := exec.Command("objdump", "-t", etcdBinaryPath)
	out, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}
	lines := bytes.Split(out, []byte("\n"))
	addrs := map[string]string{}
	for _, target := range targets {
		addr, err := findAddrInObjdump(lines, target)
		if err != nil {
			log.Fatal(err)
		}
		addrs[target] = addr
	}

	reqReg := arm64ReqReg
	if runtime.GOARCH == "amd64" {
		reqReg = arm64ReqReg
	}

	tmpl, err := template.ParseFiles(templateFile)
	if err != nil {
		log.Fatal(err)
	}
	file, err := os.Create(outFile)
	if err != nil {
		log.Fatal(err)
	}
	params := templateParams{
		EtcdBinaryPath: etcdBinaryPath,
		Addrs:          addrs,
		ReqReg:         reqReg,
	}
	tmpl.Execute(file, params)
}

func findAddrInObjdump(lines [][]byte, symbol string) (string, error) {
	for _, line := range lines {
		if bytes.Contains(line, []byte(symbol)) {
			return string(bytes.Split(line, []byte(" "))[0]), nil
		}
	}
	return "", fmt.Errorf("Symbol %s not found in objdump", symbol)
}
