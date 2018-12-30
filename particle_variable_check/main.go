package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
        "log"
        "net/http"
//	"github.com/sensu/sensu-go/types"
	"github.com/spf13/cobra"
)

var (
        access_token string 
	device_id    string
	product_id   string
        variable     string
        urlStr          string
        verbose bool 
)

type CoreInfo struct {
	Name string
	DeviceID string
	Connected bool
	Last_handshake_at string
        Last_app string
}
type ParticleVar struct {
	Name string
        Result string
        CoreInfo CoreInfo
}

func main() {
	rootCmd := configureRootCommand()
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func configureRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "particle_variable_check",
		Short: "Retrieve particle variable value and output in desired metric format",
		RunE:  run,
	}

	cmd.Flags().StringVarP(&device_id,
		"device",
		"d",
		"",
		"Particle Device ID")

	_ = cmd.MarkFlagRequired("device")

	cmd.Flags().StringVarP(&access_token,
		"token",
		"t",
		"",
		"Particle Access Token")
	_ = cmd.MarkFlagRequired("token")

	cmd.Flags().StringVarP(&variable,
		"variable",
		"v",
		"",
		"Particle Access Token")
	_ = cmd.MarkFlagRequired("variable")
	

	cmd.Flags().StringVarP(&product_id,
		"product",
		"p",
		"",
		"Optional Particle Product ID")

        cmd.Flags().BoolVar(&verbose,
		"verbose",
		false,
		"Enable verbose output")

	return cmd
}

func run(cmd *cobra.Command, args []string) error {
	if len(args) != 0 {
		_ = cmd.Help()
		return fmt.Errorf("invalid argument(s) received")
	}

	return particleDevicePing()
}

func particleDevicePing() error {
        urlStr="https://api.particle.io/v1/"
        if product_id != "" {
        urlStr+="products/"+product_id+"/"
        }         
        urlStr+="devices/"+device_id+"/"+variable+"?access_token="+access_token
	if verbose {
          fmt.Printf("Device:%s Token:%s Variable:%s\n", device_id,access_token,variable)
          if product_id != "" {
            fmt.Printf("  Product:%s\n", product_id)
          }
          fmt.Printf("  Url:%s\n", urlStr)
        }
        MakeRequest(urlStr)
	return nil
}


func MakeRequest(urlStr string) {
	resp, err := http.Get(urlStr)
	if err != nil {
		log.Fatalln(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}	
        fmt.Printf("Response: %s\n", body)

        var output ParticleVar 
        json.Unmarshal([]byte(body), &output)
        fmt.Printf("Var:%s Val:%s\n", output.Name,output.Result)
}

