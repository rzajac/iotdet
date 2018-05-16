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

// log provides global logger.
var log *logrus.Entry

func init() {
    logrus.SetFormatter(&logrus.JSONFormatter{TimestampFormat: "2006-01-02T15:04:05.999999-07:00"})
    log = logrus.WithFields(logrus.Fields{"service": "iotdet"})

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
        log.Error(err)
        os.Exit(1)
    }
}

// rootCmd is the main command for the iotdet binary.
var rootCmd = &cobra.Command{
    Use:     "iothq",
    Version: getVersion(),
    Short:   "IoT HQ.",
    Long:    `IoT HQ.`,
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
    cfg := &hq.Config{}

    // New agent detection configuration.
    cfg.DetItfName = viper.GetString("hq.detect.itf")
    if cfg.DetItfName == "" {
        return nil, errors.New("you must provide WiFi interface name")
    }
    cfg.DetApPass = viper.GetString("hq.detect.ap_pass")
    cfg.DetAgentIP = viper.GetString("hq.detect.agent_ip")
    cfg.DetUseIP = viper.GetString("hq.detect.use_ip")
    cfg.DetCmdPort = viper.GetInt("hq.detect.cmd_port")
    cfg.DetInterval = viper.GetDuration("hq.detect.scan_interval") * time.Second

    // Cipher configuration.
    cn := viper.GetString("hq.cipher")
    switch cn {
    case "aes":
        key, err := hex.DecodeString(viper.GetString("hq.cipher_aes.key"))
        if err != nil {
            return nil, errors.New("invalid AES key value")
        }

        vi, err := hex.DecodeString(viper.GetString("hq.cipher_aes.vi"))
        if err != nil {
            return nil, errors.New("invalid AES vi value")
        }
        cfg.Cipher = hq.NewAesDrv(key, vi)

    case "none":
        cfg.Cipher = &hq.Noop{}

    default:
        return nil, errors.Errorf("invalid cipher name %s", cn)
    }

    // HQ access point configuration.
    cfg.HQApName = viper.GetString("hq.access_point.name")
    cfg.HQApPass = viper.GetString("hq.access_point.pass")

    // MQTT configuration.
    cfg.MQTTIP = viper.GetString("hq.mqtt.ip")
    cfg.MQTTPort = viper.GetInt("hq.mqtt.port")
    cfg.MQTTUser = viper.GetString("hq.mqtt.user")
    cfg.MQTTPass = viper.GetString("hq.mqtt.password")

    // Setup logger.
    cfg.Log = log

    return cfg, nil
}
