# venue

Golang API and server for TouchOSC to control Avid™ VENUE software

- <https://www.avid.com/US/products/venue-software>
- <http://hexler.net/software/touchosc>

This software package enables audio engineers to control an Avid™ VENUE system
using a mobile device. The initial version focuses on the workflow of a
monitoring engineer needing to quickly set and maintain levels for multiple
monitor mixes, while having the freedom to stand on the stage with the
performers while doing so.

Although an engineer can always use a VNC client to perform this function,
doing so is cumbersome due to the small UI elements and lack of design for
mobile usage. This software should make that easier.

## Project links

[![Go Reference](https://pkg.go.dev/badge/github.com/kward/venue.svg)](https://pkg.go.dev/github.com/kward/venue)
[![CI](https://github.com/kward/venue/actions/workflows/ci.yml/badge.svg)](https://github.com/kward/venue/actions/workflows/ci.yml)
[![Coverage Status](https://coveralls.io/repos/github/kward/venue/badge.svg?branch=master)](https://coveralls.io/github/kward/venue?branch=master)

Design doc: <https://docs.google.com/document/d/1RZKLrLeaKxsrjHwHxDp2NOQgQHE1-Px-9bGcWrq0g6E/pub>

## Requirements

This software requires ONE of the following:

- A real VENUE console, with remote access via Ethernet (e.g. via an add-on
  Ethernet card).
- A Windows machine running the free [Avid VENUE Software][VENUE], and a VNC
  server to export the display. The Windows instance can easily be run as a
  virtual machine so that a separate physical machine is not required.

Software was tested on the following:

- An [Avid VENUE Profile][Profile] System.
- Windows 7, running on [VMware Fusion][Fusion] (which provides a built-in VNC
  server) on OS X Yosemite.

## Setup

These instructions assume no knowledge of writing software with the Go language.
If you have  experience, feel free to follow your preferred standards.

### Installation

This code is written in Golang (<http://golang.org/>).

1. Install Golang. Follow the instructions at <http://golang.org/doc/install>.
2. Setup environment. Note, the exports must either be run each time the
   software will be used, or they can be added to your `~/.bashrc` file.
   (Examples are for OS X or Linux.)

    ```sh
    $ mkdir -p "${HOME}/opt/go/bin"
    $ export GOPATH="${HOME}/opt/go"
    $ export GOBIN="${GOPATH}/bin"
    $ export PATH="${PATH}:/usr/local/go/bin:${GOBIN}"
    ```

3. Download software.

    ```sh
  $ go get github.com/kward/venue
  $ go get github.com/kward/go-osc
  $ go get github.com/kward/go-vnc
  $ go get github.com/golang/glog
    ```

4. Test the client software. This will "randomly" select an input channel every
   few seconds. It is simply to test that a connection can be made and that the
   console can be controlled.

    ```sh
    $ cd "${GOPATH}/src/github.com/kward/venue"
    $ go run client/rand_inputs/main.go --venue_host <hostname/IP> --venue_passwd <passwd>
    Press CTRL-C to exit.
    ```

5. Install the TouchOSC layout. TODO(kward): Document this.

6. Test the server software. Configure TouchOSC to connect to the hostname/IP
   of your machine (not the host running VENUE).

   ```sh
   $ go run venue.go --venue_host <hostname/IP> --venue_passwd <passwd>
   ```

Notes:

- If you are not using the default VNC port of 5900, the `--venue_port` option
  should be added.
- If you are not using the default TouchOSC port of 8000, the
  `--osc_server_port` option should be added for `venue.go`.

### Updates

To update the software, repeat Installation step 3, with slight modifications.
Here's a simple script to do the updates.

```sh
$ for pkg in \
  github.com/kward/{go-osc,go-vnc,venue}
  github.com/golang/glog
do
  go get -u ${pkg}
done
```

_Avid™ is a registered trademark of Avid, Inc._

<!--- Links -->
[Fusion]: http://www.vmware.com/products/fusion/
[Profile]: https://www.avid.com/US/products/profile-system
[VENUE]: http://www.avid.com/us/products/venue-software

