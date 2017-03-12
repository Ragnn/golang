//
// test for firstTry in /home/stanley/Documents/Internship/go/
//
// Made by Stanley Stephens
// Login   <stanley.stephens@epitech.eu>
//
// Started on  Fri Mar 10 23:43:37 2017 Stanley Stephens
// Last update Sun Mar 12 13:05:12 2017 Stanley Stephens
//

package main

import (
	"github.com/TheThingsNetwork/ttn/mqtt"
	"github.com/TheThingsNetwork/go-utils/log/apex"
	"github.com/TheThingsNetwork/go-utils/log"
	"github.com/TheThingsNetwork/ttn/core/types"
	"github.com/spf13/viper"
	"fmt"
	"os"
)

func main() {
	ctx := apex.Stdout().WithField("Test", "Go Client")
	log.Set(ctx)

	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println("Config file not found...")
		os.Exit(1)
	}
	appId := viper.GetString("config.appId")
	devId := viper.GetString("config.devId")
	accessKey := viper.GetString("config.accessKey")
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
