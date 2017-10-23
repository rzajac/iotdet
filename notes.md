## Useful commands

```
wpa_passphrase IOT_5CCF7F80CE79 password > /tmp/iot_wpa_supplicant.conf

wpa_supplicant -D nl80211 -i wlx000f55a93e30 -c /tmp/iot_wpa_supplicant.conf
wpa_supplicant -h # show drivers

echo -n '{"cmd": "setAp", "name": "AccessPoint", "pass": "password"}' | nc 192.168.42.1 7802
echo -n '{"cmd": "setSrv", "ip": "192.168.1.149", "port": 1883,  "user": "username", "pass": "password"}' | nc 192.168.42.1 7802

nc -u -l -k 7802

```
