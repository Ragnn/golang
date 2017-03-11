//
// test for firstTry in /home/stanley/Documents/Internship/go/
//
// Made by Stanley Stephens
// Login   <stanley.stephens@epitech.eu>
//
// Started on  Fri Mar 10 23:43:37 2017 Stanley Stephens
// Last update Sat Mar 11 22:05:12 2017 Stanley Stephens
//

package main

import (
	"github.com/TheThingsNetwork/ttn/mqtt"
	"github.com/TheThingsNetwork/go-utils/log/apex"
	"github.com/TheThingsNetwork/go-utils/log"
	"github.com/TheThingsNetwork/ttn/core/types"
	"fmt"
)

func main() {
	ctx := apex.Stdout().WithField("Test", "Go Client")
	log.Set(ctx)

	accessKey := "ttn-account-v2.OfuuW9smtu33PjpPtVAs54Bmc2dcgHEOywtuAT1oqzk"
	appId, devId := "office-app", "office-hq"
	client := mqtt.NewClient(ctx, "ttnctl", appId, accessKey, "tcp://eu.thethings.network:1883")
	if err := client.Connect(); err != nil {
		ctx.WithError(err).Fatal("Could not connect")
	}
	token := client.SubscribeDeviceUplink(appId, devId, func(client mqtt.Client, appID string, devID string, req types.UplinkMessage) {
	})
	token.Wait()

	if err := token.Error(); err != nil {
		ctx.WithError(err).Fatal("Could not subscribe")
	}

	var input string
	fmt.Print("Enter to quit ")
	fmt.Scanln(&input)
	client.Disconnect()
}