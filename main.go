package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"log"
	"os"
	"os/exec"
	"strings"
)

var logger *log.Logger
var tagDelimiter string

type execResult struct {
	out    bytes.Buffer
	stdErr bytes.Buffer
	err    error
}

// Tag is the AWS tag of a resource
type Tag struct {
	Key   string
	Value string
}

func runCommand(cmd string, args []string) execResult {
	var result execResult
	exec := exec.Command(cmd, args...)
	exec.Stdout = &result.out
	exec.Stderr = &result.stdErr
	result.err = exec.Run()

	return result
}

func main() {
	flag.StringVar(&tagDelimiter, "tagDelimiter", "k8s-label_", "delimiter which identifies a tag on the ec2 instance which should be used for node label")
	flag.Parse()

	logger = log.New(os.Stderr, "[AWS Node Labels]", log.LstdFlags)
	result := runCommand("kubectl", []string{"get", "nodes", "-o", "name"})

	if result.err != nil {
		logger.Println(result.stdErr.String())
		os.Exit(1)
	}

	nodes := strings.Split(result.out.String(), "\n")

	for _, node := range nodes {
		if len(node) > 0 {
			dns := strings.TrimLeft(node, "node/")

			if len(dns) > 0 {
				// Get the AWS EC2 instance's tags for the node that is in the cluster
				result = runCommand("aws", []string{"ec2", "describe-instances", "--filters=" + "Name=private-dns-name,Values=" + dns, "--query=Reservations[].Instances[].Tags[]"})

				if result.err != nil {
					logger.Println(result.stdErr.String())
					continue
				}

				var tags []*Tag
				err := json.Unmarshal(result.out.Bytes(), &tags)

				if err != nil {
					logger.Println(err)
					continue
				}

				for _, tag := range tags {
					// Check if the tag was intended to be a node label in kubernetes
					if strings.Contains(tag.Key, tagDelimiter) {
						label := strings.Split(tag.Key, tagDelimiter)

						if len(label) > 0 {
							result = runCommand("kubectl", []string{"label", "nodes", dns, label[1] + "=" + tag.Value})

							if result.err != nil {
								logger.Println(result.stdErr.String())
							}
						}
					}
				}
			}
		}
	}
}
