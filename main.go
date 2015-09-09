package main

import (
	"flag"
	"fmt"
	"github.com/spf13/viper"
	"gopkg.in/ldap.v1"
	"log"
	"net"
	"time"
)

const APP_VERSION = "0.5"
const DEBUG = false

// The flag package provides a default help printer via -h switch
var versionFlag *bool = flag.Bool("v", false, "Print the version number.")

// Initialize Configuration
func InitializeConfig() {
	viper.SetConfigName("config")
	viper.AddConfigPath("/etc/cnmonitor2graphite/")
	viper.AddConfigPath("$HOME/.cnmonitor2graphite/")
	viper.Set("Verbose", true)
	viper.SetConfigType("toml")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
}

//Fetch Data from ldap
func FetchData(ldapuri, user, passwd, baseDN, filter string, Attributes []string) *ldap.SearchResult {
	if DEBUG {
		log.Println("Connect to ldap and read data")
	}
	l, err := ldap.Dial("tcp", ldapuri)
	if err != nil {
		log.Fatalf("ERROR: %s\n", err.Error())
	}
	defer l.Close()
	// l.Debug = true

	err = l.Bind(user, passwd)
	if err != nil {
		log.Printf("ERROR: Cannot bind: %s\n", err.Error())
		return nil
	}
	search := ldap.NewSearchRequest(
		baseDN,
		ldap.ScopeBaseObject, ldap.NeverDerefAliases, 0, 0, false,
		filter,
		Attributes,
		nil)

	sr, err := l.Search(search)
	if err != nil {
		log.Fatalf("ERROR: %s\n", err.Error())
		return nil
	}

	if DEBUG {
		log.Printf("BaseDN: %s  -- Search: %s -> num of entries = %d\n", baseDN, search.Filter, len(sr.Entries))
	}
	return sr
}

//Sent data to graphite
func SentData(graphite, prefix string, result *ldap.SearchResult) {
	t := time.Now().Unix()
	if DEBUG {
		log.Println("Connect to", graphite, "and sent data")
	}
	conn, err := net.Dial("tcp", graphite)
	if err != nil {
		if DEBUG {
			log.Printf("ERROR: Cannot connect: %s\n", err.Error())
		}
		return
	}
	defer conn.Close()
	for _, entry := range result.Entries {
		for _, attr := range entry.Attributes {
			fmt.Fprintf(conn, "%s.%s %s %d\n", prefix, attr.Name, attr.ByteValues[0], t)
		}
	}
}

//Show data for graphite
func ShowData(graphite, prefix string, result *ldap.SearchResult) {
	t := time.Now().Unix()
	for _, entry := range result.Entries {
		for _, attr := range entry.Attributes {
			fmt.Printf("%s.%s %s %d\n", prefix, attr.Name, attr.ByteValues[0], t)
		}
	}
}

func main() {
	flag.Parse() // Scan the arguments list
	InitializeConfig()
	graphite := fmt.Sprintf("%s%s%s", viper.GetString("graphite.host"), ":", viper.GetString("graphite.port"))
	prefix := fmt.Sprintf("%s", viper.GetString("graphite.prefix"))

	if *versionFlag {
		fmt.Println("cnm2g: Cn=monitor to Graphite")
		fmt.Println("Version:", APP_VERSION)
		fmt.Println("Config File >>>", viper.ConfigFileUsed())
		fmt.Println("Graphite >>>", graphite)
		fmt.Println("Prefix >>>", prefix)
		return
	}
	ldapmap := viper.GetStringMap("ldap")
	dnmap := viper.GetStringMap("dn")
	for ldap, _ := range ldapmap {
		ldapuri := viper.GetString(fmt.Sprintf("ldap.%s.uri", ldap))
		ldapuser := viper.GetString(fmt.Sprintf("ldap.%s.user", ldap))
		ldappass := viper.GetString(fmt.Sprintf("ldap.%s.pass", ldap))
		for dn, _ := range dnmap {
			prefixldap := fmt.Sprintf("%s.%s.%s", prefix, ldap, dn)
			data := viper.GetStringSlice(fmt.Sprintf("dn.%s.data", dn))
			basedn := viper.GetString(fmt.Sprintf("dn.%s.dn", dn))
			ldapresult := FetchData(ldapuri, ldapuser, ldappass, basedn, "(objectclass=*)", data)
			if DEBUG {
				ShowData(graphite, prefixldap, ldapresult)
			} else {
				SentData(graphite, prefixldap, ldapresult)
			}
		}
	}
}
