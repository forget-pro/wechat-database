package chatlog

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

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

	// 添加守护进程命令
	rootCmd.AddCommand(serverDaemonCmd)
	serverDaemonCmd.Flags().StringVarP(&serverAddr, "addr", "a", "127.0.0.1:5030", "server address")
	serverDaemonCmd.Flags().StringVarP(&serverDataDir, "data-dir", "d", "", "data dir")
	serverDaemonCmd.Flags().StringVarP(&serverWorkDir, "work-dir", "w", "", "work dir")
	serverDaemonCmd.Flags().StringVarP(&serverPlatform, "platform", "p", runtime.GOOS, "platform")
	serverDaemonCmd.Flags().IntVarP(&serverVer, "version", "v", 3, "version")
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

		// 如果是后台模式，创建独立的后台进程
		if serverBackground {
			if err := startBackgroundServer(result); err != nil {
				result.Success = false
				result.Message = fmt.Sprintf("failed to start background server: %v", err)
				outputResult(result, serverJSON)
				return
			}
			result.Success = true
			result.Message = fmt.Sprintf("HTTP server started successfully on %s", serverAddr)
			outputResult(result, serverJSON)
			return
		}

		// 前台模式
		m, err := chatlog.New("")
		if err != nil {
			result.Success = false
			result.Message = fmt.Sprintf("failed to create chatlog instance: %v", err)
			outputResult(result, serverJSON)
			return
		}

		if err := m.CommandHTTPServer(serverAddr, serverDataDir, serverWorkDir, serverPlatform, serverVer, false, serverJSON); err != nil {
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

// serverDaemonCmd 是实际运行的守护进程命令
var serverDaemonCmd = &cobra.Command{
	Use:    "server-daemon",
	Short:  "Run HTTP server daemon (internal use)",
	Hidden: true, // 隐藏此命令，不在帮助中显示
	Run: func(cmd *cobra.Command, args []string) {
		m, err := chatlog.New("")
		if err != nil {
			fmt.Printf("failed to create chatlog instance: %v\n", err)
			os.Exit(1)
		}

		// 运行前台服务器（在守护进程中）
		if err := m.CommandHTTPServer(serverAddr, serverDataDir, serverWorkDir, serverPlatform, serverVer, false, false); err != nil {
			fmt.Printf("failed to start server: %v\n", err)
			os.Exit(1)
		}
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

// startBackgroundServer 启动后台服务器进程
func startBackgroundServer(result ServerResult) error {
	// 构建命令参数
	args := []string{
		"server-daemon", // 使用特殊的守护进程命令
		"-a", result.Address,
		"-p", result.Platform,
		"-v", fmt.Sprintf("%d", result.Version),
	}
	
	if result.DataDir != "" {
		args = append(args, "-d", result.DataDir)
	}
	
	if result.WorkDir != "" {
		args = append(args, "-w", result.WorkDir)
	}

	var cmd *exec.Cmd
	
	// 获取当前可执行文件路径
	executable, err := os.Executable()
	if err != nil || strings.Contains(executable, "go-build") {
		// 如果是 go run 模式，先编译后台程序
		// 创建临时目录
		tempDir := os.TempDir()
		tempExe := fmt.Sprintf("%s\\chatlog-daemon-%d.exe", tempDir, time.Now().Unix())
		
		// 编译程序
		buildCmd := exec.Command("go", "build", "-o", tempExe, ".")
		buildCmd.Dir, _ = os.Getwd()
		if err := buildCmd.Run(); err != nil {
			// 如果编译失败，回退到 go run 模式
			cmd = exec.Command("go", append([]string{"run", "."}, args...)...)
			cmd.Dir, _ = os.Getwd()
			cmd.Stderr = nil // 抑制错误输出
		} else {
			// 使用编译后的可执行文件
			cmd = exec.Command(tempExe, args...)
		}
	} else {
		// 使用可执行文件启动后台进程
		cmd = exec.Command(executable, args...)
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start background process: %v", err)
	}

	// 等待一小段时间让服务启动
	time.Sleep(1 * time.Second)
	
	// 检查服务是否真的启动了
	return checkHTTPServiceRunning(result.Address)
}

// checkHTTPServiceRunning 检查 HTTP 服务是否在运行
func checkHTTPServiceRunning(addr string) error {
	for i := 0; i < 10; i++ { // 最多等待 5 秒
		conn, err := net.DialTimeout("tcp", addr, 500*time.Millisecond)
		if err == nil {
			conn.Close()
			return nil
		}
		time.Sleep(500 * time.Millisecond)
	}
	return fmt.Errorf("HTTP service failed to start or is not accessible")
}

// getHost 从地址中提取主机名
func getHost(addr string) string {
	parts := strings.Split(addr, ":")
	if len(parts) >= 2 {
		return parts[0]
	}
	return addr
}

// getPort 从地址中提取端口
func getPort(addr string) string {
	parts := strings.Split(addr, ":")
	if len(parts) >= 2 {
		return parts[1]
	}
	return ""
}
