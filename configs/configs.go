package configs

import (
	"encoding/json"
	"os"
)

var (
	config     map[string]interface{} = make(map[string]interface{})
	UseSSL     bool
	ServePort  string
	CertFile   string
	KeyFile    string
	ChannelDir string
)

func SetConfigFile(filePath string) (err error) {
	var (
		jsonByte []byte
	)

	if jsonByte, err = os.ReadFile(filePath); err != nil {
		return
	}

	if err = json.Unmarshal(jsonByte, &config); err != nil {
		return
	}

	if useSSL, isExist := config["UseSSL"]; isExist {
		UseSSL = useSSL.(bool)
	}

	if servePort, isExist := config["ServePort"]; isExist {
		ServePort = servePort.(string)
	}

	if channelDir, isExist := config["ChannelDir"]; isExist {
		ChannelDir = channelDir.(string)
	}

	if certFile, isExist := config["CertFile"]; isExist {
		CertFile = certFile.(string)
	}

	if keyFile, isExist := config["KeyFile"]; isExist {
		KeyFile = keyFile.(string)
	}

	return
}
