package main

import (
	_ "embed"
	"encoding/json"
	"fmt"
)

// Config struct untuk menyimpan konfigurasi
type Config struct {
	GENAI_API_KEY            string `json:"GENAI_API_KEY"`
	GENAI_MODEL_ID           string `json:"GENAI_MODEL_ID"`
	GENAI_SYSTEM_INSTRUCTION string `json:"GENAI_SYSTEM_INSTRUCTION"`
	CLIENT_JID               string `json:"CLIENT_JID"`
}

//go:embed config.json
var configFile []byte

var config Config

// InitConfig untuk inisialisasi konfigurasi
func InitConfig() error {
	err := json.Unmarshal(configFile, &config)
	if err != nil {
		return err
	}
	fmt.Println("Config", config)
	return nil
}

// GetConfig untuk mendapatkan nilai konfigurasi berdasarkan kunci
func GetConfig(key string) (string, error) {
	switch key {
	case "GENAI_API_KEY":
		return config.GENAI_API_KEY, nil
	case "GENAI_MODEL_ID":
		return config.GENAI_MODEL_ID, nil
	case "GENAI_SYSTEM_INSTRUCTION":
		return config.GENAI_SYSTEM_INSTRUCTION, nil
	case "CLIENT_JID":
		return config.CLIENT_JID, nil
	default:
		return "", fmt.Errorf("key %s not found in configuration", key)
	}
}
