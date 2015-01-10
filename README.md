driver-kodi
====================

Ninja Sphere - Kodi (Xbmc) Driver

To Do
- Notifications
- Media metadata
- Remote control? (up down left right etc...)
- Seeking
- Repeat + Shuffle
- Low battery alert
- Show playlist? Jump to item?

Bugs
- Sometimes the Kodi instance isn't discovered using ZeroConf. Add upnp?

## Building the driver

The follwoing instructions were collected from the discussion in  [#3](/../../issues/3). See the issue thread for more details.

### Requirements

#### GoLang
Install golang with cross compilation capabilities, since you will be compiling the driver for linx OS and arm architecture (this is of course not needed if you are buidling from a linux/arm machine).

Some pointers for installing golang with cross compilation capability:
- https://coderwall.com/p/pnfwxg/cross-compiling-golang
- http://dave.cheney.net/2012/09/08/an-introduction-to-cross-compilation-with-go

#### Loging in to the Spheramid

If your network automatically maps hotsnames to IPs, then you should be able to reach your spheramid via the `ninjasphere` hostname. Otherwise you need to lookup it's IP address.

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
```

In order to copy the drivers to the spheramid you need a drop-off point where the user ninja can write to.
In the spheramid:
```bash
sudo -i
mkdir /data/dropoff
chown ninja /data/dropoff
```

#### Install a text editor

VI comes installed in the spheramid, but VIM is easyer to use, so you can optionally install VIM via `apt-get`:
On ths spheramid:
```bash
sudo -i
apt-get -y install vim
```

### Compile

In the driver project directory:
```bash
make all
```

### Deploy the compiled driver

In the driver project directory:
```bash
scp bin/driver-kodi package.json ninja@ninjasphere:/data/dropoff
```

### Install the new driver

In the spheramid:
```bash
sudo -i
with-rw bash
mv /data/dropoff/* /opt/ninjablocks/drivers/driver-kodi/
vim /opt/ninjablocks/config/default.json
<add "driver-kodi" to the list of drivers in section homecloud.autostart -- don't forget the comma!>
```

