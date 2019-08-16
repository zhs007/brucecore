package brucecore

import (
	"io/ioutil"
	"os"

	"go.uber.org/zap/zapcore"
	"gopkg.in/yaml.v2"

	brucecorebase "github.com/zhs007/brucecore/base"
)

// Config - config
type Config struct {

	//------------------------------------------------------------------
	// adarender configuration

	// AdaCoreServAddr - Ada core service address
	AdaCoreServAddr string
	// AdaCoreToken - This is a valid adacoreserv token
	AdaCoreToken string

	//------------------------------------------------------------------
	// jarviscrawler configuration

	// JarvisCrawlerServAddr - Jarvis Crawler service address
	JarvisCrawlerServAddr string
	// JarvisCrawlerToken - This is a valid jarviscrawler token
	JarvisCrawlerToken string

	//------------------------------------------------------------------
	// brucenode service configuration

	// ClientTokens - There are the valid clienttokens for this node
	ClientTokens []string

	//------------------------------------------------------------------
	// logger configuration

	Log struct {
		// LogPath - log path
		LogPath string
		// LogLevel - log level, it can be debug, info, warn, error
		LogLevel string
		// LogConsole - it can be output to console
		LogConsole bool
	}
}

func getLogLevel(str string) zapcore.Level {
	switch str {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	default:
		return zapcore.ErrorLevel
	}
}

func checkConfig(cfg *Config) error {
	if cfg.AdaCoreServAddr == "" {
		return brucecorebase.ErrConfigNoAdaCoreServAddr
	}

	if cfg.AdaCoreToken == "" {
		return brucecorebase.ErrConfigNoAdaCoreToken
	}

	return nil
}

// LoadConfig - load config
func LoadConfig(filename string) (*Config, error) {
	fi, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	defer fi.Close()
	fd, err := ioutil.ReadAll(fi)
	if err != nil {
		return nil, err
	}

	cfg := &Config{}

	err = yaml.Unmarshal(fd, cfg)
	if err != nil {
		return nil, err
	}

	err = checkConfig(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

// HasToken - has token
func (cfg *Config) HasToken(token string) bool {
	// for _, v := range cfg.ClientTokens {
	// 	if v == token {
	// 		return true
	// 	}
	// }

	return false
}

// InitLogger - init logger
func InitLogger(cfg *Config) {
	brucecorebase.InitLogger(getLogLevel(cfg.Log.LogLevel), cfg.Log.LogConsole, cfg.Log.LogPath)
}
