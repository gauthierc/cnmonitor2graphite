cnmonitor2graphite
==================
cnmonitor2graphite est un programme qui interroge un ou plusieurs serveurs ldap et qui envoie les informations de connexions sur un serveur graphite.
## configuration
Copier le fichier config.toml.exemple dans /etc/cnmonitor2graphite/ ou dans le répertoire $HOME/.cnmonitor2graphite/ de l'utilisateur et renommer le en config.toml
### Graphite
Remplir la partie concernant votre serveur graphite:
```
[graphite]
host = "graphite"
port = "2003"
prefix = "cnmonitor"
```
### ldap
Il est possible de renseigner plusieurs serveurs ldap comme ceci:

```
[ldap.monserveurldap]
	uri = "ldap://ldap1.mydomain.io"
    user = "cn=Directory Manager"
    pass = "secret"
```
Le champs user et pass peuvent être vide si la ressource ldap est accessible en anonyme.

### dn
Cette partie indique le chemin et les attributs à grapher.

```
[dn.snmp]
        dn = "cn=snmp,cn=monitor"
        data = [
                "anonymousbinds",
                "unauthbinds",
                "simpleauthbinds",
                "strongauthbinds",
                "bindsecurityerrors",
                "inops",
                "readops",
                ...
```

### Whisper
Organisation des données dans Whisper avec une configuration du style :
```
[Graphite]
prefix = "cnmonitor"
...
[ldap.monserveurldap]
...
[dn.snmp]
        dn = "cn=snmp,cn=monitor"
        data = [
				"anonymousbinds",
                "simpleauthbinds",
```
Seront stockés dans Whisper comme ceci:
```
cnmonitor.monserveurldap.snmp.anonymousbinds
cnmonitor.monserveurldap.snmp.simpleauthbinds
```

## Utilisation
Lancer le programme dans un cron toutes les minutes.
