package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	cmd := os.Args
	if len(cmd) < 2 {
		fmt.Fprintf(os.Stderr, "Erorr no command present")
	}
	command := cmd[1]
	sc, err := NewSwayConnection()
	if err != nil {
		panic(err)
	}

	switch command {
	case "list":
		WorkspacesList(0, sc)
	case "focus":
		WorkspacesList(1, sc)
	case "urgent":
		WorkspacesList(2, sc)
	case "fww":
		WorkspacesList(3, sc)
	case "listen":
		Listener(sc)
	case "-c":
		if len(cmd) >= 2 {
			cmda := cmd[2:]
			all_command := strings.Join(cmda, " ")
			_, err := sc.RunSwayCommand(fmt.Sprint(all_command))
			if err != nil {
				fmt.Print(err)
				os.Exit(1)
			}
		} else {
			fmt.Fprintf(os.Stderr, "Command not present")
		}
	default:
		fmt.Println("Command not found")
	}
}

func Listener(sc *SwayConnection) {
	pp, err := sc.SendCommand(IPC_SUBSCRIBE, `["window"]`)
	if err != nil {
		panic(err)
	}
	printStr(pp)
	s := sc.Subscribe()
	defer s.Close()

	for {
		select {
		case event := <-s.Events:
			fmt.Println(event.Change)
			if event.Change == "new" {
				sc.RunSwayCommand(fmt.Sprintf("[con_id=%d] split h", event.Container.ID))
				sc.RunSwayCommand(fmt.Sprintf("[con_id=%d] move down", event.Container.ID))
			}
		case err := <-s.Errors:
			fmt.Println("Error:", err)
			break
		}
	}
}

func WorkspacesList(Types int, sc *SwayConnection) {
	switch Types {
	case 1:
		ws, err := sc.GetFocusedWorkspace()
		if err != nil {
			panic(err)
		}
		fmt.Print(ws.Name)
	case 2:
		ws, err := sc.GetUrgentWorkspace()
		if err != nil {
			panic(err)
		}
		if ws != nil {
			fmt.Print(ws.Name)
		}
	case 3:
		windows, err := sc.GetFocusedWorkspaceWindows()
		if err != nil {
			panic(err)
		}
		for _, window := range windows {
			fmt.Println(window.Name)
		}
	default:
		ws, err := sc.GetWorkspaces()
		if err != nil {
			panic(err)
		}
		for _, workspace := range ws {
			fmt.Println(workspace.Name)
		}
	}
}

func printStr(bytes []byte) {
	fmt.Print(string(bytes))
}
