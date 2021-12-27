package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/mitchellh/go-ps"
	"github.com/olekukonko/ts"
)

type processType struct {
	Name string
	Pid  int
}

func contains(processList []processType, pid int) bool {
	for _, v := range processList {
		if v.Pid == pid {
			return true
		}
	}

	return false
}

func getProccess() []processType {
	var procList []processType

	processList, err := ps.Processes()

	if err != nil {
		log.Println("Task aborted unexpectedly.")
		os.Exit(2)
	}

	for x := range processList {
		process := processList[x]
		item := processType{Name: process.Executable(), Pid: process.PPid()}

		if !contains(procList, item.Pid) {
			procList = append(procList, item)
		}
	}

	return procList
}

func confirmBack() {
	options := []string{"Back", "Quit"}

	templates := &promptui.SelectTemplates{
		Label:    "{{ . | blue }}",
		Active:   "  {{ . | cyan }}",
		Inactive: "  {{ . }}",
	}

	searcher := func(input string, index int) bool {
		option := options[index]
		name := strings.Replace(strings.ToLower(option), " ", "", -1)
		input = strings.Replace(strings.ToLower(input), " ", "", -1)

		return strings.Contains(name, input)
	}

	prompt := promptui.Select{
		Label:     "Process killed successfully!",
		Items:     options,
		Templates: templates,
		Searcher:  searcher,
		Size:      2,
	}

	_, result, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		os.Exit(1)
	}

	if result == "Back" {
		main()
	} else {
		os.Exit(0)
	}
}

func confirmKill(pid int) {

	options := []string{"Yes", "No"}

	templates := &promptui.SelectTemplates{
		Label:    "{{ . | red }}",
		Active:   "  {{ . | cyan }}",
		Inactive: "  {{ . }}",
	}

	searcher := func(input string, index int) bool {
		option := options[index]
		name := strings.Replace(strings.ToLower(option), " ", "", -1)
		input = strings.Replace(strings.ToLower(input), " ", "", -1)

		return strings.Contains(name, input)
	}

	prompt := promptui.Select{
		Label:     "Are you sure you want to kill this process?",
		Items:     options,
		Templates: templates,
		Searcher:  searcher,
		Size:      2,
	}

	_, result, promptErr := prompt.Run()

	if promptErr != nil {
		fmt.Printf("Prompt failed %v\n", promptErr)
		os.Exit(0)
	}

	if result == "Yes" {
		proc, procErr := os.FindProcess(pid)

		if procErr != nil {
			fmt.Printf("Prompt failed %v\n", procErr)
			os.Exit(1)
		}

		procErr = proc.Kill()

		if procErr != nil {
			fmt.Printf("Prompt failed %v\n", procErr)
			os.Exit(1)
		}

		confirmBack()
	} else {
		main()
	}
}

func main() {

	procList := getProccess()
	size, _ := ts.GetSize()

	templates := &promptui.SelectTemplates{
		Label:    "{{ . |  green }}",
		Active:   "\u25b8 {{ .Name | cyan }}",
		Inactive: "  {{ .Name }}",
		Selected: "  {{ .Name | green }}",
		Details: `
{{ "--------- Details ----------" | bold }}
{{ "Process Name:" | magenta | faint }} {{ .Name }}
{{ "Process ID:" | magenta | faint }} {{ .Pid }}`,
	}

	searcher := func(input string, index int) bool {
		process := procList[index]
		name := strings.Replace(strings.ToLower(process.Name), " ", "", -1)
		input = strings.Replace(strings.ToLower(input), " ", "", -1)

		return strings.Contains(name, input)
	}

	prompt := promptui.Select{
		Label:     "Process to kill:",
		Items:     procList,
		Templates: templates,
		Size:      (size.Row() * 5) / 10,
		Searcher:  searcher,
	}

	i, _, promptErr := prompt.Run()

	if promptErr != nil {
		fmt.Printf("Prompt failed %v\n", promptErr)
		os.Exit(0)
	}

	confirmKill(procList[i].Pid)
}
