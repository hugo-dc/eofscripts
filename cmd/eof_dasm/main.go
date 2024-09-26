package main

import (
	"encoding/hex"
	"fmt"
	"os"
	"strconv"

	"github.com/hugo-dc/eofscripts/common"
)

func showUsage() {
	fmt.Println("eof_dasm - Disassemble EOF bytecode showing each section")
	fmt.Println("Usage:")
	fmt.Println("\teof_dasm <eof_code> [option]")
	fmt.Println("Options:")
	fmt.Println("\t-m | -magic\t\tRemove magic")
	fmt.Println("\t-m | +magic\t\tShow magic")
	fmt.Println("\t-v | -version\t\tRemove version")
	fmt.Println("\t-v | +version\t\tShow version")
	fmt.Println("\t-t | -types-header\t\tRemove types section header")
	fmt.Println("\t-t | +types-header\t\tShow types section header")
	fmt.Println("\t-c | -code-header\t\tRemove code section header")
	fmt.Println("\t-c | +code-header\t\tShow code section header")
	fmt.Println("\t-d | -data-header\t\tRemove data section header")
	fmt.Println("\t-d | +data-header\t\tShow data section header")
	fmt.Println("\t-tm | -terminator\t\tRemove terminator ")
	fmt.Println("\t-ts [n] | -type-section [n]\tRemove type section n (default:0)")
	fmt.Println("\t-ts [n] | +type-section [n]\tShow type section n (default:0)")
	fmt.Println("\t~ts <n> <type> | ~type-section <n> <type>\tReplace type section n (default:0)")
	fmt.Println("\t-cs [n] | -code-section [n]\tRemove code section n (default:0)")
	fmt.Println("\t-cs [n] | +code-section [n]\tShow code section n (default:0)")
	fmt.Println("\t-ds | -data-section\t\tRemove data section")
	fmt.Println("\t-ds | +data-section\t\tShow data section")
}

type ContractSections struct {
	Magic           bool
	Version         bool
	TypeHeader      bool
	CodeHeader      bool
	ContainerHeader bool
	DataHeader      bool
	Terminator      bool
	TypeSection     map[int]bool
	CodeSection     map[int]bool
	DataSection     bool
}

func main() {
	if len(os.Args) == 1 {
		showUsage()
		return
	}

	cs := ContractSections{
		Magic:           true,
		Version:         true,
		TypeHeader:      true,
		CodeHeader:      true,
		ContainerHeader: true,
		DataHeader:      true,
		Terminator:      true,
		TypeSection:     make(map[int]bool),
		CodeSection:     make(map[int]bool),
		DataSection:     true,
	}
	modTypeSections := make(map[int]string)
	eofStr := os.Args[1]
	eofBytecode, err := hex.DecodeString(eofStr)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	eofObject, err := common.ParseEOF(eofBytecode)

	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	for i, a := range os.Args {
		switch {
		case a == "-m" || a == "-magic":
			cs.Magic = false
			break
		case a == "-v" || a == "-version":
			cs.Version = false
			break
		case a == "-t" || a == "-types-header":
			cs.TypeHeader = false
			break
		case a == "-c" || a == "-code-header":
			cs.CodeHeader = false
			break
		case a == "-d" || a == "-data-header":
			cs.DataHeader = false
			break
		case a == "-tm" || a == "-terminator":
			cs.Terminator = false
			break
		case a == "-ts" || a == "-types-section":
			section := -1
			if len(os.Args) > i+1 {
				section_, err := strconv.ParseInt(os.Args[i+1], 10, 64)
				if err == nil {
					section = int(section_)
				}
			}
			cs.TypeSection[section] = false
			break
		case a == "~ts" || a == "~types-section":
			section := 0
			if len(os.Args) > i+1 {
				section_, err := strconv.ParseInt(os.Args[i+1], 10, 64)

				if err == nil {
					section = int(section_)
				}
			}
			cs.TypeSection[section] = true
			if len(os.Args[i+2]) != 8 {
				fmt.Println("Error: type must be in format 00000000")
				return
			}
			modTypeSections[section] = os.Args[i+2]
			break
		case a == "-cs" || a == "-code-section":
			section := -1
			if len(os.Args) > i+1 {
				section_, err := strconv.ParseInt(os.Args[i+1], 10, 64)

				if err == nil {
					section = int(section_)
				}
				cs.CodeSection[section] = false
			}
			break
		case a == "-ds" || a == "-data-section":
			cs.DataSection = false
			break
		}
	}

	if cs.Magic {
		fmt.Print("ef00")
	}
	if cs.Version {
		fmt.Printf("%02d", eofObject.Version)
	}
	if cs.TypeHeader {
		fmt.Printf("%02x%04x", common.CTypeId, len(eofObject.Types)*4)
	}
	if cs.CodeHeader {
		ch := fmt.Sprintf("%02x%04x", common.CCodeId, len(eofObject.CodeSections))

		for _, cs := range eofObject.CodeSections {
			ch += fmt.Sprintf("%04x", len(cs)/2)
		}
		fmt.Print(ch)
	}
	if cs.ContainerHeader && len(eofObject.ContainerSections) > 0 {
		fmt.Printf("%02x%04x", common.CContainerId, len(eofObject.ContainerSections))

		for _, cs := range eofObject.ContainerSections {
			fmt.Printf("%04x", len(cs)/2)
		}
	}

	if cs.DataHeader {
		fmt.Printf("%02x%04x", common.CDataId, len(eofObject.Data)/2)
	}
	if cs.Terminator {
		fmt.Print("00")
	}

	hideAll := false
	if _, ok := cs.TypeSection[-1]; ok {
		hideAll = true
	}

	if hideAll == false {
		for i, t := range eofObject.Types {
			showSection := true
			if v, ok := cs.TypeSection[i]; ok {
				showSection = v
			} else {
				if len(cs.TypeSection) == 0 {
					showSection = true
				}
			}
			if showSection {
				if s, ok := modTypeSections[i]; ok {
					fmt.Printf("%s", s)
				} else {
					fmt.Printf("%02x%02x%04x", t[0], t[1], t[2])
				}
			}
		}
	}

	hideAll = false
	if _, ok := cs.CodeSection[-1]; ok {
		hideAll = true
	}

	if hideAll != true {
		for i, c := range eofObject.CodeSections {
			showSection := false
			if v, ok := cs.CodeSection[i]; ok {
				showSection = v
			} else {
				if len(cs.CodeSection) == 0 {
					showSection = true
				}
			}

			if showSection {
				fmt.Printf("%s", c)
			}
		}
	}

	if cs.DataSection {
		fmt.Println(eofObject.Data)
	}
}
