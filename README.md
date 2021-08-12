# geoip
A microservice that performs a Geo lookup on an IP or a PFSense Firewall Filter log entry and sends the results to an 
InfluxDB instance.

## Pre-Reqs
1. [Maxmind Geo-Lite](https://dev.maxmind.com/geoip/geolite2-free-geolocation-data) account and license key
2. Running instance of InfluxDB 2.0 

## Building and Installing
Makefile Targets:

| Target | Description|
|---|---|
|deps | Download and Install any missing dependecies|
|build  | Install missing dependencies. Builds binaries for linux and darwin in ./dist |
|dist | Creates a distribution  |
|tidy |                    Verifies and downloads all required dependencies|
|fmt   |                   Runs gofmt on all source files|
|test   |                  Tests code coverage|
|testwithcoverge |         Tests code coverage|
|missing         |         Displays lines of code missing from coverage. Puts report in ./build/coverage.out|
|vet  |                    Run go vet.  Puts report in ./build/vet.out|
|reports |                 Runs vet, coverage, and missing reports|
|clean |                   Removes build, dist and report dirs|
|gencerts |                Generates a sample self signed cert and key to enable TLS|
|debug  |                  Print make env information|

## Configuring
By default the binary will look for a config file named config.toml in the following locations:<br/>
```go
viper.AddConfigPath(".")                          // Local dir
viper.AddConfigPath("/etc/" + version.APP_NAME)   // Looks in etc/geoip
viper.AddConfigPath("$HOME/." + version.APP_NAME) // Looks in $HOME/geoip
```

The location and name of the config file can be overridden.  See below for details.

The file config.toml contains a default config file.  A default config file can also be generated.  See below for 
details.

```toml
# Print the current config to console
printConfig = true 

# Logging Config
[logging]
    # Logging level TRACE if true, INFO if false
    tracelogging = false

    # name of the logfile.  Will also log to stdout if run from the cli
    logfile = "/var/log/geoip/geoip.log"

    # HTTP Logging
    [logging.http]
        # enable http logging
        enabled = true
    
        # log http to stdout
        stdout = true
    
        # Log http to file
        fileout = false
    
        # Name fo the Http logfile
        logfile = "/var/log/geoip/geoip_http.log"

# Http server config
[http]
    # IP to bind to and listen
    hostname = "0.0.0.0"

    # Port to bind to and listen
    port = "8082"
    
    # Use HTTPS
    useHttps = false

    # Minimum TLS version to use
    tlsMinVersion = "1.2"

    # Use Strict TLS Ciphers
    httpTLSStrictCiphers = false

    # TLS Certificate file to use
    tlsCert = "/etc/ssl/geoip.crt"

    # TLS Key file to use
    tlsKey = "/etc/ssl/geoip.key"

    # Enable CORS support
    enableCORS = false

    # JWT Secret for protected  api endpoints.  Currently not used
    jwtSecret = "There is a mouse in my house"

# GEO IP Lookup Config
[geoip]

# How often the MaxMind DB should be refreshed from MaxMind.
# See: https://dev.maxmind.com/geoip/updating-databases?lang=en
RefreshDuration = "24h"

# Your Maxmind account id
accountId = 123

# Your MaxMind license Key
LicenseKey = "testing"

# Directory to store the Maxmind geo database
DatabaseDirectory = "data/geoip"

# Name of the MaxMind geo database
DatabaseName = "GeoLite2-City.mmdb"

# verbose output for database updating
Verbose = false

# InfluxDB config
[influxdb]
# URL of the InfluxV2 database 
url = "http://localhost:8086"

# InfluxV2 token
token = "mytoken"

# InfluxV2 org name
org = "myOrg"

# InfluxV2 bucket name
bucket = "myBucket"
```


## How to Use
Available Commands:

| command | description |
| --- | --- |
| config | Generates and prints a default config |
| help | Help about any command |
| serve | start http server with configured api |
| version | Print the version number |

**Generate a default config file**
```shell
./geoip config > config.toml
```

**Serve with default config**
```shell
./geoip serve
```

**Serve while specifying path to config file**
```shell
./geoip serve -c /path/to/config.toml
```

## API
<span style="color:#49cc90; font-weight: bold;">POST</span> -- /v1/geo/batch <br>
For a batch of line protocol from tail will perform a geoip lookup and send a new point to InfluxDB<br/>

**body:**<br>
<table>
<thead>
<td>Name</td>
<td>Type</td>
<td>Required</td>
<td>Description</td>
</thead>
<tr>
<td>body</td>
<td>object</td>
<td>yes</td>
<td>Batch of line protocol.  Where each line, at a minimum must contain the following:<br/>
<b><u>Tags:</u></b><br/>
<li>host -- The name of the host</li>

<b><u>Fields:</u></b><br/>
<li>dest_ip -- The ip address of the destination </li>
<li>dest_port -- The port of the destination </li>
<li>src_ip -- The ip address of the source </li>
<li>src_port -- The port of the source </li>
<li>rule -- The rule that triggered the block </li>
<li>direction -- The direction of the packet (in or out)</li>
</td>
</tr>
</table>

***Example:***<br/>
```
tail,host=localhost,path=filter.log dest_ip="173.160.205.9",iface="em0",direction="in",src_port="14553",rule="102",
src_ip="205.185.117.79",dest_port="22" 1625793903302961350
```

**Response:**
<table>
<thead>
<td>Code</td>
<td>Description</td>
</thead>
<tr>
<td>201</td>
<td>created</td>
</tr>
<tr>
<td>500</td>
<td>error</td>
</tr>
</table>
<br/>

<span style="color:#61affe; font-weight: bold;">GET</span>/v1/geo/{ipaddress}<br/>
Returns JSON GeoLocation of the ip address

**Parameters:**<br>

| Name | Type | Required |Description|
|---|---|---|---|
| ipaddress | string | yes | The ipaddress to lookup |

**Response:**<br/>
<table>
<thead>
<td>Code</td>
<td>Description</td>
</thead>
<tr>
<td>200</td>
<td>
JSON of GeoLocation<br>
<pre>
{
    "City":
    {
        "GeoNameID": 0,
        "Names": null
    },
    "Continent":
    {
        "Code": "NA",
        "GeoNameID": 6255149,
        "Names":
        {
            "de": "Nordamerika",
            "en": "North America",
            "es": "Norteamérica",
            "fr": "Amérique du Nord",
            "ja": "北アメリカ",
            "pt-BR": "América do Norte",
            "ru": "Северная Америка",
            "zh-CN": "北美洲"
        }
    },
    "Country":
    {
        "GeoNameID": 6252001,
        "IsInEuropeanUnion": false,
        "IsoCode": "US",
        "Names":
        {
            "de": "USA",
            "en": "United States",
            "es": "Estados Unidos",
            "fr": "États-Unis",
            "ja": "アメリカ合衆国",
            "pt-BR": "Estados Unidos",
            "ru": "США",
            "zh-CN": "美国"
        }
    },
    "Location":
    {
        "AccuracyRadius": 1000,
        "Latitude": 37.751,
        "Longitude": -97.822,
        "MetroCode": 0,
        "TimeZone": "America/Chicago"
    },
    "Postal":
    {
        "Code": ""
    },
    "RegisteredCountry":
    {
        "GeoNameID": 6252001,
        "IsInEuropeanUnion": false,
        "IsoCode": "US",
        "Names":
        {
            "de": "USA",
            "en": "United States",
            "es": "Estados Unidos",
            "fr": "États-Unis",
            "ja": "アメリカ合衆国",
            "pt-BR": "Estados Unidos",
            "ru": "США",
            "zh-CN": "美国"
        }
    },
    "RepresentedCountry":
    {
        "GeoNameID": 0,
        "IsInEuropeanUnion": false,
        "IsoCode": "",
        "Names": null,
        "Type": ""
    },
    "Subdivisions": null,
    "Traits":
    {
        "IsAnonymousProxy": false,
        "IsSatelliteProvider": false
    }
}
</pre>
</td>
</tr>
<tr>
<td>500</td>
<td>error</td>
</tr>
</table>
<br/>

<span style="color:#61affe; font-weight: bold;">GET</span> /v1/api<br/>
Returns JSON of API

**Response:**
<table>
<thead>
<td>Code</td>
<td>Description</td>
</thead>
<tr>
<td>200</td>
<td>
JSON of the api<br>
<pre>
[
    "POST /v1/geo/batch",
    "GET /v1/geo/{ipaddress}",
    "GET /v1/api",
    "GET /v1/ping",
    "GET /v1/version"
]
</pre>
</td>
</tr>
</table>
<br/>

<span style="color:#61affe; font-weight: bold;">GET</span> /v1/ping

**Response:**
<table>
<thead>
<td>Code</td>
<td>Description</td>
</thead>
<tr>
<td>200</td>
<td>
OK
</td>
</tr>
</table>
<br/>

<span style="color:#61affe; font-weight: bold;">GET</span> /v1/version

**Response:**
<table>
<thead>
<td>Code</td>
<td>Description</td>
</thead>
<tr>
<td>200</td>
<td>
JSON of version information in the form of:
<pre>
{
    "AppName":"geoip",
    "Version":"dev",
    "Branch":"main",
    "Commit":"b1fa4ada42cf92d156593a1517e34d9540c81cea",
    "BuildTime":"2021-07-27T15:51:09-07:00"
}
</pre>
</td>
</tr>
</table>
<br/>