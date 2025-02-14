package server

import (
	"encoding/json"
	"log"
	"os"
)



const (
    RulesPath = "server_rules.json"
    DefaultRules = "{\n\t\"initial_velocity\": 300,\n\t\"increment_velocity\": 100\n}"
)

func NewGameRules() *GameRules {
    buf, err := os.ReadFile(RulesPath)
    if err != nil {
        err := os.WriteFile(RulesPath, []byte(DefaultRules), 0644)
        if err != nil {
            log.Fatal("Failed to create rules file")
        }

        return &GameRules {
            InitVel: 300,
            IncVel: 100,
        }
    }

    rules := &GameRules{}
    json.Unmarshal(buf, rules)
    return rules
}

type GameRules struct {
    InitVel     float64 `json:"initial_velocity"`
    IncVel      float64 `json:"increment_velocity"`
}
