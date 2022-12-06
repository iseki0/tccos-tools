package cosupload

import (
	"context"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/tencentyun/cos-go-sdk-v5"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"tccos-tools/config"
	"tccos-tools/exitcode"
)

func Cmd() *cobra.Command {
	var c cobra.Command
	c.Use = "cosupload FILE... DIR"
	c.Args = cobra.MinimumNArgs(2)
	c.Run = run
	return &c
}

func run(cmd *cobra.Command, args []string) {
	var ctx = context.Background()
	var files = args[:len(args)-1]
	var target = args[len(args)-1]
	configu, e := config.Parse(strings.TrimSpace(os.Getenv("COS_BASE64")))
	if e != nil {
		fmt.Fprintln(os.Stderr, e.Error())
		exitcode.Set(1)
		return
	}
	var client = cos.NewClient(configu.BaseURL, &http.Client{Transport: configu.Auth})
	for _, file := range files {
		var targetPath string
		if filepath.IsAbs(file) {
			targetPath = filepath.Join(target, filepath.Base(file))
		} else {
			targetPath = filepath.Join(target, file)
		}
		_, response, e := client.Object.Upload(ctx, targetPath, file, &cos.MultiUploadOptions{
			ThreadPoolSize: 4,
		})
		if e != nil {
			fmt.Fprintln(os.Stderr, e.Error())
			exitcode.Set(1)
			return
		}
		if response.StatusCode >= 400 {
			fmt.Fprintln(os.Stderr, fmt.Sprintf("HTTP status: %d", response.StatusCode))
			exitcode.Set(1)
			return
		}
	}
}
