# Nagios Core 4.5.13 and Nagios Plugins 2.5 Build Commands for Debian

## Overview

This document describes the recommended build commands for:

* Nagios Core 4.5.13
* Nagios Plugins 2.5

on Debian systems, including both:

* Manual installation
* Debian package (.deb) creation

The commands follow Debian packaging best practices and avoid installing files directly into the operating system during package builds.

---

# Nagios Core 4.5.13

## Configure

```bash
./configure \
  --prefix=/usr \
  --sysconfdir=/etc/nagios \
  --localstatedir=/var/lib/nagios \
  --with-command-group=nagios
```

## Build

Single-threaded:

```bash
make all
```

Multi-threaded:

```bash
make -j$(nproc) all
```

Verify version:

```bash
src/nagios --version
```

---

# Manual Installation

## Install Binaries

```bash
make install
```

Installs:

```text
/usr/bin
/usr/sbin
/etc/nagios
/var/lib/nagios
```

---

## Install Default Configuration

```bash
make install-config
```

Creates:

```text
/etc/nagios/nagios.cfg
/etc/nagios/cgi.cfg
/etc/nagios/resource.cfg
```

---

## Install CGI Components

```bash
make install-webconf
```

Note:

This target is designed for Apache integration and is generally not used when deploying Nagios behind Nginx and fcgiwrap.

---

## Install Init Scripts

```bash
make install-daemoninit
```

Modern Debian systems should use systemd units instead.

---

# Debian Package Build (Recommended)

When creating .deb packages, avoid installing directly into the host system.

Use DESTDIR staging.

## Build

```bash
make -j$(nproc) all
```

## Install Into Package Staging Directory

```bash
make install DESTDIR=$(pwd)/debian/tmp
```

Install default configuration:

```bash
make install-config DESTDIR=$(pwd)/debian/tmp
```

Result:

```text
debian/tmp/usr/
debian/tmp/etc/
debian/tmp/var/
```

This directory is later processed by:

```bash
dpkg-buildpackage
```

---

# Nagios Plugins 2.5

## Configure

```bash
./configure \
  --prefix=/usr \
  --libexecdir=/usr/lib/nagios/plugins
```

---

## Build

```bash
make
```

or

```bash
make -j$(nproc)
```

---

## Manual Installation

```bash
make install
```

Plugins will be installed into:

```text
/usr/lib/nagios/plugins
```

---

# Debian Package Build

## Build

```bash
make -j$(nproc)
```

## Install Into Staging Directory

```bash
make install DESTDIR=$(pwd)/debian/tmp
```

Result:

```text
debian/tmp/usr/lib/nagios/plugins
```

---

# Available Make Targets

To view all available targets:

```bash
make help
```

or

```bash
grep "^install" Makefile
```

Common Nagios Core targets:

```text
install
install-config
install-commandmode
install-webconf
install-daemoninit
install-init
```

---

# Recommended Build Workflow for Debian Trixie

## Nagios Core

```bash
./configure \
  --prefix=/usr \
  --sysconfdir=/etc/nagios \
  --localstatedir=/var/lib/nagios \
  --with-command-group=nagios

make -j$(nproc)

make install DESTDIR=debian/tmp

make install-config DESTDIR=debian/tmp
```

---

## Nagios Plugins

```bash
./configure \
  --prefix=/usr \
  --libexecdir=/usr/lib/nagios/plugins

make -j$(nproc)

make install DESTDIR=debian/tmp
```

---

# Debian Package Generation

After staging all files:

```bash
dpkg-buildpackage -us -uc -b
```

Expected output:

```text
nagios-core_4.5.13-1_amd64.deb
nagios-plugins_2.5-1_amd64.deb
```

---

# Best Practices

* Always use DESTDIR when building packages.
* Never run make install directly during package creation.
* Use dpkg-buildpackage for final package generation.
* Validate packages using lintian.
* Build packages in clean environments using sbuild or pbuilder.
* Keep configuration files under /etc/nagios.
* Install plugins under /usr/lib/nagios/plugins.
* Use systemd instead of legacy init scripts.
* Prefer Nginx + fcgiwrap over Apache when building modern deployments.

---

# Validation Commands

Validate Nagios configuration:

```bash
nagios -v /etc/nagios/nagios.cfg
```

Validate package:

```bash
lintian ../nagios-core_4.5.13-1_amd64.changes
```

List installed plugins:

```bash
ls -la /usr/lib/nagios/plugins
```

Verify Nagios version:

```bash
nagios --version
```

