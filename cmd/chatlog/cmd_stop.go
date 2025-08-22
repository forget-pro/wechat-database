package chatlog

import (
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
	rootCmd.AddCommand(exitCmd)
}

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "停止聊天日志服务",
	Long:  `停止正在运行的聊天日志服务进程`,
	Example: `chatlog stop`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := stopChatlogProcess(); err != nil {
			fmt.Printf("停止服务失败: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("聊天日志服务已成功停止")
	},
}

var exitCmd = &cobra.Command{
	Use:   "exit",
	Short: "退出聊天日志程序",
	Long:  `安全退出聊天日志程序，停止所有服务并清理资源`,
	Example: `chatlog exit`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := stopChatlogProcess(); err != nil {
			log.Debug().Err(err).Msg("停止服务时出现错误")
			fmt.Println("警告: 停止服务时出现错误，但程序将退出")
		} else {
			fmt.Println("聊天日志服务已停止")
		}
		fmt.Println("聊天日志程序已退出")
		os.Exit(0)
	},
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
	// 查找chatlog进程
	cmd := exec.Command("tasklist", "/FI", "IMAGENAME eq chatlog.exe", "/FO", "CSV")
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("查找进程失败: %v", err)
	}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.Contains(line, "chatlog.exe") {
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

				// 终止进程
				killCmd := exec.Command("taskkill", "/PID", strconv.Itoa(pid), "/F")
				if err := killCmd.Run(); err != nil {
					return fmt.Errorf("终止进程 %d 失败: %v", pid, err)
				}
				return nil
			}
		}
	}
	return fmt.Errorf("未找到正在运行的chatlog进程")
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
