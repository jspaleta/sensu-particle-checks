package main

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"io/ioutil"
	"log"
	"net/http"
	//"net/url"
	"os"
	//"strings"
)

var (
	accessToken string
	deviceID    string
	productID   string
	verbose      bool
)

type ping struct {
	Online bool
	Ok     bool
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
		Use:   "particle_ping_check",
		Short: "Ping Particle device",
		Long:  `Ping Particle device and check to see if its online`,
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

	cmd.Flags().StringVarP(&productID,
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

	var output ping
	err := particlePing(&output)
	if err != nil {
		return err
	}
	fmt.Printf("Device: %s Online: %v Ok: %v\n", deviceID, output.Online,output.Ok)
        if ( !( output.Online && output.Ok )) { os.Exit(2) }
	return err
}

func particlePing(output *ping) error {
	var baseURLStr string
	baseURLStr = "https://api.particle.io/v1/"
	if productID != "" {
		baseURLStr += "products/" + productID + "/"
	}
	baseURLStr += "devices/" + deviceID + "/ping"
	if verbose {
		fmt.Printf("  Url:%s\n", baseURLStr)
	}
	body, err := makeRequest(baseURLStr, accessToken)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	json.Unmarshal([]byte(body), &output)
	if verbose {
		fmt.Printf("Response: %s\n", body)
	}
	return err
}

func makeRequest(urlStr string, accessToken string) ([]byte, error) {
        data := "Bearer "+accessToken
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPut, urlStr, nil)
	req.Header.Add("Authorization",data)
	fmt.Fprintf(os.Stderr, "Header: %v\n", req.Header)
	if err != nil {
		// handle error
		log.Fatal(err)
	}
	resp, err := client.Do(req)
	if err != nil {
		// handle error
		log.Fatal(err)
	}
	defer resp.Body.Close()
	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	//hdr := resp.Header
	if resp.StatusCode != 200 {
		err = fmt.Errorf("Failed Request %s StatusCode: %v\n%v", urlStr, resp.StatusCode,string(contents))
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	return contents, err
}
