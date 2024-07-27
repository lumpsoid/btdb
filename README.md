# Bluetooth for dual boot

Export bluetooth device keys from windows to linux.

## Usage

```sh
btdb -mnt /mnt
```

`/mnt` - where windows partition is mounted

```
Interface: <interface mac>

Device: <bluetooth < 5.1 mac>
  Key: <bluetooth key>

Device: <bluetooth 5.1 mac>
  LTK: <key>
  IRK: <key>
  ERand: <number>
  EDIV: <number>
```

After `btdb` printed your devices, you will need to go to the `/var/lib/bluetooth/<interface mac>/<bluetooth mac>/info` location, and edit specified keys.

For now, you will need to manually type printed keys on your linux machine.

## Build

```sh
git clone https://github.com/lumpsoid/btdb
cd btdb
go build -o btdb ./cmd/btdb/main.go
```
