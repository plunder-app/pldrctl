package ux

import (
	"time"

	"github.com/gookit/color"
	"github.com/plunder-app/plunder/pkg/plunderlogging"
)

//LogsGetFormat -
func LogsGetFormat(logs plunderlogging.JSONLog) {

	color.White.Printf("Logs:\n")
	for i := range logs.Entries {
		color.Green.Printf("Started: ")
		color.White.Printf("%s\t", logs.Entries[i].Created.Format(time.ANSIC))
		color.Green.Printf("Task Name: ")
		color.White.Printf("%s\n", logs.Entries[i].TaskName)
		if logs.Entries[i].Entry != "" {
			color.Green.Printf("Output:\n")
			color.White.Printf("%s\n", logs.Entries[i].Entry)
		}
		if logs.Entries[i].Err != "" {
			color.Red.Printf("Error:\n")
			color.White.Printf("%s\n", logs.Entries[i].Err)
		}
	}
	color.White.Printf("Task Status: ")
	switch logs.State {
	case "Completed":
		color.Green.Printf("Completed\n")
	case "Running":
		color.Blue.Printf("Running\n")
	case "Failed":
		color.Red.Printf("Failed\n")
	}

}
