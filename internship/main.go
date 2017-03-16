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
	"github.com/spf13/viper"
	"fmt"
	"os"
	"strings"
	"net/http"
)

func main() {
	ctx := apex.Stdout().WithField("Test", "Go Client")
	log.Set(ctx)

	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
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

		// "My First Msg" should be change with the data in req.Metadata
		reader := strings.NewReader(`{"data": "My First Msg"}`)

		// <CLIENT ID> and <PWD> are logins for devices in opensensors.io
		// <TOPIC ID> find in the "topics" section in opensensors.io
		request, err := http.NewRequest("POST", "https://realtime.opensensors.io/v1/topics/<TOPIC ID>?client-id=<CLIENT ID>&password=<PWD>", reader)

		if err != nil {
			ctx.WithError(err).Fatal("Could not do request")
		}
		http.Header.Add("Content-Type", "application/json")
		http.Header.Add("Accept", "application/json")
		// <API KEY> is find in the profile page
		http.Header.Add("Authorization", "api-key <API KEY>")
		Client := &http.Client{}
		resp, err := Client.Do(request)
		resp = resp
		if err != nil {
			ctx.WithError(err).Fatal("Could not get response")
		}
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
