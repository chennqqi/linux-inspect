package main

import (
	"fmt"
	"os"

	"github.com/gyuho/linux-inspect/inspect"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

type ssFlags struct {
	protocol  string
	program   string
	localPort int64
	top       int
}

var (
	ssCommand = &cobra.Command{
		Use:   "ss",
		Short: "Inspects '/proc/net/tcp,tcp6'",
		RunE:  ssCommandFunc,
	}
	ssCmdFlag = ssFlags{}
)

func init() {
	ssCommand.PersistentFlags().StringVarP(&ssCmdFlag.protocol, "protocol", "c", "tcp", "Specify the protocol ('tcp' or 'tcp6').")
	ssCommand.PersistentFlags().StringVarP(&ssCmdFlag.program, "program", "s", "", "Specify the program name.")
	ssCommand.PersistentFlags().Int64VarP(&ssCmdFlag.localPort, "local-port", "l", -1, "Specify the local port.")
	ssCommand.PersistentFlags().IntVarP(&ssCmdFlag.top, "top", "t", 5, "Limit the number results to return.")
}

func ssCommandFunc(cmd *cobra.Command, args []string) error {
	color.Set(color.FgMagenta)
	fmt.Fprintf(os.Stdout, "\n'ss' to inspect '/proc/net/tcp,tcp6'\n\n")
	color.Unset()

	opts := []inspect.FilterFunc{inspect.WithTCP()}
	if ssCmdFlag.protocol == "tcp6" {
		opts[0] = inspect.WithTCP6()
	} else if ssCmdFlag.protocol != "tcp" {
		fmt.Fprintf(os.Stderr, "unknown protocol %q\n", ssCmdFlag.protocol)
		os.Exit(233)
	}
	opts = append(opts, inspect.WithProgram(ssCmdFlag.program), inspect.WithLocalPort(ssCmdFlag.localPort), inspect.WithTopLimit(ssCmdFlag.top))
	sss, err := inspect.GetSS(opts...)
	if err != nil {
		return err
	}
	hd, rows := inspect.ConvertSS(sss...)
	txt := inspect.StringSS(hd, rows, -1)
	fmt.Print(txt)

	color.Set(color.FgGreen)
	fmt.Fprintf(os.Stdout, "\nDONE!\n")
	color.Unset()

	return nil
}
