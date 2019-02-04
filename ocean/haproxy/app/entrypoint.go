package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
)

const version = "0.2.1"
const internalSSLDirectory = "/etc/ssl"
const externalSSLDirectory = "/etc/letsencrypt/live/"

type entryPoint struct {
	certs        string
	email        string
	proxyChan    chan int
	proxyCounter int
}

func (ep *entryPoint) Execute() {
	log.Printf("Starting Ocean Proxy Daemon... Version: %v \r\n", version)

	log.Println("Starting HTTP HAProxy...")
	go ep.StartProxy("/usr/local/etc/haproxy/haproxy.http.cfg")

	log.Println("Checking certs...")
	ep.ReviewCerts()

	log.Println("Certificate Initialization Process Complete!")

	ep.proxyChan <- 1 // Stop Proxy

	log.Println("Starting HTTPS HAProxy...")
	go ep.StartProxy("/usr/local/etc/haproxy/haproxy.https.cfg")

	log.Println("Ocean Proxy Service Running...")

	// ticker every 20 days
	ticker := time.NewTicker(time.Hour * 480)
	for _ = range ticker.C {
		log.Println("Renewing Certs (forced)...")
		// Renew CERT
		ep.RenewCerts()

		// Restart HAProxy
		log.Println("Stopped Proxy...")
		ep.proxyChan <- 1 // Stop Proxy

		log.Println("Starting HTTPS HAProxy...")
		go ep.StartProxy("/usr/local/etc/haproxy/haproxy.https.cfg")
	}
}

func (ep *entryPoint) StartProxy(config string) {
	cmd := exec.Command("haproxy", "-f", config)

	// cmd.Stdin

	var stdOut bytes.Buffer
	cmd.Stdout = &stdOut

	var stdErr bytes.Buffer
	cmd.Stderr = &stdErr

	err := cmd.Start()
	if err != nil {
		log.Printf("Error starting Haproxy! %v", err)
		log.Printf("Err: %q \r\n", stdErr.String())
		log.Panicf("Out: %q \r\n", stdOut.String())
	}

	go func(cmd2 *exec.Cmd) {
		for {
			select {
			case _ = <-ep.proxyChan:
				// Stop Haproxy
				ep.proxyCounter++
				cmd2.Process.Kill()
			}
		}
	}(cmd)

	err = cmd.Wait()
	if err != nil {
		log.Printf("Stopped proxy... %v", err)
	}
}

func (ep *entryPoint) ExistingCerts(url string) (exists bool, path string) {
	externalDir := fmt.Sprintf("%v/%v", externalSSLDirectory, url)
	// /etc/letsencrypt/live/url/url.pem
	// TODO Under what conditions does certbot attach -####
	if _, err := os.Stat(externalDir); err == nil {
		return true, externalDir
	}
	return false, ""
}

func (ep *entryPoint) BuildCert(cert, email string) {

	cmd := exec.Command("certbot", "certonly", "--standalone", "-d", cert,
		"--non-interactive", "--agree-tos", "--email", email, "--http-01-port=8888")

	var stdOut bytes.Buffer
	cmd.Stdout = &stdOut

	var stdErr bytes.Buffer
	cmd.Stderr = &stdErr

	err := cmd.Run()
	if err != nil {
		log.Printf("Error using certbot! %v", err)
		log.Printf("Err: %q \r\n", stdErr.String())
		log.Fatalf("Out: %q \r\n", stdOut.String())
	}
}

func (ep *entryPoint) MergeFiles(outputFilename string, files ...string) {
	outputFile, err := os.Create(outputFilename)
	if err != nil {
		log.Panicf("Error writing file %v", outputFilename)
	}

	defer outputFile.Close()

	for _, file := range files {
		data, err := ioutil.ReadFile(file)
		if err != nil {
			log.Panicf("Error reading file %v", file)
		}
		outputFile.Write(data)
	}
}

func (ep *entryPoint) MergeCert(cert string) {
	dir := fmt.Sprintf("%v/%v", internalSSLDirectory, cert)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		log.Printf("Building internal SSL directory %v \n", dir)

		// Internal SSL directory missing, building
		cmd := exec.Command("mkdir", "-p", dir)
		err := cmd.Run()
		if err != nil {
			log.Panicf("Error creating internal SSL directory %v", dir)
		}
	}

	if exists, externalDir := ep.ExistingCerts(cert); exists {
		// Merge fullchain.pem and privkey.pem into one file
		privkeyFile := fmt.Sprintf("%v/privkey.pem", externalDir)
		fullchainFile := fmt.Sprintf("%v/fullchain.pem", externalDir)
		fullcertFilename := fmt.Sprintf("%v/%v/%v.pem", internalSSLDirectory, cert, cert)
		ep.MergeFiles(fullcertFilename, fullchainFile, privkeyFile)
	} else {
		log.Panicf("Unable to find certs for %v", cert)
	}
}

func (ep *entryPoint) ReviewCerts() {
	for _, cert := range strings.Split(ep.certs, ",") {
		cert = strings.TrimSpace(cert)

		if exists, _ := ep.ExistingCerts(cert); !exists {
			log.Printf("Building cert %v ...\n", cert)
			ep.BuildCert(cert, ep.email)
		}
		log.Printf("Merging cert %v ...\n", cert)
		ep.MergeCert(cert)
	}
}

func (ep *entryPoint) RenewCerts() {
	cmd := exec.Command("certbot", "renew", "--force-renewal", "--tls-sni-01-port=8888")

	// cmd.Stdin

	var stdOut bytes.Buffer
	cmd.Stdout = &stdOut

	var stdErr bytes.Buffer
	cmd.Stderr = &stdErr

	err := cmd.Run()
	if err != nil {
		log.Printf("Error renewing CERTS! %v", err)
		log.Printf("Err: %q \r\n", stdErr.String())
		log.Panicf("Out: %q \r\n", stdOut.String())
	}

	// For each cert we need to merge
	for _, cert := range strings.Split(ep.certs, ",") {
		cert = strings.TrimSpace(cert)

		if exists, _ := ep.ExistingCerts(cert); !exists {
			log.Printf("Building cert %v ...\n", cert)
			ep.BuildCert(cert, ep.email)
		}
		log.Printf("Merging cert %v ...\n", cert)
		ep.MergeCert(cert)
	}
}

func main() {
	ep := &entryPoint{
		certs:        os.Getenv("CERTS"),
		email:        os.Getenv("EMAIL"),
		proxyCounter: 0,
		proxyChan:    make(chan int),
	}
	ep.Execute()
}
