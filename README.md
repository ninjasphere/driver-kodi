# Ninja Sphere - Kodi Driver


[![Build status](https://badge.buildkite.com/b76d5dcb767cd874a8f5699f88e3888dcbc5c76d3198d12673.svg)](https://buildkite.com/ninja-blocks-inc/driver-kodi)
[![godoc](http://img.shields.io/badge/godoc-Reference-blue.svg)](https://godoc.org/github.com/ninjasphere/driver-kodi)
[![MIT License](https://img.shields.io/badge/license-MIT-yellow.svg)](LICENSE)
[![Ninja Sphere](https://img.shields.io/badge/built%20by-ninja%20blocks-lightgrey.svg)](http://ninjablocks.com)
[![Ninja Sphere](https://img.shields.io/badge/works%20with-ninja%20sphere-8f72e3.svg)](http://ninjablocks.com)

---


### Introduction
This is a Ninja Sphere driver for Kodi (previously known as Xbmc) that exposes any discovered instances as media devices, allowing control from the mobile apps and led controller.

### Supported Sphere Protocols

| Name | URI | Supported Events | Supported Methods |
| ------ | ------------- | ---- | ----------- |
| volume | [http://schema.ninjablocks.com/protocol/volume](https://github.com/ninjasphere/schemas/blob/master/protocol/volume.json) | set, volumeUp, volumeDown, mute, unmute, toggleMute | state |
| media-control | [http://schema.ninjablocks.com/protocol/media-control](https://github.com/ninjasphere/schemas/blob/master/protocol/media-control.json) | play, pause, next, previous  | playing, paused, stopped |
| battery | [http://schema.ninjablocks.com/protocol/battery](https://github.com/ninjasphere/schemas/blob/master/protocol/battery.json) |   | warning |

### To Do
* Notifications
* Media metadata
* Remote control? (up down left right etc...)
* Seeking
* Repeat + Shuffle
* Low battery alert
* Show playlist? Jump to item?

### Bugs
* Sometimes the Kodi instance isn't discovered using ZeroConf. Add upnp?

### Requirements

* Go 1.3
* Kodi (may work with recent Xbmc also)

### Dependencies

https://github.com/ninjasphere/kodi_jsonrpc
https://github.com/jonaz/mdns

### Building the driver

The follwoing instructions were collected from the discussion in  [#3](/../../issues/3). See the issue thread for more details.

### Requirements

#### GoLang
Install golang with cross compilation capabilities, since you will be compiling the driver for linux/arm (this is of course not needed if you are buidling from a linux/arm machine).

Some pointers for installing golang with cross compilation capability:
- https://coderwall.com/p/pnfwxg/cross-compiling-golang
- http://dave.cheney.net/2012/09/08/an-introduction-to-cross-compilation-with-go

#### Logging in to the Spheramid

Your Spheramid should be accessible using the hostname `ninjasphere` or `ninjasphere.local`. If it is not, you need to find its IP address from your router.

```bash
ssh ninja@ninjasphere
# type in password: <default pass is temppwd>
```

#### Prepare the Spheramid to receive the new driver

The spheramid root file system is mounted read-only. That means you won't be immediately able to deploy the driver in there.
You need to first enable read-write on the root file system, then create a directory for the new driver.

In the spheramid:
```bash
sudo -i
with-rw bash
mkdir /opt/ninjablocks/drivers/driver-kodi
chmod +w /opt/ninjablocks/drivers/driver-kodi
```
#### Compile

In the driver project directory:
```bash
GOOS=linux GOARCH=arm go build
```

#### Deploy the compiled driver

In the driver project directory:
```bash
scp driver-kodi package.json ninja@ninjasphere:/opt/ninjablocks/drivers/driver-kodi
```

#### Install the new driver

On the spheramid:
```bash
sudo -i
with-rw bash
nano /opt/ninjablocks/config/default.json
<add "driver-kodi" to the list of drivers in section homecloud.autostart -- don't forget the comma!>
```

### Options

The usual Ninja Sphere configuration and parameters apply, but these are the most useful during development.

* `--autostart` - Don't wait to be started by Ninja Sphere
* `--mqtt.host=HOST` - Override default mqtt host
* `--mqtt.port=PORT` - Override default mqtt host

### More Information

More information can be found on the [project site](http://github.com/ninjasphere/driver-go-kodi) or by visiting the Ninja Blocks [forums](https://discuss.ninjablocks.com).

### Contributing Changes

To contribute code changes to the project, please clone the repository and submit a pull-request ([What does that mean?](https://help.github.com/articles/using-pull-requests/)).

### License
This project is licensed under the MIT license, a copy of which can be found in the [LICENSE](LICENSE) file.

### Copyright
This work is Copyright (c) 2014-2015 - Ninja Blocks Inc.
