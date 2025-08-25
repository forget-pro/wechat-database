package chatlog

import (
	"encoding/json"
	"fmt"
	"runtime"

	"github.com/sjzar/chatlog/internal/chatlog"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(serverCmd)
	serverCmd.Flags().StringVarP(&serverAddr, "addr", "a", "127.0.0.1:5030", "server address")
	serverCmd.Flags().StringVarP(&serverDataDir, "data-dir", "d", "", "data dir")
	serverCmd.Flags().StringVarP(&serverWorkDir, "work-dir", "w", "", "work dir")
	serverCmd.Flags().StringVarP(&serverPlatform, "platform", "p", runtime.GOOS, "platform")
	serverCmd.Flags().IntVarP(&serverVer, "version", "v", 3, "version")
	serverCmd.Flags().BoolVarP(&serverJSON, "json", "j", false, "output result in JSON format")
	serverCmd.Flags().BoolVarP(&serverBackground, "background", "b", false, "run server in background")
}

var (
	serverAddr       string
	serverDataDir    string
	serverWorkDir    string
	serverPlatform   string
	serverVer        int
	serverJSON       bool
	serverBackground bool
)

type ServerResult struct {
	Success    bool   `json:"success"`
	Message    string `json:"message"`
	Address    string `json:"address,omitempty"`
	DataDir    string `json:"data_dir,omitempty"`
	WorkDir    string `json:"work_dir,omitempty"`
	Platform   string `json:"platform,omitempty"`
	Version    int    `json:"version,omitempty"`
	Background bool   `json:"background,omitempty"`
}

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start HTTP server",
	Run: func(cmd *cobra.Command, args []string) {
		result := ServerResult{
			Address:    serverAddr,
			DataDir:    serverDataDir,
			WorkDir:    serverWorkDir,
			Platform:   serverPlatform,
			Version:    serverVer,
			Background: serverBackground,
		}

		m, err := chatlog.New("")
		if err != nil {
			result.Success = false
			result.Message = fmt.Sprintf("failed to create chatlog instance: %v", err)
			outputResult(result, serverJSON)
			return
		}

		if err := m.CommandHTTPServer(serverAddr, serverDataDir, serverWorkDir, serverPlatform, serverVer, serverBackground, serverJSON); err != nil {
			result.Success = false
			result.Message = fmt.Sprintf("failed to start server: %v", err)
			outputResult(result, serverJSON)
			return
		}

		result.Success = true
		result.Message = fmt.Sprintf("HTTP server started successfully on %s", serverAddr)
		outputResult(result, serverJSON)
	},
}

func outputResult(result ServerResult, useJSON bool) {
	if useJSON {
		data, _ := json.MarshalIndent(result, "", "  ")
		fmt.Println(string(data))
	} else {
		if result.Success {
			fmt.Printf("✓ %s\n", result.Message)
		} else {
			fmt.Printf("✗ %s\n", result.Message)
		}
	}
}
