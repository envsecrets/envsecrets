/*
Copyright © 2023 Mrinal Wahal <mrinalwahal@gmail.com>
All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:

 1. Redistributions of source code must retain the above copyright notice,
    this list of conditions and the following disclaimer.

 2. Redistributions in binary form must reproduce the above copyright notice,
    this list of conditions and the following disclaimer in the documentation
    and/or other materials provided with the distribution.

 3. Neither the name of the copyright holder nor the names of its contributors
    may be used to endorse or promote products derived from this software
    without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE
LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR
CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF
SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS
INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN
CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE)
ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE
POSSIBILITY OF SUCH DAMAGE.
*/
package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/envsecrets/envsecrets/cli/commons"
	"github.com/envsecrets/envsecrets/cli/config"
	configCommons "github.com/envsecrets/envsecrets/cli/config/commons"
	"github.com/envsecrets/envsecrets/cli/internal/secrets"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var debug bool

var log = logrus.New()

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "envs",
	Short: "CLI-first manangement of your environment secrets and variables.",
	Long: `
envsecrets provides a centralized cloud account with rotate-able keys
to store the environment secrets and variables for all your projects in a single place
and integrate them with third-party services of your choice.

Homepage: https://envsecrets.com
Documentation: https://docs.envsecrets.com
DM me on Twitter for help: @MrinalWahal

Upgrade the CLI:

	MacOS => brew upgrade envsecrets/tap/envs
	Linux => snap refresh envs
`,
	Version: commons.VERSION,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {

		verbosity := "info"
		if debug {
			verbosity = "debug"
		}

		if err := setUpLogs(os.Stdout, verbosity); err != nil {
			return err
		}

		/* 		//	Load the project configuration.
		   		config, err := LoadConfig(commons.CONFIG_LOC)
		   		if err != nil {
		   			return err
		   		}
		fmt.Println("Project in local config", config.Project)
		*/

		//	Initialize configuration
		commons.Initialize(log)

		/* 		//	If the user is not already authenticated,
		   		//	log them in first.
		   		if args[0] != "login" {
		   			if !auth.IsLoggedIn() {
		   				loginCmd.Run(cmd, args)
		   			}
		   		}
		*/
		return nil
	},

	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

/* // LoadConfig reads configuration from file or environment variables.
func LoadConfig(path string) (config *commons.Config, err error) {
	viper.AddConfigPath(path)
	//viper.SetConfigName("app")
	viper.SetConfigType("yaml")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
*/
// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		log.Debug(err)
		os.Exit(1)
	}
}

type myFormatter struct {
	logrus.TextFormatter
}

func (f *myFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	// this whole mess of dealing with ansi color codes is required if you want the colored output otherwise you will lose colors in the log levels
	var levelColor int
	switch entry.Level {
	case logrus.DebugLevel, logrus.TraceLevel:
		levelColor = 37 // white
	case logrus.WarnLevel:
		levelColor = 33 // yellow
	case logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel:
		levelColor = 31 // red
	default:
		levelColor = 36 // cyan blue
	}
	return []byte(fmt.Sprintf("\x1b[%dm>\x1b[0m %s\n", levelColor, entry.Message)), nil
}

// setUpLogs set the log output ans the log level
func setUpLogs(out io.Writer, level string) error {
	log.SetOutput(out)
	lvl, err := logrus.ParseLevel(level)
	if err != nil {
		return err
	}
	log.SetLevel(lvl)
	log.SetFormatter(&myFormatter{logrus.TextFormatter{
		FullTimestamp:          true,
		TimestampFormat:        "2006-01-02 15:04:05",
		ForceColors:            true,
		DisableLevelTruncation: true,
	}})

	return nil
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	//	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.envsecrets.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "Print debug logs")
}

func InitializeSecret(log *logrus.Logger) {

	if log == nil {
		log = logrus.New()
	}

	var remoteConfig *secrets.RemoteConfig
	if environmentName != "" {

		log.Debug("Reading project config...")

		//	Fetch the project config
		projectConfig, err := config.GetService().Load(configCommons.ProjectConfig)
		if err != nil {

			if os.IsNotExist(err) {

				//	If the project config does not exist, begin the `init` command.
				initCmd.PreRunE(rootCmd, []string{})
				initCmd.Run(rootCmd, []string{})

				InitializeSecret(log)

			} else {
				log.Fatal(err)
			}
		} else {
			commons.ProjectConfig = projectConfig.(*configCommons.Project)
		}

		//	If the project config does not exist, throw an error.
		if commons.ProjectConfig == nil {
			log.Fatal("Project configuration not found")
		}

		remoteConfig = &secrets.RemoteConfig{
			EnvironmentName: environmentName,
			ProjectID:       commons.ProjectConfig.ProjectID,
		}
	}

	//	Initialize the local secret.
	secret, err := secrets.GetService().Init(commons.DefaultContext, commons.GQLClient, remoteConfig)
	if err != nil {
		log.Error(err)
		log.Fatal("Failed to initialize the local secret")
	}

	commons.Secret = secret
}
