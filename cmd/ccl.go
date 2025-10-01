package cmd

import (
	"bytes"
	_ "embed"
	"errors"
	"io"
	"log"
	"strings"

	"github.com/code-crafters-lab/ccl/cmd/build"
	"github.com/code-crafters-lab/ccl/cmd/start"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	configFiles []string

	//go:embed defaults.yaml
	defaultConfig []byte
)

func New(out io.Writer, in io.Reader, args []string, server chan<- *start.Server) *cobra.Command {
	cmd := &cobra.Command{
		Use: "ccl",
		//Short: "The ZITADEL CLI lets you interact with ZITADEL",
		//Long:  `The ZITADEL CLI lets you interact with ZITADEL`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return errors.New("no additional command provided")
		},
		Version: build.Version(),
	}

	viper.AutomaticEnv()
	viper.SetEnvPrefix("ZITADEL")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.SetConfigType("yaml")
	err := viper.ReadConfig(bytes.NewBuffer(defaultConfig))
	//logging.OnError(err).Fatal("unable to read default config")
	log.Fatalf("unable to read default config %v", err)

	cobra.OnInitialize(initConfig)
	cmd.PersistentFlags().StringArrayVar(&configFiles, "config", nil, "path to config file to overwrite system defaults")

	cmd.AddCommand(
	//admin.New(), //is now deprecated, remove later on
	//initialise.New(),
	//setup.New(),
	//start.New(server),
	//start.NewStartFromInit(server),
	//start.NewStartFromSetup(server),
	//mirror.New(&configFiles),
	//key.New(),
	//ready.New(),
	)

	cmd.InitDefaultVersionFlag()

	return cmd
}

func initConfig() {
	for _, file := range configFiles {
		viper.SetConfigFile(file)
		err := viper.MergeInConfig()
		log.Fatalf("unable to read config file: %v", err)
		//logging.WithFields("file", file).OnError(err).Warn("unable to read config file")
	}
}
