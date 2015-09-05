package main

import (
	"flag"
	"fmt"
	//	"github.com/spf13/cast"
	"github.com/spf13/viper"
	"gopkg.in/ldap.v1"
	"log"
	"net"
	"time"
)

const APP_VERSION = "0.1"

// The flag package provides a default help printer via -h switch
var versionFlag *bool = flag.Bool("v", false, "Print the version number.")

type LdapEntry struct {
	URI  string
	user string
	pass string
}

type Ldap []*LdapEntry
type Ldaps map[string]*Ldap

// Initialize Configuration
func InitializeConfig() {
	fmt.Println("read config files")
	viper.SetConfigName("config")
	viper.AddConfigPath("/etc/cnmonitor2graphite/")
	viper.AddConfigPath("$HOME/.cnmonitor2graphite")
	viper.Set("Verbose", true)
	viper.SetConfigType("toml")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
}

//Fetch Data from ldap
func FetchData(ldapuri, user, passwd, baseDN, filter string, Attributes []string) {
	log.Println("Connect to ldap and read data")
	l, err := ldap.Dial("tcp", ldapuri)
	if err != nil {
		log.Fatalf("ERROR: %s\n", err.Error())
	}
	defer l.Close()
	// l.Debug = true

	err = l.Bind(user, passwd)
	if err != nil {
		log.Printf("ERROR: Cannot bind: %s\n", err.Error())
		return
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
		return
	}

	log.Printf("BaseDN: %s  -- Search: %s -> num of entries = %d\n", baseDN, search.Filter, len(sr.Entries))
	sr.PrettyPrint(0)
}

//Sent data to graphite
func SentData(graphite, prefix string) {
	t := time.Now().Unix()
	fmt.Println("Connect to graphite and sent data")
	conn, err := net.Dial("tcp", graphite)
	if err != nil {
		log.Printf("ERROR: Cannot connect: %s\n", err.Error())
		return
	}
	defer conn.Close()
	fmt.Fprintf(conn, "%s.test 20 %d\n", prefix, t)
	fmt.Printf("%s.test 20 %d\n", prefix, t)
}

func main() {
	flag.Parse() // Scan the arguments list

	if *versionFlag {
		fmt.Println("Version:", APP_VERSION)
	}
	InitializeConfig()
	graphite := fmt.Sprintf("%s%s%s", viper.GetString("graphite.host"), ":", viper.GetString("graphite.port"))
	prefix := fmt.Sprintf("%s", viper.GetString("graphite.prefix"))
	fmt.Println("Graphite >>>", graphite)
	fmt.Println("Prefix >>>", prefix)
	SentData(graphite, prefix)
	ldapmap := viper.GetStringMap("ldap")
	dnmap := viper.GetStringMap("dn")
	for ldap, _ := range ldapmap {
		ldapuri := viper.GetString(fmt.Sprintf("ldap.%s.uri", ldap))
		ldapuser := viper.GetString(fmt.Sprintf("ldap.%s.user", ldap))
		ldappass := viper.GetString(fmt.Sprintf("ldap.%s.pass", ldap))
		for dn, _ := range dnmap {
			fmt.Println("dn:", dn)
			data := viper.GetStringSlice(fmt.Sprintf("dn.%s.data", dn))
			basedn := viper.GetString(fmt.Sprintf("dn.%s.dn", dn))
			FetchData(ldapuri, ldapuser, ldappass, basedn, "(objectclass=*)", data)
		}
	}
}
