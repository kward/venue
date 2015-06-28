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

## Requirements
This software requires ONE of the following:

- A real VENUE console, with remote access via Ethernet (e.g. via an add-on
  Ethernet card).
- A Windows machine running the free [Avid VENUE | Software][VENUE], and a VNC
  server to export the display. The Windows instance can easily be run as a
  virtual machine so that a separate physical machine is not required.

Software was tested on the following:

- An [Avid VENUE | Profile][Profile] System.
- Windows 7, running on [VMware Fusion][Fusion] (which provides a built-in VNC
  server), running on OS X Yosemitte.

## Installation
This code is written in Golang (<http://golang.org/>).

1. Install Golang. Follow the instructions at <http://golang.org/doc/install>.
2. Setup environment. (Examples are for OS X or Linux.)

    ```
    $ mkdir -p "${HOME}/opt/go/bin"
    $ export GOROOT="/usr/local/go"
    $ export GOPATH="${HOME}/opt/go"
    $ export GOBIN="${GOPATH}/bin"
    $ export PATH="${PATH}:${GOROOT}/bin:${GOBIN}"
    ```

3. Download software.

    ```
    $ go get github.com/howeyc/gopass
    $ go get github.com/kward/go-osc
    $ go get github.com/kward/go-vnc
    $ go get github.com/kward/venue
    ```

4. Test the client software. The default VNC port is 5900.

    ```
    $ cd "${GOPATH}/src/github.com/kward/venue"
    $ go run client/venue_cli.go --host <hostname/IP> --port 5900 --passwd <passwd>
    Press CTRL-C to exit.
    ```

5. Install the TouchOSC layout. TODO(kward): Document this.

6. Test the server software. Configure TouchOSC to connect to the hostname/IP
   of your machine (not the host running VENUE).

   ```
   $ go run server/venue_osc.go --venue_host <hostname/IP> --venue_port 5900 --venue_passwd <passwd> --osc_server_port 8000
   ```

Notes:

- If you are using the default VNC port of 5900, the `--port` or `--venue_port`
  options may be left off.
- If you are using the default TouchOSC port of 8000, the `--osc_server_port`
  option may be left off.

## Journal

* 2015-06-23: Initial checkin. The software is not yet ready for general usage.
  Only basic channel selection and send level adjustments work at this time.

Avid™ is a registerd trademark of Avid, Inc.

[Fusion]: http://www.vmware.com/products/fusion/
[Profile]: https://www.avid.com/US/products/profile-system
[VENUE]: http://www.avid.com/us/products/venue-software
