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
	appId := viper.GetString("subscribing.appId")
	devId := viper.GetString("subscribing.devId")
	accessKey := viper.GetString("subscribing.accessKey")
	client := mqtt.NewClient(ctx, "ttnctl", appId, accessKey, "tcp://eu.thethings.network:1883")
	if err := client.Connect(); err != nil {
		ctx.WithError(err).Fatal("Could not connect")
	}
	token := client.SubscribeDeviceUplink(appId, devId, func(client mqtt.Client, appID string, devID string, req types.UplinkMessage) {

		metadata := fmt.Sprintf("Frequency: %g\n", req.Metadata.Frequency)
		metadata += fmt.Sprintf("Modulation: %s\n", req.Metadata.Modulation)
		metadata += fmt.Sprintf("DataRate: %s\n", req.Metadata.DataRate)
		metadata += fmt.Sprintf("Bitrate: %d\n", req.Metadata.Bitrate)
		metadata += fmt.Sprintf("CodingRate: %s\n", req.Metadata.CodingRate)
		metadata += fmt.Sprintf("Location Lat/Long/Alt: %g/%g/%d\n",
					req.Metadata.LocationMetadata.Latitude,
					req.Metadata.LocationMetadata.Longitude,
					req.Metadata.LocationMetadata.Altitude)

		reader := strings.NewReader(`{"data": "` + metadata + `"}`)

		clientId := viper.GetString("posting.clientId")
		password := viper.GetString("posting.password")
		topicId := viper.GetString("posting.topicId")
		apiKey := viper.GetString("posting.apiKey")

		url := "https://realtime.opensensors.io/v1/topics/" + topicId + "?client-id=" + clientId + "&password=" + password
		request, err := http.NewRequest("POST", url, reader)
		if err != nil {
			ctx.WithError(err).Fatal("Could not do request")
		}
		request.Header.Set("Content-Type", "application/json")
		request.Header.Set("Accept", "application/json")
		request.Header.Set("Authorization", "api-key " + apiKey)
		resp, err := http.DefaultClient.Do(request)
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
