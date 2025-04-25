# gojolt

gojolt keeps your screen awake preventing screen blank and automatic suspension.

It uses [godbus](https://github.com/godbus/dbus) D-Bus bindings.

It works on Gnome Shell on Xorg and Wayland.

## Building from source

Simply clone this repo and run: `go build`

## Usage

It's a CLI program.

The duration of screen blank inhibition in minutes has to be passed as the first and only argument.
