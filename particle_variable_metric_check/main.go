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
	access_token string
	device_id    string
	product_id   string
	metric_name  string
	variable     string
	timestamp    string
	timeout      int = 60
	verbose      bool
)

type CoreInfo struct {
	Name              string
	DeviceID          string
	Connected         bool
	Last_handshake_at string
	Last_app          string
}
type ParticleVar struct {
	Name      string
	Result    string
	Timestamp int
	CoreInfo  CoreInfo
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
		Short: "Retrieve Particle variable",
		Long:  `Retrieve Particle string variable and output in graphite plaintext format`,
		RunE:  run,
	}

	cmd.Flags().StringVarP(&device_id,
		"device",
		"d",
		"",
		"Particle Device ID")

	_ = cmd.MarkFlagRequired("device")

	cmd.Flags().StringVarP(&access_token,
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

	cmd.Flags().StringVarP(&product_id,
		"product",
		"p",
		"",
		"Optional Particle Product ID")

	cmd.Flags().StringVarP(&metric_name,
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

	var output ParticleVar
	if timeout < 0 {
		timeout = 0
	}
	t := time.Now().Unix() 
	err := particleDeviceVariable(&output)
        if (err != nil) {
          return nil
        }
	if timeout == 0 || (int(t) - int(output.Timestamp) ) <= int(timeout) {
		if metric_name == "" {
			metric_name, err = os.Hostname()
			metric_name += "." + output.Name
		}
		fmt.Printf("%s %s %d\n", metric_name, output.Result, output.Timestamp)
	} else {
		err = fmt.Errorf("Stale Variable Measurement: %d - %d = %d\n",int(t),output.Timestamp, int(t)-output.Timestamp)
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	return err
}

func particleDeviceVariable(output *ParticleVar) error {
	var baseUrlStr string
	var variableUrlStr string
	baseUrlStr = "https://api.particle.io/v1/"
	if product_id != "" {
		baseUrlStr += "products/" + product_id + "/"
	}
	variableUrlStr = baseUrlStr + "devices/" + device_id + "/" + variable + "?access_token=" + access_token
	if verbose {
		fmt.Printf("Device:%s Token:%s Variable:%s\n", device_id, access_token, variable)
		if product_id != "" {
			fmt.Printf("  Product:%s\n", product_id)
		}
		fmt.Printf("  Url:%s\n", variableUrlStr)
	}
	body, err := MakeRequest(variableUrlStr)
        if (err != nil) {
	  fmt.Fprintf(os.Stderr, "error: %v\n", err)
	  os.Exit(1)
        }

	json.Unmarshal([]byte(body), &output)
	if timestamp != "" {
		var timestampUrlStr string
		timestampUrlStr = baseUrlStr + "devices/" + device_id + "/" + timestamp + "?access_token=" + access_token
		var tout ParticleVar
		var tbody []byte
		tbody, err = MakeRequest(timestampUrlStr)
        	if (err != nil) {
	  		fmt.Fprintf(os.Stderr, "error: %v\n", err)
	  		os.Exit(1)
        	}
		if verbose {
			fmt.Printf("Timestamp Response: %s\n", tbody)
		}
		json.Unmarshal([]byte(tbody), &tout)
		output.Timestamp, err = strconv.Atoi(tout.Result)
	}
	if verbose {
		fmt.Printf("Response: %s\n", body)
		fmt.Printf("Var:%s Val:%s Timestamp: %d\n", output.Name, output.Result, output.Timestamp)
	}
	return err
}

func MakeRequest(urlStr string) ([]byte, error) {
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
