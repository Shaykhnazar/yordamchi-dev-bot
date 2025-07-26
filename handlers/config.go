package handlers

import (
    "encoding/json"
    "fmt"
    "io"
    "math/rand"
    "os"
    "time"
)

type Config struct {
    Bot      BotConfig      `json:"bot"`
    Messages MessageConfig  `json:"messages"`
    Jokes    []string       `json:"jokes"`
    Quotes   []string       `json:"quotes"`
}

type BotConfig struct {
    Name        string `json:"name"`
    Version     string `json:"version"`
    Description string `json:"description"`
    Author      string `json:"author"`
}

type MessageConfig struct {
    Welcome        string `json:"welcome"`
    Help           string `json:"help"`
    UnknownCommand string `json:"unknown_command"`
}

func LoadConfig() (*Config, error) {
    file, err := os.Open("config.json")
    if err != nil {
        return nil, fmt.Errorf("config.json faylini ochishda xatolik: %w", err)
    }
    defer file.Close()

    data, err := io.ReadAll(file)
    if err != nil {
        return nil, fmt.Errorf("config.json faylini o'qishda xatolik: %w", err)
    }

    var config Config
    if err := json.Unmarshal(data, &config); err != nil {
        return nil, fmt.Errorf("config.json parse qilishda xatolik: %w", err)
    }

    rand.Seed(time.Now().UnixNano())
    
    return &config, nil
}

// Fixed functions to accept config parameter
func GetRandomJoke(config *Config) string {
    if len(config.Jokes) == 0 {
        return "😅 Hazillar yuklanmagan!"
    }
    return config.Jokes[rand.Intn(len(config.Jokes))]
}

func GetRandomQuote(config *Config) string {
    if len(config.Quotes) == 0 {
        return "💭 Iqtiboslar yuklanmagan!"
    }
    return config.Quotes[rand.Intn(len(config.Quotes))]
}