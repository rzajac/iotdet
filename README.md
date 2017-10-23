## Detect and configure IoT devices.

Repository contains library and command line tool for automatic detection and 
configuration of wireless IoT devices in range.

It was designed to work with ESP8266 C implementation 
at https://github.com/rzajac/esp-det repository. 

##  Configuration.

See [iotdet.ini](iotdet/iotdet.ini) file.

## AES key generation.

You may use [gen_aes.sh](gen_aes.sh) helper script to generate keys using `openssl rand`.

## TODO

- Remove configuration file written by wpa_supplicant to temp directory.
- Check wpa_supplicant command log and proceed only when message 
like "CTRL-EVENT-CONNECTED - Connection to 5e:cf:7f:80:ce:79 completed" is seen.

## License

[Apache License Version 2.0](LICENSE) unless stated otherwise.
