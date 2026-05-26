package output

import (
	"encoding/json"
	"fmt"
	"os"
)

const (
	ExitSuccess  = 0
	ExitError    = 1
	ExitAuth     = 2
	ExitNotFound = 3
	ExitNetwork  = 4
	ExitArgs     = 5
)

type Envelope struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
	Code    int         `json:"code"`
}

type Printer struct {
	JSON bool
}

func NewPrinter(json bool) *Printer {
	return &Printer{JSON: json}
}

func (p *Printer) PrintSuccess(data interface{}) {
	if p.JSON {
		p.printJSON(Envelope{Success: true, Data: data, Code: 0})
	}
}

func (p *Printer) PrintError(msg string, code int) int {
	if p.JSON {
		p.printJSON(Envelope{Success: false, Error: msg, Code: code})
		return code
	}
	fmt.Fprintf(os.Stderr, "Error: %s\n", msg)
	return code
}

func (p *Printer) printJSON(env Envelope) {
	bytes, _ := json.MarshalIndent(env, "", "  ")
	fmt.Println(string(bytes))
}
