package types

type CertificateConfig struct {
	Certificate          string `yaml:"certificate"`
	PrivateKey           string `yaml:"private_key"`
	CertificateAuthority string `yaml:"certificate_authority"`
}

type AgentConfig struct {
	Port        int    `yaml:"port"`
	Profile     bool   `yaml:"profile"`
	Console     bool   `yaml:"console"`
	Verbose     bool   `yaml:"verbose"`
	Mediator    string `yaml:"mediator"`
	ProfilePort int    `yaml:"profile_port"`
	LogPath     string `yaml:"log_path"`
	LogFileSize int64  `yaml:"log_file_size"`

	CertificateConfig
}

type MediatorConfig struct {
	Port        int  `yaml:"port"`
	Profile     bool `yaml:"profile"`
	ProfilePort int  `yaml:"profile_port"`

	CertificateConfig
}
