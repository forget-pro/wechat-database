package chatlog

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(stopCmd)
	stopCmd.Flags().BoolVarP(&stopJSON, "json", "j", false, "output result in JSON format")
	
	rootCmd.AddCommand(exitCmd)
	exitCmd.Flags().BoolVarP(&exitJSON, "json", "j", false, "output result in JSON format")
}

var (
	stopJSON bool
	exitJSON bool
)

// StopResponse 停止命令的响应结构
type StopResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Action  string `json:"action"` // "stop" 或 "exit"
}

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "停止聊天日志服务",
	Long:  `停止正在运行的聊天日志服务进程`,
	Example: `chatlog stop`,
	Run: func(cmd *cobra.Command, args []string) {
		response := StopResponse{
			Action: "stop",
		}

		if err := stopChatlogProcess(); err != nil {
			// 检查是否是"未找到进程"的错误
			if strings.Contains(err.Error(), "未找到正在运行的chatlog进程") {
				response.Success = true
				response.Message = "没有正在运行的聊天日志服务"
			} else {
				response.Success = false
				response.Message = fmt.Sprintf("停止服务失败: %v", err)
				outputStopResult(response, stopJSON)
				os.Exit(1)
			}
		} else {
			response.Success = true
			response.Message = "聊天日志服务已成功停止"
		}
		
		outputStopResult(response, stopJSON)
	},
}

var exitCmd = &cobra.Command{
	Use:   "exit",
	Short: "退出聊天日志程序",
	Long:  `安全退出聊天日志程序，停止所有服务并清理资源`,
	Example: `chatlog exit`,
	Run: func(cmd *cobra.Command, args []string) {
		response := StopResponse{
			Action: "exit",
		}

		if err := stopChatlogProcess(); err != nil {
			// 检查是否是"未找到进程"的错误
			if strings.Contains(err.Error(), "未找到正在运行的chatlog进程") {
				response.Success = true
				response.Message = "没有正在运行的聊天日志服务，程序已退出"
			} else {
				log.Debug().Err(err).Msg("停止服务时出现错误")
				response.Success = true // 即使停止失败，exit 也应该成功退出
				response.Message = "警告: 停止服务时出现错误，但程序已退出"
			}
		} else {
			response.Success = true
			response.Message = "聊天日志服务已停止，程序已退出"
		}
		
		outputStopResult(response, exitJSON)
		os.Exit(0)
	},
}

// outputStopResult 输出停止命令的结果
func outputStopResult(response StopResponse, useJSON bool) {
	if useJSON {
		data, _ := json.MarshalIndent(response, "", "  ")
		fmt.Println(string(data))
	} else {
		if response.Success {
			fmt.Printf("✓ %s\n", response.Message)
		} else {
			fmt.Printf("✗ %s\n", response.Message)
		}
	}
}

// stopChatlogProcess 查找并停止聊天日志进程
func stopChatlogProcess() error {
	if runtime.GOOS == "windows" {
		return stopProcessWindows()
	} else {
		return stopProcessUnix()
	}
}

// stopProcessWindows 在Windows上停止进程
func stopProcessWindows() error {
	// 查找所有 chatlog 相关进程 (包括 chatlog.exe 和 chatlog-daemon-*.exe)
	cmd := exec.Command("tasklist", "/FO", "CSV")
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("查找进程失败: %v", err)
	}

	lines := strings.Split(string(output), "\n")
	var foundPID int = -1
	
	for _, line := range lines {
		// 查找包含 chatlog 的进程名
		if strings.Contains(line, "chatlog") && (strings.Contains(line, ".exe") || strings.Contains(line, "chatlog-daemon")) {
			// 解析PID
			fields := strings.Split(line, ",")
			if len(fields) >= 2 {
				pidStr := strings.Trim(fields[1], "\"")
				pid, err := strconv.Atoi(pidStr)
				if err != nil {
					continue
				}

				// 排除当前进程
				if pid == os.Getpid() {
					continue
				}

				foundPID = pid
				break
			}
		}
	}
	
	if foundPID == -1 {
		return fmt.Errorf("未找到正在运行的chatlog进程")
	}

	// 终止进程
	killCmd := exec.Command("taskkill", "/PID", strconv.Itoa(foundPID), "/F")
	if err := killCmd.Run(); err != nil {
		return fmt.Errorf("终止进程 %d 失败: %v", foundPID, err)
	}
	return nil
}

// stopProcessUnix 在Unix系统上停止进程
func stopProcessUnix() error {
	// 查找chatlog进程
	cmd := exec.Command("pgrep", "-f", "chatlog")
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("查找进程失败: %v", err)
	}

	pids := strings.Fields(string(output))
	currentPid := strconv.Itoa(os.Getpid())

	for _, pidStr := range pids {
		if pidStr == currentPid {
			continue // 跳过当前进程
		}

		// 发送TERM信号
		killCmd := exec.Command("kill", "-TERM", pidStr)
		if err := killCmd.Run(); err != nil {
			// 如果TERM失败，使用KILL
			killCmd = exec.Command("kill", "-KILL", pidStr)
			if err := killCmd.Run(); err != nil {
				return fmt.Errorf("终止进程 %s 失败: %v", pidStr, err)
			}
		}
		return nil
	}
	return fmt.Errorf("未找到正在运行的chatlog进程")
}
