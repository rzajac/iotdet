package cmd

import (
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
    "regexp"
    "fmt"
)

// cfgFile holds path to the configuration file.
var cfgFile string

func init() {
    cobra.OnInitialize(onInitialize)
    rootCmd.SetVersionTemplate(`{{.Version}}`)
    rootCmd.Flags().StringVarP(&cfgFile, "config", "c", "", "path to configuration file (default is ./iotdet.yaml)")
    rootCmd.Flags().BoolP("version", "v", false, "version")
    rootCmd.Flags().BoolP("debug", "d", false, "run in debug mode")
    viper.BindPFlag("debug", rootCmd.Flags().Lookup("debug"))
}

// Execute executes root command.
func Execute() {
    if err := rootCmd.Execute(); err != nil {
        fmt.Fprintf(os.Stderr, "%s\n", err.Error())
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
func getHQ() (*hq.HQ, error) {
    c := hq.NewHQ()

    // configure new agent detection.
    detItfName := viper.GetString("hq.detect.itf")
    detApNamePat := regexp.MustCompile(viper.GetString("hq.detect.agent_ap_name"))
    detApPass := viper.GetString("hq.detect.ap_pass")
    if err := c.SetDet(detItfName, detApNamePat, detApPass); err != nil {
        return nil, err
    }

    // Set IPs to use during the detection stage.
    detAgentIP := viper.GetString("hq.detect.agent_ip")
    detUseIP := viper.GetString("hq.detect.use_ip")
    if err := c.SetDetIPs(detUseIP, detAgentIP); err != nil {
        return nil, err
    }

    // Set TCP port agents use for command server.
    detCmdPort := viper.GetInt("hq.detect.cmd_port")
    if err := c.SetDetCmdPort(detCmdPort); err != nil {
        return nil, err
    }

    // Set detection interval.
    detInterval := viper.GetDuration("hq.detect.scan_interval") * time.Second
    if err := c.SetDetInterval(detInterval); err != nil {
        return nil, err
    }

    // Set cipher configuration for agents communication.
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

        c.SetCipher(hq.NewCipherAES(key, vi))

    case hq.CIPHER_NONE:
        c.SetCipher(hq.NewNoopCipher())

    default:
        return nil, errors.Errorf("invalid cipher name %s", ci)
    }

    // Set main access point configuration.
    apName := viper.GetString("hq.access_point.name")
    apPass := viper.GetString("hq.access_point.pass")
    if err := c.SetAccessPoint(apName, apPass); err != nil {
        return nil, err
    }

    // Set MQTT configuration.
    clientID := viper.GetString("hq.mqtt.client_id")
    mqttIP := viper.GetString("hq.mqtt.ip")
    mqttPort := viper.GetInt("hq.mqtt.port")
    mqttUser := viper.GetString("hq.mqtt.user")
    mqttPass := viper.GetString("hq.mqtt.pass")

    mqttCfg, err := hq.NewMQTTConfig(clientID, mqttIP, mqttPort, mqttUser, mqttPass)
    if err != nil {
        return nil, err
    }

    c.SetMQTTConfig(mqttCfg)

    // Setup logger.
    if viper.GetBool("debug") {
        c.SetLogger(hq.NewLog().DebugOn())
    }

    return c, nil
}
