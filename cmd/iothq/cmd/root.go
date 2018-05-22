package cmd

import (
    "github.com/spf13/viper"
    "strings"
    "github.com/spf13/cobra"
    "github.com/rzajac/iotdet/version"
    "encoding/json"
    "os"
    "encoding/hex"
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
    rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "path to configuration file (default is ./iotdet.yaml)")
    rootCmd.PersistentFlags().BoolP("version", "v", false, "version")
    rootCmd.PersistentFlags().BoolP("debug", "d", false, "run in debug mode")
    viper.BindPFlag("debug", rootCmd.PersistentFlags().Lookup("debug"))
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
    // name of the configuration file and where to look for it.
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

// getConfiguredHQ returns configured HQ structure.
func getConfiguredHQ() (*hq.HQ, error) {
    h, err := hq.NewHQ(config{})
    if err != nil {
        return nil, err
    }

    hq.ERROR = &Log{}
    hq.INFO = &Log{}

    // Setup logger.
    if viper.GetBool("debug") {
        hq.DEBUG = &Log{}
    }

    return h, nil
}

// config implements hq.Configurator interface.
type config struct{}

func (config) GetDetItfName() string {
    return viper.GetString("hq.detect.itf")
}

func (config) GetDetItfIP() string {
    return viper.GetString("hq.detect.itf_ip")
}

func (config) GetDetInterval() time.Duration {
    return viper.GetDuration("hq.detect.scan_interval") * time.Second
}

func (config) GetBeaconNamePat() *regexp.Regexp {
    return regexp.MustCompile(viper.GetString("hq.detect.beacon_name"))
}

func (config) GetBeaconPass() string {
    return viper.GetString("hq.detect.beacon_pass")
}

func (config) GetBeaconIP() string {
    return viper.GetString("hq.detect.beacon_ip")
}

func (config) GetBeaconPort() int {
    return viper.GetInt("hq.detect.cmd_port")
}

func (config) GetCipher() hq.Cipher {
    ci := viper.GetString("hq.cipher")
    switch ci {
    case hq.CipherAES:
        key, err := hex.DecodeString(viper.GetString("hq.cipher_aes.key"))
        if err != nil {
            panic("invalid AES key value")
        }

        vi, err := hex.DecodeString(viper.GetString("hq.cipher_aes.vi"))
        if err != nil {
            panic("invalid AES vi value")
        }
        return hq.NewCipherAES(key, vi)

    case hq.CipherNoop:
        return hq.NewNoopCipher()

    default:
        panic("invalid cipher name " + ci)
    }
}

func (config) GetAPName() string {
    return viper.GetString("hq.access_point.name")
}

func (config) GetAPPass() string {
    return viper.GetString("hq.access_point.pass")
}

func (config) GetMQTTIP() string {
    return viper.GetString("hq.mqtt.ip")
}

func (config) GetMQTTPort() int {
    return viper.GetInt("hq.mqtt.port")
}

func (config) GetMQTTUser() string {
    return viper.GetString("hq.mqtt.user")
}

func (config) GetMQTTPass() string {
    return viper.GetString("hq.mqtt.pass")
}

func (config) GetMQTTClientID() string {
    return viper.GetString("hq.mqtt.client_id")
}

type Log struct{}

func (l *Log) Println(v ...interface{}) {
    fmt.Println(v...)
}

func (l *Log) Printf(format string, v ...interface{}) {
    fmt.Printf(format, v...)
}
