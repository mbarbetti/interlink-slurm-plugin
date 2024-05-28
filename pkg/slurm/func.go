package slurm

import (
	"context"
	"flag"
	"fmt"
	"os"

	"k8s.io/client-go/kubernetes"

	"github.com/containerd/containerd/log"
	"gopkg.in/yaml.v2"
)

var SlurmConfigInst SlurmConfig
var Clientset *kubernetes.Clientset

// TODO: implement factory design

// NewSlurmConfig returns a variable of type SlurmConfig, used in many other functions and the first encountered error.
func NewSlurmConfig() (SlurmConfig, error) {
	if !SlurmConfigInst.set {
		var path string
		verbose := flag.Bool("verbose", false, "Enable or disable Debug level logging")
		errorsOnly := flag.Bool("errorsonly", false, "Prints only errors if enabled")
		SlurmConfigPath := flag.String("SlurmConfigpath", "", "Path to InterLink config")
		flag.Parse()

		if *verbose {
			SlurmConfigInst.VerboseLogging = true
			SlurmConfigInst.ErrorsOnlyLogging = false
		} else if *errorsOnly {
			SlurmConfigInst.VerboseLogging = false
			SlurmConfigInst.ErrorsOnlyLogging = true
		}

		if *SlurmConfigPath != "" {
			path = *SlurmConfigPath
		} else if os.Getenv("SLURMCONFIGPATH") != "" {
			path = os.Getenv("SLURMCONFIGPATH")
		} else {
			path = "/etc/interlink/SlurmConfig.yaml"
		}

		if _, err := os.Stat(path); err != nil {
			log.G(context.Background()).Error("File " + path + " doesn't exist. You can set a custom path by exporting INTERLINKCONFIGPATH. Exiting...")
			return SlurmConfig{}, err
		}

		log.G(context.Background()).Info("Loading SLURM config from " + path)
		yfile, err := os.ReadFile(path)
		if err != nil {
			log.G(context.Background()).Error("Error opening config file, exiting...")
			return SlurmConfig{}, err
		}
		yaml.Unmarshal(yfile, &SlurmConfigInst)

		if os.Getenv("SIDECARPORT") != "" {
			SlurmConfigInst.Sidecarport = os.Getenv("SIDECARPORT")
		}

		if os.Getenv("SBATCHPATH") != "" {
			SlurmConfigInst.Sbatchpath = os.Getenv("SBATCHPATH")
		}

		if os.Getenv("SCANCELPATH") != "" {
			SlurmConfigInst.Scancelpath = os.Getenv("SCANCELPATH")
		}

		if os.Getenv("TSOCKS") != "" {
			if os.Getenv("TSOCKS") != "true" && os.Getenv("TSOCKS") != "false" {
				fmt.Println("export TSOCKS as true or false")
				return SlurmConfig{}, err
			}
			if os.Getenv("TSOCKS") == "true" {
				SlurmConfigInst.Tsocks = true
			} else {
				SlurmConfigInst.Tsocks = false
			}
		}

		if os.Getenv("TSOCKSPATH") != "" {
			path = os.Getenv("TSOCKSPATH")
			if _, err := os.Stat(path); err != nil {
				log.G(context.Background()).Error("File " + path + " doesn't exist. You can set a custom path by exporting TSOCKSPATH. Exiting...")
				return SlurmConfig{}, err
			}

			SlurmConfigInst.Tsockspath = path
		}

		SlurmConfigInst.set = true
	}
	return SlurmConfigInst, nil
}
