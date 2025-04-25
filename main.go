package main

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/godbus/dbus/v5"
)

const (
	dbusDest      = "org.gnome.SessionManager"
	dbusPath      = "/org/gnome/SessionManager"
	inhibitReason = "Preventing screen blank"
	inhibitFlags  = 0x8 // inhibit flag for idle
)

func main() {

	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}
}

func run() error {

	if len(os.Args) != 2 {
		return fmt.Errorf("usage: %s <minutes>", os.Args[0])
	}
	minutes, err := strconv.Atoi(os.Args[1])
	if err != nil || minutes < 1 {
		return errors.New("please enter a positive number of minutes")
	}
	conn, err := dbus.SessionBus()
	if err != nil {
		return fmt.Errorf("connect to session bus: %w", err)
	}
	defer conn.Close()
	obj := conn.Object(dbusDest, dbus.ObjectPath(dbusPath))
	var cookie uint32
	err = obj.Call("org.gnome.SessionManager.Inhibit", 0,
		os.Args[0], uint32(0), inhibitReason, uint32(inhibitFlags)).Store(&cookie)
	if err != nil {
		return fmt.Errorf("inhibit screen: %w", err)
	}
	defer func() {
		if err := obj.Call("org.gnome.SessionManager.Uninhibit", 0, cookie).Err; err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to uninhibit: %s\n", err)
		}
	}()
	fmt.Printf("Inhibiting screen blank for %d minutes...\n", minutes)
	fmt.Println("Press Ctrl+C to stop")
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt)
	timer := time.NewTimer(time.Duration(minutes) * time.Minute)
	defer timer.Stop()
	select {
	case <-timer.C:
	case <-done:
		fmt.Println("\nStopping early...")
	}
	return nil
}
