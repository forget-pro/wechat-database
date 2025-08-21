package chatlog

import (
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/sjzar/chatlog/internal/chatlog"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

func init() {
	// Status command
	rootCmd.AddCommand(statusCmd)
	statusCmd.Flags().BoolVarP(&statusJSON, "json", "j", false, "output in JSON format")

	// Get secret command
	rootCmd.AddCommand(getSecretCmd)
	getSecretCmd.Flags().IntVarP(&getSecretPID, "pid", "p", 0, "Process ID of WeChat instance")
	getSecretCmd.Flags().BoolVarP(&getSecretJSON, "json", "j", false, "output in JSON format")

	// Decrypt data command
	rootCmd.AddCommand(decryptDataCmd)
	decryptDataCmd.Flags().StringVarP(&decryptDataDir, "data-dir", "d", "", "WeChat data directory")
	decryptDataCmd.Flags().StringVarP(&decryptWorkDir, "work-dir", "w", "", "Work directory for decrypted files")
	decryptDataCmd.Flags().StringVarP(&decryptKey, "key", "k", "", "Encryption key")
	decryptDataCmd.Flags().StringVarP(&decryptDataPlatform, "platform", "p", runtime.GOOS, "Platform (darwin/windows)")
	decryptDataCmd.Flags().IntVarP(&decryptDataVer, "version", "v", 3, "WeChat version (3 or 4)")
	decryptDataCmd.Flags().BoolVarP(&decryptDataJSON, "json", "j", false, "output in JSON format")

	// Auto decrypt command
	rootCmd.AddCommand(autoDecryptCmd)
	autoDecryptCmd.Flags().BoolVarP(&autoDecryptJSON, "json", "j", false, "output in JSON format")
	autoDecryptCmd.Flags().BoolVarP(&autoDecryptStop, "stop", "s", false, "stop auto decryption")

	// Work dir command
	rootCmd.AddCommand(workDirCmd)
	workDirCmd.Flags().StringVarP(&workDirPath, "path", "p", "", "Work directory path")
	workDirCmd.Flags().BoolVarP(&workDirJSON, "json", "j", false, "output in JSON format")

	// Set secret command
	rootCmd.AddCommand(setSecretCmd)
	setSecretCmd.Flags().StringVarP(&setSecretKey, "key", "k", "", "Encryption key to set")
	setSecretCmd.Flags().BoolVarP(&setSecretJSON, "json", "j", false, "output in JSON format")

	// Data dir command
	rootCmd.AddCommand(dataDirCmd)
	dataDirCmd.Flags().StringVarP(&dataDirPath, "path", "p", "", "Data directory path")
	dataDirCmd.Flags().BoolVarP(&dataDirJSON, "json", "j", false, "output in JSON format")
}

// Variables for all commands
var (
	// Status command variables
	statusJSON bool

	// Get secret command variables
	getSecretPID  int
	getSecretJSON bool

	// Decrypt data command variables
	decryptDataDir      string
	decryptWorkDir      string
	decryptKey          string
	decryptDataPlatform string
	decryptDataVer      int
	decryptDataJSON     bool

	// Auto decrypt command variables
	autoDecryptJSON bool
	autoDecryptStop bool

	// Work dir command variables
	workDirPath string
	workDirJSON bool

	// Set secret command variables
	setSecretKey  string
	setSecretJSON bool

	// Data dir command variables
	dataDirPath string
	dataDirJSON bool
)

// Response types
type StatusResponse struct {
	Account       string `json:"account,omitempty"`
	PID           int    `json:"pid,omitempty"`
	Status        string `json:"status,omitempty"`
	ExePath       string `json:"exe_path,omitempty"`
	Platform      string `json:"platform,omitempty"`
	Version       string `json:"version,omitempty"`
	Session       string `json:"session,omitempty"`
	DataKey       string `json:"data_key,omitempty"`
	DataUsage     string `json:"data_usage,omitempty"`
	DataDir       string `json:"data_dir,omitempty"`
	WorkUsage     string `json:"work_usage,omitempty"`
	WorkDir       string `json:"work_dir,omitempty"`
	HTTPServer    string `json:"http_server,omitempty"`
	AutoDecrypt   string `json:"auto_decrypt,omitempty"`
	HTTPAddr      string `json:"http_addr,omitempty"`
	Message       string `json:"message"`
	Success       bool   `json:"success"`
}

type GetSecretResponse struct {
	Key     string `json:"key,omitempty"`
	Message string `json:"message"`
	Success bool   `json:"success"`
}

type DecryptDataResponse struct {
	DataDir  string `json:"data_dir,omitempty"`
	WorkDir  string `json:"work_dir,omitempty"`
	Platform string `json:"platform,omitempty"`
	Version  int    `json:"version,omitempty"`
	Message  string `json:"message"`
	Success  bool   `json:"success"`
}

type AutoDecryptResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Success bool   `json:"success"`
}

type WorkDirResponse struct {
	WorkDir string `json:"work_dir,omitempty"`
	Message string `json:"message"`
	Success bool   `json:"success"`
}

type SetSecretResponse struct {
	Message string `json:"message"`
	Success bool   `json:"success"`
}

type DataDirResponse struct {
	DataDir string `json:"data_dir,omitempty"`
	Message string `json:"message"`
	Success bool   `json:"success"`
}

// Status command - get WeChat process information
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Get WeChat process information",
	Run: func(cmd *cobra.Command, args []string) {
		m, err := chatlog.New("")
		if err != nil {
			if statusJSON {
				resp := StatusResponse{Success: false, Message: err.Error()}
				jsonData, _ := json.Marshal(resp)
				fmt.Println(string(jsonData))
			} else {
				log.Err(err).Msg("failed to create chatlog instance")
			}
			return
		}

		// Get current status information
		status := m.GetCurrentStatus()
		
		// Format session time
		sessionStr := ""
		if !status.LastSession.IsZero() {
			sessionStr = status.LastSession.Format("2006-01-02 15:04:05")
		}

		// Format HTTP server status
		httpServerStr := "[未启动]"
		if status.HTTPEnabled {
			httpServerStr = fmt.Sprintf("[已启动] %s", status.HTTPAddr)
		}

		// Format auto decrypt status
		autoDecryptStr := "[未开启]"
		if status.AutoDecrypt {
			autoDecryptStr = "[已开启]"
		}

		// Handle empty values with defaults
		dataUsage := status.DataUsage
		if dataUsage == "" {
			dataUsage = "0 B"
		}
		workUsage := status.WorkUsage
		if workUsage == "" {
			workUsage = "0 B"
		}

		if statusJSON {
			resp := StatusResponse{
				Account:     status.Account,
				PID:         status.PID,
				Status:      status.Status,
				ExePath:     status.ExePath,
				Platform:    status.Platform,
				Version:     status.FullVersion,
				Session:     sessionStr,
				DataKey:     status.DataKey,
				DataUsage:   dataUsage,
				DataDir:     status.DataDir,
				WorkUsage:   workUsage,
				WorkDir:     status.WorkDir,
				HTTPServer:  httpServerStr,
				AutoDecrypt: autoDecryptStr,
				HTTPAddr:    status.HTTPAddr,
				Success:     true,
				Message:     "Status information retrieved successfully",
			}
			jsonData, err := json.Marshal(resp)
			if err != nil {
				resp := StatusResponse{Success: false, Message: err.Error()}
				jsonData, _ := json.Marshal(resp)
				fmt.Println(string(jsonData))
				return
			}
			fmt.Println(string(jsonData))
		} else {
			// Format output to match terminal UI display
			fmt.Printf("Account:      %-25s PID:          %-25d\n", status.Account, status.PID)
			fmt.Printf("Status:       %-25s ExePath:      %-25s\n", status.Status, status.ExePath)
			fmt.Printf("Platform:     %-25s Version:      %-25s\n", status.Platform, status.FullVersion)
			fmt.Printf("Session:      %-25s Data Key:     %-25s\n", sessionStr, status.DataKey)
			fmt.Printf("Data Usage:   %-25s Data Dir:     %-25s\n", dataUsage, status.DataDir)
			fmt.Printf("Work Usage:   %-25s Work Dir:     %-25s\n", workUsage, status.WorkDir)
			fmt.Printf("HTTP Server:  %-25s Auto Decrypt: %-25s\n", httpServerStr, autoDecryptStr)
		}
	},
}

// Get secret command - get data encryption key
var getSecretCmd = &cobra.Command{
	Use:   "get-secret",
	Short: "Get data encryption key",
	Run: func(cmd *cobra.Command, args []string) {
		m, err := chatlog.New("")
		if err != nil {
			if getSecretJSON {
				resp := GetSecretResponse{Success: false, Message: err.Error()}
				jsonData, _ := json.Marshal(resp)
				fmt.Println(string(jsonData))
			} else {
				log.Err(err).Msg("failed to create chatlog instance")
			}
			return
		}

		key, err := m.CommandKey(getSecretPID)
		if err != nil {
			if getSecretJSON {
				resp := GetSecretResponse{Success: false, Message: err.Error()}
				jsonData, _ := json.Marshal(resp)
				fmt.Println(string(jsonData))
			} else {
				log.Err(err).Msg("failed to get key")
			}
			return
		}

		if getSecretJSON {
			resp := GetSecretResponse{
				Key:     key,
				Success: true,
				Message: "Key retrieved successfully",
			}
			jsonData, err := json.Marshal(resp)
			if err != nil {
				resp := GetSecretResponse{Success: false, Message: err.Error()}
				jsonData, _ := json.Marshal(resp)
				fmt.Println(string(jsonData))
				return
			}
			fmt.Println(string(jsonData))
		} else {
			fmt.Println(key)
		}
	},
}

// Decrypt data command - decrypt WeChat database files
var decryptDataCmd = &cobra.Command{
	Use:   "decrypt-data",
	Short: "Decrypt WeChat database files",
	Run: func(cmd *cobra.Command, args []string) {
		m, err := chatlog.New("")
		if err != nil {
			if decryptDataJSON {
				resp := DecryptDataResponse{Success: false, Message: err.Error()}
				jsonData, _ := json.Marshal(resp)
				fmt.Println(string(jsonData))
			} else {
				log.Err(err).Msg("failed to create chatlog instance")
			}
			return
		}

		err = m.CommandDecrypt(decryptDataDir, decryptWorkDir, decryptKey, decryptDataPlatform, decryptDataVer)
		if err != nil {
			if decryptDataJSON {
				resp := DecryptDataResponse{Success: false, Message: err.Error()}
				jsonData, _ := json.Marshal(resp)
				fmt.Println(string(jsonData))
			} else {
				log.Err(err).Msg("failed to decrypt")
			}
			return
		}

		if decryptDataJSON {
			resp := DecryptDataResponse{
				DataDir:  decryptDataDir,
				WorkDir:  decryptWorkDir,
				Platform: decryptDataPlatform,
				Version:  decryptDataVer,
				Success:  true,
				Message:  "Data decryption completed successfully",
			}
			jsonData, err := json.Marshal(resp)
			if err != nil {
				resp := DecryptDataResponse{Success: false, Message: err.Error()}
				jsonData, _ := json.Marshal(resp)
				fmt.Println(string(jsonData))
				return
			}
			fmt.Println(string(jsonData))
		} else {
			fmt.Println("Data decryption completed successfully")
		}
	},
}

// Auto decrypt command - start or stop automatic decryption
var autoDecryptCmd = &cobra.Command{
	Use:   "auto-decrypt",
	Short: "Start or stop automatic decryption",
	Run: func(cmd *cobra.Command, args []string) {
		m, err := chatlog.New("")
		if err != nil {
			if autoDecryptJSON {
				resp := AutoDecryptResponse{Success: false, Message: err.Error()}
				jsonData, _ := json.Marshal(resp)
				fmt.Println(string(jsonData))
			} else {
				log.Err(err).Msg("failed to create chatlog instance")
			}
			return
		}

		if autoDecryptStop {
			// Stop auto decryption
			err := m.StopAutoDecrypt()
			if err != nil {
				if autoDecryptJSON {
					resp := AutoDecryptResponse{Success: false, Message: err.Error()}
					jsonData, _ := json.Marshal(resp)
					fmt.Println(string(jsonData))
				} else {
					log.Err(err).Msg("failed to stop auto decryption")
				}
				return
			}

			if autoDecryptJSON {
				resp := AutoDecryptResponse{
					Status:  "stopped",
					Success: true,
					Message: "Auto decryption stopped successfully",
				}
				jsonData, err := json.Marshal(resp)
				if err != nil {
					resp := AutoDecryptResponse{Success: false, Message: err.Error()}
					jsonData, _ := json.Marshal(resp)
					fmt.Println(string(jsonData))
					return
				}
				fmt.Println(string(jsonData))
			} else {
				fmt.Println("Auto decryption stopped successfully")
			}
		} else {
			// Start auto decryption
			err := m.StartAutoDecrypt()
			if err != nil {
				if autoDecryptJSON {
					resp := AutoDecryptResponse{Success: false, Message: err.Error()}
					jsonData, _ := json.Marshal(resp)
					fmt.Println(string(jsonData))
				} else {
					log.Err(err).Msg("failed to start auto decryption")
				}
				return
			}

			if autoDecryptJSON {
				resp := AutoDecryptResponse{
					Status:  "started",
					Success: true,
					Message: "Auto decryption started successfully. Press Ctrl+C to stop.",
				}
				jsonData, err := json.Marshal(resp)
				if err != nil {
					resp := AutoDecryptResponse{Success: false, Message: err.Error()}
					jsonData, _ := json.Marshal(resp)
					fmt.Println(string(jsonData))
					return
				}
				fmt.Println(string(jsonData))
			} else {
				fmt.Println("Auto decryption started successfully. Press Ctrl+C to stop.")
			}

			// Wait for interrupt signal
			c := make(chan os.Signal, 1)
			signal.Notify(c, os.Interrupt, syscall.SIGTERM)
			<-c

			// Stop auto decryption on exit
			err = m.StopAutoDecrypt()
			if err != nil {
				if !autoDecryptJSON {
					log.Err(err).Msg("failed to stop auto decryption on exit")
				}
			} else {
				if !autoDecryptJSON {
					fmt.Println("\nAuto decryption stopped")
				}
			}
		}
	},
}

// Work dir command - set or get working directory
var workDirCmd = &cobra.Command{
	Use:   "work-dir",
	Short: "Set working directory",
	Run: func(cmd *cobra.Command, args []string) {
		m, err := chatlog.New("")
		if err != nil {
			if workDirJSON {
				resp := WorkDirResponse{Success: false, Message: err.Error()}
				jsonData, _ := json.Marshal(resp)
				fmt.Println(string(jsonData))
			} else {
				log.Err(err).Msg("failed to create chatlog instance")
			}
			return
		}

		if workDirPath == "" {
			if workDirJSON {
				resp := WorkDirResponse{
					Success: false,
					Message: "Work directory path is required. Use --path flag",
				}
				jsonData, _ := json.Marshal(resp)
				fmt.Println(string(jsonData))
			} else {
				fmt.Println("Work directory path is required. Use --path flag")
			}
			return
		}

		err = m.SetWorkDir(workDirPath)
		if err != nil {
			if workDirJSON {
				resp := WorkDirResponse{Success: false, Message: err.Error()}
				jsonData, _ := json.Marshal(resp)
				fmt.Println(string(jsonData))
			} else {
				log.Err(err).Msg("failed to set work directory")
			}
			return
		}

		if workDirJSON {
			resp := WorkDirResponse{
				WorkDir: workDirPath,
				Success: true,
				Message: "Work directory set successfully",
			}
			jsonData, err := json.Marshal(resp)
			if err != nil {
				resp := WorkDirResponse{Success: false, Message: err.Error()}
				jsonData, _ := json.Marshal(resp)
				fmt.Println(string(jsonData))
				return
			}
			fmt.Println(string(jsonData))
		} else {
			fmt.Printf("Work directory set to: %s\n", workDirPath)
		}
	},
}

// Set secret command - set data encryption key
var setSecretCmd = &cobra.Command{
	Use:   "set-secret",
	Short: "Set data encryption key",
	Run: func(cmd *cobra.Command, args []string) {
		m, err := chatlog.New("")
		if err != nil {
			if setSecretJSON {
				resp := SetSecretResponse{Success: false, Message: err.Error()}
				jsonData, _ := json.Marshal(resp)
				fmt.Println(string(jsonData))
			} else {
				log.Err(err).Msg("failed to create chatlog instance")
			}
			return
		}

		if setSecretKey == "" {
			if setSecretJSON {
				resp := SetSecretResponse{
					Success: false,
					Message: "Encryption key is required. Use --key flag",
				}
				jsonData, _ := json.Marshal(resp)
				fmt.Println(string(jsonData))
			} else {
				fmt.Println("Encryption key is required. Use --key flag")
			}
			return
		}

		err = m.SetDataKey(setSecretKey)
		if err != nil {
			if setSecretJSON {
				resp := SetSecretResponse{Success: false, Message: err.Error()}
				jsonData, _ := json.Marshal(resp)
				fmt.Println(string(jsonData))
			} else {
				log.Err(err).Msg("failed to set encryption key")
			}
			return
		}

		if setSecretJSON {
			resp := SetSecretResponse{
				Success: true,
				Message: "Encryption key set successfully",
			}
			jsonData, err := json.Marshal(resp)
			if err != nil {
				resp := SetSecretResponse{Success: false, Message: err.Error()}
				jsonData, _ := json.Marshal(resp)
				fmt.Println(string(jsonData))
				return
			}
			fmt.Println(string(jsonData))
		} else {
			fmt.Println("Encryption key set successfully")
		}
	},
}

// Data dir command - set data directory
var dataDirCmd = &cobra.Command{
	Use:   "data-dir",
	Short: "Set data directory",
	Run: func(cmd *cobra.Command, args []string) {
		m, err := chatlog.New("")
		if err != nil {
			if dataDirJSON {
				resp := DataDirResponse{Success: false, Message: err.Error()}
				jsonData, _ := json.Marshal(resp)
				fmt.Println(string(jsonData))
			} else {
				log.Err(err).Msg("failed to create chatlog instance")
			}
			return
		}

		if dataDirPath == "" {
			if dataDirJSON {
				resp := DataDirResponse{
					Success: false,
					Message: "Data directory path is required. Use --path flag",
				}
				jsonData, _ := json.Marshal(resp)
				fmt.Println(string(jsonData))
			} else {
				fmt.Println("Data directory path is required. Use --path flag")
			}
			return
		}

		err = m.SetDataDir(dataDirPath)
		if err != nil {
			if dataDirJSON {
				resp := DataDirResponse{Success: false, Message: err.Error()}
				jsonData, _ := json.Marshal(resp)
				fmt.Println(string(jsonData))
			} else {
				log.Err(err).Msg("failed to set data directory")
			}
			return
		}

		if dataDirJSON {
			resp := DataDirResponse{
				DataDir: dataDirPath,
				Success: true,
				Message: "Data directory set successfully",
			}
			jsonData, err := json.Marshal(resp)
			if err != nil {
				resp := DataDirResponse{Success: false, Message: err.Error()}
				jsonData, _ := json.Marshal(resp)
				fmt.Println(string(jsonData))
				return
			}
			fmt.Println(string(jsonData))
		} else {
			fmt.Printf("Data directory set to: %s\n", dataDirPath)
		}
	},
}
