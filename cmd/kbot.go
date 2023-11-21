/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	polygon "github.com/polygon-io/client-go/rest"
	"github.com/polygon-io/client-go/rest/models"
	"github.com/spf13/cobra"
	telebot "gopkg.in/telebot.v3"
)

var (
	// TeleToken bot
	TeleToken = os.Getenv("TELE_TOKEN")
)

var (
	// Polygon API key
	PolygonAPIKey = os.Getenv("POLYGON_API_KEY")
)

// kbotCmd represents the kbot command
var kbotCmd = &cobra.Command{
	Use:     "kbot",
	Aliases: []string{"start"},
	Short:   "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("kbot %s started", appVersion)

		kbot, err := telebot.NewBot(telebot.Settings{
			URL:    "",
			Token:  TeleToken,
			Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
		})

		polygonClient := polygon.New(PolygonAPIKey)

		if err != nil {
			log.Fatalf("Please check TELE_TOKEN env variable %s", err)
			return
		}

		kbot.Handle("/ticker", func(c telebot.Context) error {
			ticker := c.Message().Payload

			log.Printf("Ticker: %s\n", ticker)

			params := models.GetPreviousCloseAggParams{
				Ticker: ticker,
			}.WithAdjusted(true)

			// make request
			res, err := polygonClient.GetPreviousCloseAgg(context.Background(), params)
			if err != nil {
				log.Fatal(err)
				return err
			}

			if res.ResultsCount < 0 {
				sendErr := c.Send("Ticker wasn't found!")
				return sendErr
			}

			c.Send(fmt.Sprintf("Open price for %s: %.2f$", ticker, res.Results[0].Open))
			c.Send(fmt.Sprintf("Close price for %s: %.2f$", ticker, res.Results[0].Close))
			return err
		})

		kbot.Handle(telebot.OnText, func(m telebot.Context) error {

			log.Print(m.Message().Payload, m.Text())
			payload := m.Message().Payload

			switch payload {
			case "hello":
				err := m.Send(fmt.Sprintf("Hello I'm Kbot %s!", appVersion))
				return err
			}

			return err

		})

		kbot.Start()
	},
}

func init() {
	rootCmd.AddCommand(kbotCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// kbotCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// kbotCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
