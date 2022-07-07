package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"

	"github.com/oxodao/metaprint/config"
	"github.com/oxodao/metaprint/pulse"
)

var (
    AUTHOR        = "Oxodao"
    CONTRIBUTORS  = []string{"digitallyserviced"}
    VERSION       = "DEV"
    COMMIT        = "XXXXXXXXX"
    SOFTWARE_NAME = "metaprint"
)

var otherCommands = []string{"pulseaudio-infos", "version", "debug"}

func main() {
    cfg := config.Load()

    hasOtherCommands := hasOtherCommand()
    if len(os.Args) < 3 && !hasOtherCommands {
        printUsage()
        os.Exit(1)
    }

    if hasOtherCommands {
        switch os.Args[1] {
        case "pulseaudio-infos":
            pulse.PrintInfos()
        case "version":
            fmt.Printf("%v %v (Commit %v) by %v\nhttps://github.com/%v/%v\nAdd'l contrib'd from: %v\n", SOFTWARE_NAME, VERSION, COMMIT, AUTHOR, strings.ToLower(AUTHOR), strings.ToLower(SOFTWARE_NAME), strings.Join(CONTRIBUTORS, ","))
            os.Exit(1)
        case "debug":
            printModuleConfig()
            os.Exit(1)
        default:
            printUsage()
            os.Exit(1)
        }
        return
    }

    module, err := cfg.FindModule(os.Args[1], os.Args[2])

    if module == nil && err == nil {
        printUsage()
        return
    } else if module == nil {
        fmt.Println(err)
        os.Exit(1)
    }

    text := module.Print(os.Args[3:])

    if len(module.GetPrefix()) > 0 {
        text = module.GetPrefix() + " " + text
    }

    if len(module.GetSuffix()) > 0 {
        text += " " + module.GetSuffix()
    }

    fmt.Println(text)
}

func printModuleConfig() {
    fmts, ok := config.Load().GetModuleFormatsAvailable()
    if ok != nil {
        panic(fmt.Errorf("error loading modules and formats: %s", ok))
    }
    writer := tabwriter.NewWriter(os.Stdout, 2, 4, 1, ' ', tabwriter.Debug)
    fmt.Fprintln(writer)
    fmt.Fprintln(writer, "Available modules: ")
    // newdump.Print(fmts)
    for n, m := range fmts {
        fmt.Fprintln(writer, fmt.Sprintf(" - ( %s ) formats:  ", n))
        for fname, fmtvals := range m.(map[string]interface{}) {
            fmt.Fprintln(writer, fmt.Sprintf("\t - %s: ", fname))
            if fname == "NA" {
                fmt.Fprintln(writer, fmt.Sprintf("\t\t - %s: ", fmtvals.(string)))
                continue
            }
            for fkname, fmtkstr := range fmtvals.(map[string]string) {
                if len(fmtkstr) != 0 {
                    fmtkstr = strconv.Quote(fmtkstr)
                } else {
                    fmtkstr = "--"
                }
                fmt.Fprintln(writer, fmt.Sprintf("\t\t - %s:\t( %s )", strconv.Quote(fkname), fmtkstr))
            }
        }
    }
    writer.Flush()
}

func printUsage() {
    fmt.Println("Usage: metaprint <module> <name> [params]")
    fmt.Println("       metaprint other-command")
    fmt.Println()
    fmt.Println("Other commands: ")
    for _, cmd := range otherCommands {
        fmt.Println("\t- " + cmd)
    }
    fmt.Println()
    fmt.Println("Available modules: ")
    for _, mod := range config.GetModulesAvailable() {
        fmt.Println("\t- " + mod)
    }
}

func hasOtherCommand() bool {
    if len(os.Args) != 2 {
        return false
    }

    for _, val := range otherCommands {
        if os.Args[1] == val {
            return true
        }
    }

    return false
}
