## Detect and configure IoT wireless devices.

Repository contains library and command line tool for automatic detection and 
configuration of wireless IoT devices, called agents, in range.

It was designed to work with ESP8266 C implementation 
at https://github.com/rzajac/esp-det repository. 

## Help.

```
IoT HQ.
   
   Usage:
     iothq [command]
   
   Available Commands:
     configure   Configure agent.
     detect      Detect new agents
     help        Help about any command
   
   Flags:
     -c, --config string   path to configuration file (default is ./iothq.yaml)
     -d, --debug           run in debug mode
     -h, --help            help for iothq
     -v, --version         version
   
   Use "iothq [command] --help" for more information about a command.
```

##  Configuration.

See [hq.yaml](cmd/iothq/hq.yaml) file.

## AES key generation.

You may use [gen_aes.sh](gen_aes.sh) helper script to generate keys 
using `openssl rand`.

## License

[Apache License Version 2.0](LICENSE) unless stated otherwise.
