package config

type Config struct {
	HTTP    HTTP
	Storage Storage
	Broker  Broker
}

type HTTP struct {
	Port int `default:"8888"`
}

type Storage struct {
	Type     string `default:"memory"`
	Database Database
}

type Database struct {
	Driver string `default:"sqlite3"`
	DSN    string `default:"state.db"`
}

type Broker struct {
	DSN      string `default:"tcp://localhost:1883"`
	ClientID string `default:"elsa"`
}
