# IPTracker

This small application use this [http://ifconfig.me](http://ifconfig.me) service to discover your current IP.
Logs are organized in `YYYY/MM/DD.log` structure.

Example:
```bash
~/ipmon$ tree
.
└── 2017
    └── 12
        ├── 06.log
        └── 07.log
```

Log structure:
```
2017-12-07 00:20:30 |ERROR| Get http://ifconfig.me/all.json: read tcp 192.168.0.101:51010->153.121.72.212:80: read: connection reset by peer
2017-12-07 00:25:01 |INFO| 192.168.1.1
```
# Usage

Using other services that return json with `ip_addr` field is possible. However field name is fixed.

```bash
$ iptracker -help
Usage of iptracker:
  -path string
        path to save logs (default "$HOME/ipmon")
  -url string
        URL to get IP address (default "http://ifconfig.me/all.json")
```

# Installation

Default installation will build and move binary to `$HOME/bin` folder. Use with care.

One could install this application using:
```
make install
```

To build only use:
```
make build
```
Output is saved in `./bin` directory.

# Example cron configuration

Use command `crontab -e`. You *NEED* to specify full path to binary. With this configuration cron will run every 5 minutes.

```config
*/5 * * * * /<your_full_path_to_home_dir>/bin/iptracker
```