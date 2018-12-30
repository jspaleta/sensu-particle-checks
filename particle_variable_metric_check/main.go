package main

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"
)

var (
	accessToken string
	deviceID    string
	productID   string
	metricName  string
	variable     string
	timestamp    string
	timeout      = 60
	verbose      bool
)

type coreInfo struct {
	DeviceID          string
	ProductID         int    `json:"product_id"`
	Connected         bool
	LastHandshakeAt   string `json:"last_handshake_at"`
	LastApp           string `json:"last_app"`
	LastHeard         string `json:"last_heard"`
}
type particleVar struct {
	Name      string
	Result    string
	Timestamp int
	CoreInfo  coreInfo
}

func main() {
	rootCmd := configureRootCommand()
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	} else {
        }
}

func configureRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "particle_variable_check",
		Short: "Retrieve Particle variable as graphite metric",
		Long:  `Retrieve Particle string variable and output in graphite plaintext format`,
		RunE:  run,
	}

	cmd.Flags().StringVarP(&deviceID,
		"device",
		"d",
		"",
		"Particle Device ID")

	_ = cmd.MarkFlagRequired("device")

	cmd.Flags().StringVarP(&accessToken,
		"access_token",
		"a",
		"",
		"Particle Access Token")
	_ = cmd.MarkFlagRequired("access_token")

	cmd.Flags().StringVarP(&variable,
		"variable",
		"v",
		"",
		"Particle Variable name, must hold string value")
	_ = cmd.MarkFlagRequired("variable")

	cmd.Flags().StringVarP(&timestamp,
		"timestamp",
		"t",
		"",
		"Optional Particle Timestamp Variable, must hold string representation of Unix Epoch integer")

	cmd.Flags().IntVarP(&timeout,
		"timeout",
		"T",
		60,
		"Optional particle Metric Timestamp Timeout (seconds)")

	cmd.Flags().StringVarP(&productID,
		"product",
		"p",
		"",
		"Optional Particle Product ID")

	cmd.Flags().StringVarP(&metricName,
		"metric",
		"m",
		"",
		"Optional Metric Name, if not set will be determined from hostname.variable")

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

	var output particleVar
	if timeout < 0 {
		timeout = 0
	}
	t := time.Now().Unix() 
	err := particleDeviceVariable(&output)
        if (err != nil) {
          return nil
        }
	if timestamp == "" || timeout == 0 || (int(t) - int(output.Timestamp) ) <= int(timeout) {
		if metricName == "" {
			metricName, err = os.Hostname()
			metricName += "." + output.Name
		}
		fmt.Printf("%s %s %d\n", metricName, output.Result, output.Timestamp)
	} else {
		err = fmt.Errorf("stale variable measurement: %d - %d = %d",int(t),output.Timestamp, int(t)-output.Timestamp)
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	return err
}

func particleDeviceVariable(output *particleVar) error {
	var baseURLStr string
	var variableURLStr string
        var err error
	var tout particleVar
	var body []byte

	baseURLStr = "https://api.particle.io/v1/"
	if productID != "" {
		baseURLStr += "products/" + productID + "/"
	}
	variableURLStr = baseURLStr + "devices/" + deviceID + "/" + variable + "?access_token=" + accessToken
	if verbose {
		fmt.Printf("Device:%s Token:%s Variable:%s\n", deviceID, accessToken, variable)
		if productID != "" {
			fmt.Printf("  Product:%s\n", productID)
		}
		fmt.Printf("  Url:%s\n", variableURLStr)
	}

	if timestamp != "" {
		var timestampURLStr string
		timestampURLStr = baseURLStr + "devices/" + deviceID + "/" + timestamp + "?access_token=" + accessToken
		var tbody []byte
		tbody, err = makeRequest(timestampURLStr)
        	if (err != nil) {
	  		fmt.Fprintf(os.Stderr, "error: %v\n", err)
	  		os.Exit(2)
        	}
		if verbose {
			fmt.Printf("Timestamp Response: %s\n", tbody)
		}
		json.Unmarshal([]byte(tbody), &tout)
	}

	body, err = makeRequest(variableURLStr)
        if (err != nil) {
	  fmt.Fprintf(os.Stderr, "error: %v\n", err)
	  os.Exit(2)
        }
	json.Unmarshal([]byte(body), &output)

	if timestamp == "" { 
        	output.Timestamp = int(time.Now().Unix()) 
        } else { output.Timestamp, err = strconv.Atoi(tout.Result) }
	if verbose {
		fmt.Printf("Response: %s\n", body)
		fmt.Printf("Var:%s Val:%s Timestamp: %d\n", output.Name, output.Result, output.Timestamp)
		fmt.Printf("DeviceID:%v LastHeard:%v LastHandshakeAt:%v ProductID: %v\n", output.CoreInfo.DeviceID,output.CoreInfo.LastHeard,output.CoreInfo.LastHandshakeAt,output.CoreInfo.ProductID)
	}
	return err
}

func makeRequest(urlStr string) ([]byte, error) {
	resp, err := http.Get(urlStr)
	if err != nil {
	  	fmt.Fprintf(os.Stderr, "error: %v\n", err)
	  	os.Exit(1)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
	  	fmt.Fprintf(os.Stderr, "error: %v\n", err)
	  	os.Exit(1)
	}
	if resp.StatusCode != 200 {
                err = fmt.Errorf("Failed Request %s StatusCode: %v", urlStr,resp.StatusCode)
	  	fmt.Fprintf(os.Stderr, "error: %v\n", err)
	  	os.Exit(1)
	}
	return body, err
}
