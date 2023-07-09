package cmd

import (
	"bytes"
	"fmt"
	"html/template"
	"runtime"

	"github.com/microctar/licorice/app/config"
	"github.com/microctar/licorice/app/constant"
	"github.com/spf13/cobra"
)

const (
	versionTmpl = `licorice:
  Version:      {{.Version}}
  Go version:   {{.GoVer}}
  Git commit:   {{.GitCommit}}
  Built:        {{.BuildTime}}
  OS/Arch:      {{.OSAndArch}}
`
)

var (
	port    uint16
	client  string
	confDir string

	clashRule     string
	clashRulePath string
	inputFile     string
	outputFile    string
)

var (
	rootCmd = &cobra.Command{
		Use:          "licorice [OPTIONS] COMMAND",
		Version:      constant.Version,
		Short:        "a utility to create configuration for tunnel",
		Long:         "licorice - a utility to create configuration for rule-based tunnel in go",
		SilenceUsage: true,
	}

	verCmd = &cobra.Command{
		Use:   "version",
		Short: "Show licorice version information",
		RunE: func(cmd *cobra.Command, args []string) error {
			verData := map[string]any{
				"Version":   constant.Version,
				"GoVer":     runtime.Version(),
				"GitCommit": constant.GitCommit,
				"BuildTime": constant.BuildTime,
				"OSAndArch": fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
			}

			t := template.Must(template.New("").Parse(versionTmpl))
			buf := &bytes.Buffer{}

			if err := t.Execute(buf, verData); err != nil {
				return err
			}

			fmt.Print(buf.String())

			return nil
		},
	}

	srvCmd = &cobra.Command{
		Use:   "server",
		Short: "Run in server mode",
		Run: func(cmd *cobra.Command, args []string) {
			runServer()
		},
	}

	conCmd = &cobra.Command{
		Use:   "console",
		Short: "Run in console mode",
		Run: func(cmd *cobra.Command, args []string) {
			runConsole()
		},
	}

	utilCmd = &cobra.Command{
		Use:   "util",
		Short: "Miscellaneous tools",
	}

	encCmd = &cobra.Command{
		Use:   "encode",
		Short: "Unpadded alternate base64 encoding",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Print(RawURLEncoding(args[0]))
		},
	}
)

// Execute executes the root cmd.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	clashRule = config.DefaultClashRuleFile
	clashRulePath = config.DefaultClashRulePath

	rootCmd.AddCommand(verCmd)
	rootCmd.AddCommand(srvCmd, conCmd)
	rootCmd.AddCommand(utilCmd)
	rootCmd.SetVersionTemplate(fmt.Sprintf("licorice version %s, build %s", constant.Version, constant.GitCommit))

	srvCmd.Flags().Uint16VarP(&port, "port", "p", 6060, "Set server port")
	srvCmd.Flags().StringVarP(&confDir, "confdir", "d", config.GetDefaultConfigDirectory(), "Specify the acl files directory")

	conCmd.Flags().StringVarP(&client, "client", "c", "clash", "Set target tunnel")
	conCmd.Flags().StringVarP(&confDir, "confdir", "d", config.GetDefaultConfigDirectory(), "Specify the acl files directory")
	conCmd.Flags().StringVarP(&inputFile, "input", "i", "stdin", "Read subscription from the specified file instead of stdin")
	conCmd.Flags().StringVarP(&outputFile, "output", "o", "stdout", "Write output to the specified file instead of stdout")

	utilCmd.AddCommand(encCmd)
}
