package configdb

import (
	"encoding/json"
	"log"
	"os"
	"path"
	"sync"
)

type Mu struct {
	Mu sync.Mutex
}

type GlobalOptions struct {
	IP                 string
	Gateway            string
	LeaseDuration      int
	AuthToken          string
	ListenInterface    string
	HTTPListenAddress  string
	TLSEnabled         bool
	TLSPrivateKeyFile  string
	TLSCertificateFile string
	Options            interface{}
}

type Clients struct {
	Clients []Client
}

type Client struct {
	Hostname string
	MAC      string
	IP       string
	Options  interface{}
}

var DB *Clients
var Config *GlobalOptions
var ConfigPath string

func Init(configPath string) error {
	ConfigPath = configPath
	content := []byte(`{}`)
	dataFile := path.Join(configPath, "rest-dhcpd-clients.json")
	_, err := os.Stat(dataFile)
	if !os.IsNotExist(err) {
		content, err = os.ReadFile(dataFile)
		if err != nil {
			return err
		}
	}
	if len(content) == 0 {
		content = []byte(`{}`)
	}
	err = json.Unmarshal(content, &DB)
	if err != nil {
		return err
	}

	cfg, err := os.ReadFile(path.Join(configPath, "rest-dhcpd-config.json"))
	if err != nil {
		return err
	}
	err = json.Unmarshal(cfg, &Config)
	if err != nil {
		return err
	}
	log.Printf("DB init done.")
	return nil
}

func (m *Mu) Save() error {
	content, err := json.MarshalIndent(DB, "", "  ")
	if err != nil {
		log.Printf("%s", err)
	}
	m.Mu.Lock()
	err = os.WriteFile(path.Join(ConfigPath, "rest-dhcpd-clients.json"), content, 0644)
	m.Mu.Unlock()
	return err
}
