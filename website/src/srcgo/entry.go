package main

import (
	"encoding/json"
	"fmt"
	"syscall/js"
	"time"

	"github.com/osdc/resrap"
	"github.com/osdc/resrap/templates"
)

type browserTemplate struct {
	Name        string `json:"name"`
	Language    string `json:"language"`
	Description string `json:"description"`
	Grammar     string `json:"grammar"`
}

func main() {
	println("Resrap WASM starting")
	par := resrap.NewResrap()

	js.Global().Set("generateText", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if len(args) < 4 {
			return "Error: need grammar, startpoint, seed, tokenlen"
		}
		grammar := args[0].String()
		startpoint := args[1].String()
		seed := browserSeed(args[2].Int())
		tokenlen := args[3].Int()
		if tokenlen <= 0 {
			return "Error: tokenlen must be greater than zero"
		}
		if err := par.ParseGrammar("sample", grammar); err != nil {
			return "Error: " + err.Error()
		}
		return par.GenerateCompleteWithSeeded("sample", startpoint, seed, tokenlen)
	}))

	js.Global().Set("generateTemplate", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if len(args) < 4 {
			return "Error: need template, startpoint, seed, tokenlen"
		}
		tmpl, err := templates.Get(args[0].String())
		if err != nil {
			return "Error: " + err.Error()
		}
		seed := browserSeed(args[2].Int())
		tokenlen := args[3].Int()
		if tokenlen <= 0 {
			return "Error: tokenlen must be greater than zero"
		}
		if err := tmpl.Load(par); err != nil {
			return "Error: " + err.Error()
		}
		return par.GenerateCompleteWithSeeded(tmpl.Name, args[1].String(), seed, tokenlen)
	}))

	js.Global().Set("getTemplates", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		list, err := templates.List()
		if err != nil {
			return "[]"
		}
		out := make([]browserTemplate, 0, len(list))
		for _, tmpl := range list {
			out = append(out, browserTemplate{
				Name:        tmpl.Name,
				Language:    tmpl.Language,
				Description: tmpl.Description,
				Grammar:     tmpl.Grammar,
			})
		}
		body, err := json.Marshal(out)
		if err != nil {
			fmt.Println("marshal templates:", err)
			return "[]"
		}
		return string(body)
	}))

	js.Global().Set("goTest", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		return "Hello from Resrap WASM!"
	}))

	js.Global().Set("goAdd", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if len(args) < 2 {
			return 0
		}
		return args[0].Int() + args[1].Int()
	}))

	println("Resrap WASM ready")
	<-make(chan bool)
}

// browserSeed treats zero as "pick a random seed", matching the UI hint.
func browserSeed(n int) uint64 {
	if n == 0 {
		return uint64(time.Now().UnixNano())
	}
	return uint64(n)
}
