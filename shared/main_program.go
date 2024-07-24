package shared

type AppConfig struct {
	Port int        `yaml:"port" json:"port"`
	Mqtt MqttConfig `yaml:"mqtt" json:"mqtt"`
}

type MqttConfig struct {
	Broker   string `yaml:"broker" json:"broker"`
	Port     int    `yaml:"port" json:"port"`
	Username string `yaml:"username" json:"username"`
	Password string `yaml:"password" json:"password"`
}

type MainProgram interface {
	Start(cfg AppConfig)

	Stop()
}
