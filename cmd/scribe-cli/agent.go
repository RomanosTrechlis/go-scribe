package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/RomanosTrechlis/go-scribe/internal/util/fs"
	"github.com/RomanosTrechlis/go-scribe/scribe"
	"github.com/RomanosTrechlis/go-scribe/types"
	"gopkg.in/yaml.v2"
)

type AgentConfig types.AgentConfig

func getAgentQs() []q {
	qs := make([]q, 0)
	qs = append(qs, q{1, "port", "What is Agent's port", "8080", -1})
	qs = append(qs, q{2, "profile", "Does Agent provides profile info", "false", -1})
	qs = append(qs, q{3, "profile_port", "What is Agent's profile port", "2222", 2})
	qs = append(qs, q{4, "console", "Should Agent write logs on console", "false", -1})
	qs = append(qs, q{5, "verbose", "Should Agent be verbose", "false", -1})
	qs = append(qs, q{6, "mediator", "Where is the Mediator, if any", "", -1})
	qs = append(qs, q{7, "log_path", "Where should logs be written", "logs", -1})
	qs = append(qs, q{8, "log_file_size", "What's the maximum size of log files should be", "10MB", -1})
	qs = append(qs, q{9, "certificate", "Certificate's path", "", -1})
	qs = append(qs, q{10, "private_key", "Private Key path", "", -1})
	qs = append(qs, q{11, "certificate_authority", "Certificate Authority path", "", -1})
	return qs
}

func (ac *AgentConfig) addValueToAgentConfig(field, val string) error {
	if field == "port" {
		v, err := strconv.Atoi(val)
		if err != nil {
			return err
		}
		ac.Port = v
		return nil
	}
	if field == "profile" {
		v, err := strconv.ParseBool(val)
		if err != nil {
			return err
		}
		ac.Profile = v
	}
	if field == "profile_port" {
		v, err := strconv.Atoi(val)
		if err != nil {
			return err
		}
		ac.ProfilePort = v
	}
	if field == "console" {
		v, err := strconv.ParseBool(val)
		if err != nil {
			return err
		}
		ac.Console = v
	}
	if field == "verbose" {
		v, err := strconv.ParseBool(val)
		if err != nil {
			return err
		}
		ac.Verbose = v
	}
	if field == "log_path" {
		ac.LogPath = val
	}
	if field == "log_file_size" {
		size, err := scribe.LexicalToNumber(val)
		if err != nil {
			return err
		}
		ac.LogFileSize = size
	}

	if field == "certificate" {
		ac.Certificate = val
	}
	if field == "private_key" {
		ac.PrivateKey = val
	}
	if field == "certificate_authority" {
		ac.CertificateAuthority = val
	}
	return nil
}

func createAgentConfig(write bool) error {
	homeDir, err := getUserHomeDir()
	if err != nil {
		return err
	}

	ac := new(AgentConfig)
	scanner := bufio.NewScanner(os.Stdin)
	for _, question := range getAgentQs() {
		fmt.Printf("%s (default='%s')? ", question.question, question.defaultValue)
		scanner.Scan()
		val := strings.Trim(scanner.Text(), " ")
		if val == "" {
			val = question.defaultValue
		}
		err := ac.addValueToAgentConfig(question.name, val)
		if err != nil {
			return err
		}
	}

	b, err := yaml.Marshal(ac)
	if err != nil {
		return err
	}

	if !write {
		fmt.Fprint(os.Stdout, string(b))
		return nil
	}
	path := filepath.Join(homeDir, SCRIBE_DIR)
	err = fs.CreateFolderIfNotExist(path)
	if err != nil {
		return err
	}
	file := filepath.Join(path, SCRIBE_CONFIG)
	err = ioutil.WriteFile(file, b, 0644)
	if err != nil {
		return err
	}
	return nil
}
