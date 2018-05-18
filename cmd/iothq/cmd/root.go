package cmd

import (
    "github.com/sirupsen/logrus"
    "github.com/spf13/viper"
    "strings"
    "github.com/spf13/cobra"
    "github.com/rzajac/iotdet/version"
    "encoding/json"
    "os"
    "encoding/hex"
    "github.com/pkg/errors"
    "github.com/rzajac/iotdet/pkg/hq"
    "time"
)

// cfgFile holds path to the configuration file.
var cfgFile string

// cfg represents global configuration.
var cfg *hq.Config

func init() {
    cobra.OnInitialize(onInitialize)
    rootCmd.SetVersionTemplate(`{{.Version}}`)
    rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "path to configuration file (default is ./iotdet.yaml)")
    rootCmd.PersistentFlags().BoolP("version", "v", false, "version")
    rootCmd.PersistentFlags().BoolP("debug", "d", false, "run nin debug mode")
    viper.BindPFlag("debug", rootCmd.PersistentFlags().Lookup("debug"))
}

// Execute executes root command.
func Execute() {
    if err := rootCmd.Execute(); err != nil {
        logrus.WithFields(logrus.Fields{"service": "hq"}).Error(err)
        os.Exit(1)
    }
}

// rootCmd is the main command for the iotdet binary.
var rootCmd = &cobra.Command{
    Use:           "iothq",
    Version:       getVersion(),
    Short:         "IoT HQ.",
    Long:          `IoT HQ.`,
    SilenceUsage:  true,
    SilenceErrors: true,
    PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
        var err error
        cfg, err = config()
        return err
    },
}

// onInitialize runs before command Execute function is run.
func onInitialize() {
    // Add a prefix while reading from the environment variables.
    viper.SetEnvPrefix("HQ")
    // Replace dot with underscore when looking for environmental variables.
    viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
    // Name of the configuration file and where to look for it.
    viper.SetConfigName("hq")
    viper.AddConfigPath("/usr/etc/hq")
    viper.AddConfigPath(".")
    // Load config file if it was explicitly passed.
    if cfgFile != "" {
        viper.SetConfigFile(cfgFile)
    }
    viper.AutomaticEnv()

    // If a config file is found, read it in.
    viper.ReadInConfig()
}

// getVersion returns JSON formatted application version.
func getVersion() string {
    v := struct {
        Version, BuildDate, GitHash, GitTreeState string
    }{
        version.Version,
        version.BuildDate,
        version.GitHash,
        version.GitTreeState,
    }

    j, _ := json.Marshal(v)
    return string(j)
}

// config returns validated configuration structure.
func config() (*hq.Config, error) {
    c := &hq.Config{}

    // New agent detection configuration.
    c.DetItfName = viper.GetString("hq.detect.itf")
    if c.DetItfName == "" {
        return nil, errors.New("you must provide WiFi interface name")
    }
    c.DetApPass = viper.GetString("hq.detect.ap_pass")
    c.DetAgentIP = viper.GetString("hq.detect.agent_ip")
    c.DetUseIP = viper.GetString("hq.detect.use_ip")
    c.DetCmdPort = viper.GetInt("hq.detect.cmd_port")
    c.DetInterval = viper.GetDuration("hq.detect.scan_interval") * time.Second

    // Cipher configuration.
    ci := viper.GetString("hq.cipher")
    switch ci {
    case hq.CIPHER_AES:
        key, err := hex.DecodeString(viper.GetString("hq.cipher_aes.key"))
        if err != nil {
            return nil, errors.New("invalid AES key value")
        }

        vi, err := hex.DecodeString(viper.GetString("hq.cipher_aes.vi"))
        if err != nil {
            return nil, errors.New("invalid AES vi value")
        }
        c.Cipher = hq.NewCipherAES(key, vi)

    case hq.CIPHER_NONE:
        c.Cipher = hq.NewNoopCipher()

    default:
        return nil, errors.Errorf("invalid cipher name %s", ci)
    }

    // HQ access point configuration.
    c.HQApName = viper.GetString("hq.access_point.name")
    c.HQApPass = viper.GetString("hq.access_point.pass")

    // MQTT configuration.
    c.MQTTIP = viper.GetString("hq.mqtt.ip")
    c.MQTTPort = viper.GetInt("hq.mqtt.port")
    c.MQTTUser = viper.GetString("hq.mqtt.user")
    c.MQTTPass = viper.GetString("hq.mqtt.pass")

    // Setup logger.
    c.Log = logrus.WithFields(logrus.Fields{"service": "hq"})

    return c, nil
}
