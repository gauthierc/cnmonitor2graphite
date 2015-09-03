package main 

import (
    "flag"
    "github.com/spf13/viper"
    "fmt"
)

const APP_VERSION = "0.1"

// The flag package provides a default help printer via -h switch
var versionFlag *bool = flag.Bool("v", false, "Print the version number.")

func main() {
    flag.Parse() // Scan the arguments list 

    if *versionFlag {
        fmt.Println("Version:", APP_VERSION)
    }
    viper.SetConfigName("config")
    viper.AddConfigPath("/etc/cnmonitor2graphite/")
    viper.AddConfigPath("$HOME/.cnmonitor2graphite")
    viper.Set("Verbose", true)
    viper.SetConfigType("toml")
    err := viper.ReadInConfig()
    if err != nil {
    	panic(fmt.Errorf("Fatal error config file: %s \n", err))
    }
    fmt.Println("Graphite >>>", viper.GetString("graphite.host"),":",viper.GetString("graphite.port"))
    fmt.Println("dn >>>", viper.GetString("snmp.dn"))
    data := viper.GetStringSlice("snmp.data")
    fmt.Println("data >>>", data)
    
    for i := 0; i < len(data); i++ {
    	fmt.Println(data[i])
    }
}

