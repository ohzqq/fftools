package tool

import (
	"bytes"
	"log"
	"os"
	"os/exec"

	"github.com/ohzqq/avtools/ffmpeg"
	"github.com/ohzqq/avtools/media"
)

type Cmd struct {
	Flag
	Media     *media.Media
	isVerbose bool
	cwd       string
	Batch     []*exec.Cmd
	batch     []Command
	tmpFile   string
	num       int
}

type Command interface {
	Build() (*exec.Cmd, error)
	String() string
	ParseArgs() ([]string, error)
	Run() ([]byte, error)
}

func NewCmd() *Cmd {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	return &Cmd{
		cwd: cwd,
	}
}

func (c *Cmd) Input(i string) *Cmd {
	c.Media = media.NewMedia(i)
	return c
}

func (c *Cmd) Exec(bin string, args []string) *Cmd {
	cmd := exec.Command(bin, args...)
	c.AddCmd(cmd)
	return c
}

func (c *Cmd) AddCmd(cmd *exec.Cmd) *Cmd {
	c.Batch = append(c.Batch, cmd)
	return c
}

func (c *Cmd) Verbose() *Cmd {
	c.isVerbose = true
	return c
}

func (c *Cmd) SetFlags(f Flag) *Cmd {
	c.Flag = f
	c.Media = f.Media()
	return c
}

func (c *Cmd) NewFFmpegCmd() *ffmpeg.Cmd {
	cmd := Cfg().GetProfile("default").FFmpegCmd()

	if c.Flag.Args.HasProfile() {
		cmd = Cfg().GetProfile(c.Flag.Args.Profile).FFmpegCmd()
	}

	if c.Bool.Verbose {
		cmd.LogLevel("info")
	}

	if c.Bool.Overwrite {
		cmd.AppendPreInput("y")
	}

	if c.Args.HasStart() {
		cmd.AppendPreInput("ss", c.Args.Start)
	}

	if c.Args.HasEnd() {
		cmd.AppendPreInput("to", c.Args.End)
	}

	if c.Media != nil {
		cmd.Input(c.Media.Input)
	}

	if c.Args.HasMeta() {
		cmd.FFmeta(c.Args.Meta)
	}

	var out *Output
	if c.Args.HasOutput() {
		out = NewOutput(c.Args.Output)
	} else {
		out = NewOutput(c.Args.Input)
	}
	cmd.Output(out.String())

	return cmd
}

func (cmd *Cmd) Tmp(f string) *Cmd {
	cmd.tmpFile = f
	return cmd
}

func (c Cmd) RunBatch() []byte {
	var buf bytes.Buffer
	for _, cmd := range c.batch {
		out, err := cmd.Run()
		if err != nil {
			log.Fatal(err)
		}

		_, err = buf.Write(out)
		if err != nil {
			log.Fatal(err)
		}
	}
	return buf.Bytes()
}

//func (cmd Cmd) Run() []byte {
//  if cmd.tmpFile != "" {
//    defer os.Remove(cmd.tmpFile)
//  }

//  var (
//    stderr bytes.Buffer
//    stdout bytes.Buffer
//  )
//  cmd.exec.Stderr = &stderr
//  cmd.exec.Stdout = &stdout

//  err := cmd.exec.Run()
//  if err != nil {
//    log.Fatal("Command finished with error: %v\n", cmd.exec.String())
//    fmt.Printf("%v\n", stderr.String())
//  }

//  if len(stdout.Bytes()) > 0 {
//    return stdout.Bytes()
//  }

//  if cmd.isVerbose {
//    fmt.Println(cmd.exec.String())
//  }
//  return nil
//}
