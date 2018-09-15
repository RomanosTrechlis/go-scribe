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
	"gopkg.in/yaml.v2"
)

type q struct {
	order        int
	name         string
	question     string
	defaultValue string
	dependency   int
}

func getMediatorQs() []q {
	qs := make([]q, 0)
	qs = append(qs, q{1, "port", "What is Mediator's port", "8000", -1})
	qs = append(qs, q{2, "profile", "Does Mediator provides profile info", "false", -1})
	qs = append(qs, q{3, "profile_port", "What is Mediator's profile port", "2222", 2})
	qs = append(qs, q{4, "certificate", "Certificate's path", "", -1})
	qs = append(qs, q{5, "private_key", "Private Key path", "", -1})
	qs = append(qs, q{6, "certificate_authority", "Certificate Authority path", "", -1})
	return qs
}

type MediatorConfig struct {
	Port  int  `yaml:"port"`
	Profile bool `yaml:"profile"`
	ProfilePort int  `yaml:"profile_port"`

	Certificate string `yaml:"certificate"`
	PrivateKey  string `yaml:"private_key"`
	CertificateAuthority  string `yaml:"certificate_authority"`
}

func (mc *MediatorConfig) addValueToMediatorConfig(field, val string) error {
	if field == "port" {
		v, err := strconv.Atoi(val)
		if err != nil {
			return err
		}
		mc.Port = v
		return nil
	}
	if field == "profile" {
		v, err := strconv.ParseBool(val)
		if err != nil {
			return err
		}
		mc.Profile = v
	}
	if field == "profile_port" {
		v, err := strconv.Atoi(val)
		if err != nil {
			return err
		}
		mc.ProfilePort = v
	}
	if field == "certificate" {
		mc.Certificate = val
	}
	if field == "private_key" {
		mc.PrivateKey = val
	}
	if field == "certificate_authority" {
		mc.CertificateAuthority = val
	}
	return nil
}

func createMediatorConfig(write bool) error {
	homeDir, err := getUserHomeDir()
	if err != nil {
		return err
	}

	mc := new(MediatorConfig)
	scanner := bufio.NewScanner(os.Stdin)
	for _, question := range getMediatorQs() {
		fmt.Printf("%s (default='%s')? ", question.question, question.defaultValue)
		scanner.Scan()
		val := strings.Trim(scanner.Text(), " ")
		if val == "" {
			val = question.defaultValue
		}
		mc.addValueToMediatorConfig(question.name, val)
	}

	b, err := yaml.Marshal(mc)
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
	file := filepath.Join(path, MEDIATOR_CONFIG)
	err = ioutil.WriteFile(file, b, 0644)
	if err != nil {
		return err
	}
	return nil
}
