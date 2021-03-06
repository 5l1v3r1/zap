package main

import (
	"fmt"
	"io"
	"os"

	"github.com/containerd/console"
	"github.com/urfave/cli"
	"github.com/wybiral/zap/pkg/repl"
)

const version = "0.1.0"

func main() {
	// hide default flags
	cli.HelpFlag = &cli.StringFlag{Hidden: true}
	cli.VersionFlag = &cli.StringFlag{Hidden: true}
	// setup CLI app
	c := cli.NewApp()
	c.CommandNotFound = func(ctx *cli.Context, command string) {
		fmt.Printf("Command not found: %v\n", command)
		os.Exit(1)
	}
	c.Version = version
	c.Usage = "MicroPython CLI tool"
	c.Commands = []*cli.Command{
		&cli.Command{
			Name:      "cat",
			Usage:     "Read file",
			Action:    cmdCat,
			ArgsUsage: "file",
		},
		&cli.Command{
			Name:      "cd",
			Usage:     "Change directory",
			Action:    cmdCd,
			ArgsUsage: "path",
		},
		&cli.Command{
			Name:   "download",
			Usage:  "Copy all files from device to local directory",
			Action: cmdDownload,
		},
		&cli.Command{
			Name:      "get",
			Usage:     "Copy a file from the device",
			Action:    cmdGet,
			ArgsUsage: "dst src",
		},
		&cli.Command{
			Name:      "help",
			Usage:     "Shows all commands or help for one command",
			ArgsUsage: "[command]",
			Action:    cmdHelp,
		},
		&cli.Command{
			Name:   "ls",
			Usage:  "List files",
			Action: cmdLs,
		},
		&cli.Command{
			Name:      "mkdir",
			Usage:     "Make directory",
			Action:    cmdMkdir,
			ArgsUsage: "dir",
		},
		&cli.Command{
			Name:      "put",
			Usage:     "Copy a file to the device",
			Action:    cmdPut,
			ArgsUsage: "dst src",
		},
		&cli.Command{
			Name:   "pwd",
			Usage:  "Print working directory",
			Action: cmdPwd,
		},
		&cli.Command{
			Name:   "reboot",
			Usage:  "Perform a soft reboot",
			Action: cmdReboot,
		},
		&cli.Command{
			Name:   "repl",
			Usage:  "Open the MicroPython REPL",
			Action: cmdRepl,
		},
		&cli.Command{
			Name:      "rm",
			Usage:     "Delete file",
			Action:    cmdRm,
			ArgsUsage: "file",
		},
		&cli.Command{
			Name:      "rmdir",
			Usage:     "Remove directory",
			Action:    cmdRmdir,
			ArgsUsage: "dir",
		},
		&cli.Command{
			Name:   "upload",
			Usage:  "Copy all files from local directory to device",
			Action: cmdUpload,
		},
		&cli.Command{
			Name:  "version",
			Usage: "Print zap version",
			Action: func(ctx *cli.Context) error {
				fmt.Println(c.Version)
				return nil
			},
		},
	}
	c.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:     "device",
			Aliases:  []string{"d"},
			Usage:    "Serial device name of MicroPython board",
			Required: true,
			EnvVars:  []string{"PYBOARD_DEVICE"},
		},
		&cli.IntFlag{
			Name:    "baudrate",
			Aliases: []string{"b"},
			Value:   115200,
			Usage:   "Baudrate of serial device",
			EnvVars: []string{"PYBOARD_BAUDRATE"},
		},
	}
	// run CLI app
	err := c.Run(os.Args)
	if err != nil {
		fmt.Println("\nERROR:", err)
	}
}

func cmdCat(ctx *cli.Context) error {
	r, err := repl.Connect(ctx.String("device"), ctx.Int("baudrate"))
	if err != nil {
		return err
	}
	err = r.EnterRawMode()
	if err != nil {
		return err
	}
	defer r.ExitRawMode()
	return r.Cat(os.Stdout, ctx.Args().Get(0))
}

func cmdCd(ctx *cli.Context) error {
	r, err := repl.Connect(ctx.String("device"), ctx.Int("baudrate"))
	if err != nil {
		return err
	}
	err = r.EnterRawMode()
	if err != nil {
		return err
	}
	defer r.ExitRawMode()
	return r.Cd(ctx.Args().Get(0))
}

func cmdDownload(ctx *cli.Context) error {
	r, err := repl.Connect(ctx.String("device"), ctx.Int("baudrate"))
	if err != nil {
		return err
	}
	err = r.EnterRawMode()
	if err != nil {
		return err
	}
	defer r.ExitRawMode()
	return r.Download()
}

func cmdGet(ctx *cli.Context) error {
	r, err := repl.Connect(ctx.String("device"), ctx.Int("baudrate"))
	if err != nil {
		return err
	}
	err = r.EnterRawMode()
	if err != nil {
		return err
	}
	defer r.ExitRawMode()
	args := ctx.Args()
	dst := args.Get(0)
	src := dst
	if args.Len() > 1 {
		src = args.Get(1)
	}
	return r.Get(dst, src)
}

func cmdHelp(ctx *cli.Context) error {
	args := ctx.Args()
	if args.Present() {
		cli.ShowCommandHelp(ctx, args.First())
		return nil
	}
	cli.ShowAppHelp(ctx)
	return nil
}

func cmdLs(ctx *cli.Context) error {
	r, err := repl.Connect(ctx.String("device"), ctx.Int("baudrate"))
	if err != nil {
		return err
	}
	err = r.EnterRawMode()
	if err != nil {
		return err
	}
	defer r.ExitRawMode()
	fs, err := r.Ls()
	if err != nil {
		return err
	}
	for _, f := range fs {
		fmt.Print(f + "  ")
	}
	fmt.Print("\n")
	return nil
}

func cmdMkdir(ctx *cli.Context) error {
	r, err := repl.Connect(ctx.String("device"), ctx.Int("baudrate"))
	if err != nil {
		return err
	}
	err = r.EnterRawMode()
	if err != nil {
		return err
	}
	defer r.ExitRawMode()
	return r.Mkdir(ctx.Args().Get(0))
}

func cmdPut(ctx *cli.Context) error {
	r, err := repl.Connect(ctx.String("device"), ctx.Int("baudrate"))
	if err != nil {
		return err
	}
	err = r.EnterRawMode()
	if err != nil {
		return err
	}
	defer r.ExitRawMode()
	args := ctx.Args()
	dst := args.Get(0)
	src := dst
	if args.Len() > 1 {
		src = args.Get(1)
	}
	return r.Put(dst, src)
}

func cmdPwd(ctx *cli.Context) error {
	r, err := repl.Connect(ctx.String("device"), ctx.Int("baudrate"))
	if err != nil {
		return err
	}
	err = r.EnterRawMode()
	if err != nil {
		return err
	}
	defer r.ExitRawMode()
	cwd, err := r.Cwd()
	if err != nil {
		return err
	}
	fmt.Println(cwd)
	return nil
}

func cmdReboot(ctx *cli.Context) error {
	r, err := repl.Connect(ctx.String("device"), ctx.Int("baudrate"))
	if err != nil {
		return err
	}
	err = r.EnterRawMode()
	if err != nil {
		return err
	}
	defer r.ExitRawMode()
	return r.SoftReboot()
}

func cmdRepl(ctx *cli.Context) error {
	r, err := repl.Connect(ctx.String("device"), ctx.Int("baudrate"))
	if err != nil {
		return err
	}
	current := console.Current()
	defer current.Reset()
	err = current.SetRaw()
	if err != nil {
		return err
	}
	go io.Copy(os.Stdout, r.Port)
	io.Copy(r.Port, os.Stdin)
	return nil
}

func cmdRm(ctx *cli.Context) error {
	r, err := repl.Connect(ctx.String("device"), ctx.Int("baudrate"))
	if err != nil {
		return err
	}
	err = r.EnterRawMode()
	if err != nil {
		return err
	}
	defer r.ExitRawMode()
	return r.Rm(ctx.Args().Get(0))
}

func cmdRmdir(ctx *cli.Context) error {
	r, err := repl.Connect(ctx.String("device"), ctx.Int("baudrate"))
	if err != nil {
		return err
	}
	err = r.EnterRawMode()
	if err != nil {
		return err
	}
	defer r.ExitRawMode()
	return r.Rmdir(ctx.Args().Get(0))
}

func cmdUpload(ctx *cli.Context) error {
	r, err := repl.Connect(ctx.String("device"), ctx.Int("baudrate"))
	if err != nil {
		return err
	}
	err = r.EnterRawMode()
	if err != nil {
		return err
	}
	defer r.ExitRawMode()
	return r.Upload()
}
