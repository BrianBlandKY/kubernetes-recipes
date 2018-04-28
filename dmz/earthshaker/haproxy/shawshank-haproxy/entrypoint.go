package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
)

const version = "0.1.0"
const internalSSLDirectory = "/etc/ssl"
const externalSSLDirectory = "/etc/letsencrypt/live/"

/*

	Two Goals:
		Spin up Haproxy instance
		Iterate over certs and create/initialize if needed
		// Instead of cron job, maybe we can use this process to trigger
		// cert update.
*/

type entryPoint struct {
	certs string
	email string
}

func (ep *entryPoint) Execute() {
	log.Printf("Starting Shawshank Custom System Provision... Verison: %v \r\n", version)

	log.Println("Starting HAProxy...")
	// This will panic on failure (don't program for recovery, let process die)
	ep.StartProxy() // HAProxy

	log.Println("Checking certs...")

	for _, cert := range strings.Split(ep.certs, ",") {
		cert = strings.TrimSpace(cert)

		if exists, _ := ep.ExistingCerts(cert); !exists {
			log.Printf("Building cert %v ...\n", cert)
			ep.BuildCert(cert, ep.email)
		}
		log.Printf("Merging cert %v ...\n", cert)
		ep.MergeCert(cert)
	}

	log.Println("Certificate Initialization Process Complete!")

	for {
		// Idle
	}
	// TODO
	// Start renewal process
	// Kill and restart haproxy
}

func (ep *entryPoint) StartProxy() {
	cmd := exec.Command("haproxy", "-f", "/usr/local/etc/haproxy/haproxy.cfg", "-D")

	// cmd.Stdin

	var stdOut bytes.Buffer
	cmd.Stdout = &stdOut

	var stdErr bytes.Buffer
	cmd.Stderr = &stdErr

	err := cmd.Run()
	if err != nil {
		log.Printf("Error starting Haproxy! %v", err)
		log.Printf("Err: %q \r\n", stdErr.String())
		log.Fatalf("Out: %q \r\n", stdOut.String())
	}
	// May need a way to hold on to this process and kill/restart on cert renewal
}

func (ep *entryPoint) ExistingCerts(url string) (exists bool, path string) {
	externalDir := fmt.Sprintf("%v/%v", externalSSLDirectory, url)
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
		fullcertFilename := fmt.Sprintf("%v/%v.pem", internalSSLDirectory, cert)
		ep.MergeFiles(fullcertFilename, fullchainFile, privkeyFile)
	} else {
		log.Panicf("Unable to find certs for %v", cert)
	}
}

func main() {
	ep := &entryPoint{
		certs: os.Getenv("CERTS"),
		email: os.Getenv("EMAIL"),
	}
	ep.Execute()
}
