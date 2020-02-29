package main

import (
	"bytes"
	"encoding/xml"
	"log"
	"os"
	"path/filepath"

	"golang.org/x/arch/x86/x86csv"

	"github.com/Maratyszcza/Opcodes/go/opcodesxml"
)

func main() {
	f, err := os.Open(filepath.Join("./", "x86.v0.2.csv"))
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	x := x86csv.NewReader(f)
	insts, err := x.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	opcodes := make(map[string]string)
	for _, inst := range insts {
		opcodes[inst.GNUOpcode()] = inst.GoOpcode()
	}

	instSet, err := opcodesxml.ReadFile(filepath.Join("./", "x86_64.golden.xml"))
	if err != nil {
		log.Fatal(err)
	}
	for i, inst := range instSet.Instructions {
		for j, form := range inst.Forms {
			if instSet.Instructions[i].Forms[j].GoName == "" {
				if goOpcode, ok := opcodes[form.GASName]; ok {
					instSet.Instructions[i].Forms[j].GoName = goOpcode
				}
			}
		}
	}

	out, err := xml.MarshalIndent(instSet, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	header := []byte("<?xml version='1.0' encoding='utf-8'?>\n")
	out = append(header, out...)

	out = bytes.ReplaceAll(out, []byte("></Operand>"), []byte("/>"))
	out = bytes.ReplaceAll(out, []byte("></REX>"), []byte("/>"))
	out = bytes.ReplaceAll(out, []byte("></Immediate>"), []byte("/>"))
	out = bytes.ReplaceAll(out, []byte("></Prefix>"), []byte("/>"))
	out = bytes.ReplaceAll(out, []byte("></Opcode>"), []byte("/>"))
	out = bytes.ReplaceAll(out, []byte("></ModRM>"), []byte("/>"))
	out = bytes.ReplaceAll(out, []byte("></ISA>"), []byte("/>"))
	out = bytes.ReplaceAll(out, []byte("></VEX>"), []byte("/>"))
	out = bytes.ReplaceAll(out, []byte("></ImplicitOperand>"), []byte("/>"))
	out = bytes.ReplaceAll(out, []byte("></CodeOffset>"), []byte("/>"))
	out = bytes.ReplaceAll(out, []byte("></EVEX>"), []byte("/>"))
	out = bytes.ReplaceAll(out, []byte("></RegisterByte>"), []byte("/>"))

	_, err = os.Stdout.Write(append(out, '\n'))
	if err != nil {
		log.Fatal(err)
	}
}
