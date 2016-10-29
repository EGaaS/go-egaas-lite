// Copyright 2016 The go-daylight Authors
// This file is part of the go-daylight library.
//
// The go-daylight library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-daylight library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-daylight library. If not, see <http://www.gnu.org/licenses/>.

package main

import (
	//	"fmt"
	"flag"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"

	"github.com/EGaaS/go-egaas-lite/system"
	"github.com/go-thrust/lib/bindings/window"
	"github.com/go-thrust/lib/commands"
	"github.com/go-thrust/thrust"
)

const GETPOOLURL = `http://node0.egaas.org/`

var (
	DevTools int64
)

func init() {
	flag.Int64Var(&DevTools, "dev", 0, "Devtools in thrust-shell")
	flag.Parse()
}

func main() {
	var (
		thrustWindow *window.Window
		//		mainWin      bool
	)

	runtime.LockOSThread()
	//	if utils.Desktop() && (winVer() >= 6 || winVer() == 0) {
	thrust.Start()

	thrust.NewEventHandler("*", func(cr commands.CommandResponse) {
		if cr.Type == "closed" {
			//			if mainWin || !isClosed {
			system.FinishThrust(0)
			os.Exit(0)
			/*			} else {
						close(chIntro)
						mainWin = true
					}*/
		}
	})

	thrustWindow = thrust.NewWindow(thrust.WindowOptions{
		Title:   "EGaaS Lite",
		RootUrl: GETPOOLURL,
		Size:    commands.SizeHW{Width: 1024, Height: 800},
	})
	/*	thrustWindow.HandleEvent("*", func(cr commands.EventResult) {
		fmt.Println("HandleEvent", cr)
	})*/
	if DevTools != 0 {
		thrustWindow.OpenDevtools()
	}
	accfile := filepath.Join(GetCurrentDir(), `accounts.txt`)
	thrustWindow.HandleRemote(func(er commands.EventResult, this *window.Window) {
		//		fmt.Println("RemoteMessage Recieved:", er.Message.Payload)
		if len(er.Message.Payload) > 7 && er.Message.Payload[:2] == `[{` {
			ioutil.WriteFile(accfile, []byte(er.Message.Payload), 0644)
		} else if er.Message.Payload == `ACCOUNTS` {
			accounts, _ := ioutil.ReadFile(accfile)
			this.SendRemoteMessage(string(accounts))
		} else {
			ShellExecute(er.Message.Payload)
		}
	})
	thrustWindow.Show()
	thrustWindow.Focus()
	for {
		time.Sleep(3600 * time.Second)
	}
	system.FinishThrust(0)
}

func ShellExecute(cmdline string) {
	time.Sleep(500 * time.Millisecond)
	switch runtime.GOOS {
	case "linux":
		exec.Command("xdg-open", cmdline).Start()
	case "windows":
		exec.Command(`rundll32.exe`, `url.dll,FileProtocolHandler`, cmdline).Start()
	case "darwin":
		exec.Command("open", cmdline).Start()
	}
}

func GetCurrentDir() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return "."
	}
	return dir
}
