/*
Copyright © 2026 ANTONIO RODRIGUEZ <kontakt@antoniorodriguez.no>
*/
package cmd

import (
	"bufio"
	"fmt"
	fzf "github.com/junegunn/fzf/src"
	"github.com/spf13/cobra"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
)

var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search users in the Folkomaten test database",
	Long: `Search users in the Folkomaten test database using fzf.

	You can copy information of the selected user to the clipboard by pressing Enter and then choose what to copy`,
	RunE: func(cmd *cobra.Command, args []string) error {
		databasePath, err := cmd.Flags().GetString("file")

		if err != nil {
			fmt.Println("Error reading flag:", err)
			return err
		}

		if databasePath == "" {
			databasePath = loadConfig()
		}

		getTestUsers(databasePath)
		return nil
	},
}

func init() {
	// FIX: The flag is not working as intended
	// TODO: Find a way to change the "active database" from the config
	// TODO: Parse UTF-16 files
	searchCmd.Flags().StringP("file", "f", "", "change the database file to populate the folkctl database with")
	rootCmd.AddCommand(searchCmd)
}

func getTestUsers(filePath string) {
	filePath, err := filepath.Abs(filePath)

	if err != nil {
		fmt.Println("Error resolving resource path:", err)
		return
	}

	content, err := os.ReadFile(filePath)
	clipboardCmd := detectClipboardCmd()

	selected, err := startFzf(content)

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	if selected == "" {
		fmt.Println("No selection made")
		return
	}

	value := readInput(selected)
	if value == "" {
		fmt.Println("Nothing to copy")
		return
	}

	err = copyToClipboard(value, clipboardCmd)

	if err != nil {
		fmt.Printf("Clipboard failed. Value: %s\n", value)
	} else {
		fmt.Printf("Copied: %s\n", value)
	}
}

func detectClipboardCmd() string {
	return detectClipboardCmdForOS(runtime.GOOS)
}

func detectClipboardCmdForOS(goos string) string {
	switch goos {
	case "linux":
		return "xclip -selection clipboard"
	case "windows":
		return "clip"
	default:
		return "pbcopy"
	}
}

func readInput(selected string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Copy (f)nr, (n)ame, or (a)ll? ")
	choice, _ := reader.ReadString('\n')
	choice = strings.TrimSpace(strings.ToLower(choice))

	fields := strings.Split(selected, ",")

	if len(fields) < 2 {
		fmt.Fprintf(os.Stderr, "Invalid line format: %s\n", selected)
		return ""
	}

	var value string
	switch choice {
	case "f":
		value = fields[0]
	case "n":
		value = fields[1]
	case "a":
		value = selected
	default:
		fmt.Println("Invalid choice, copying all")
		value = selected
	}

	return value

}

func copyToClipboard(value, cmd string) error {
	parts := strings.Fields(cmd)
	c := exec.Command(parts[0], parts[1:]...)
	stdin, err := c.StdinPipe()

	if err != nil {
		return err
	}

	go func() {
		io.WriteString(stdin, value)
		stdin.Close()
	}()

	return c.Run()
}

func startFzf(content []byte) (string, error) {
	wg := sync.WaitGroup{}
	inputChan := make(chan string)
	go func() {
		contentLines := strings.Split(string(content), "\n")
		for _, user := range contentLines {
			if user != "" {
				inputChan <- user
			}
		}
		close(inputChan)
	}()

	var selected string
	outputChan := make(chan string)
	wg.Go(func() {
		for s := range outputChan {
			selected = s
		}
	})

	options, err := fzf.ParseOptions(
		true,
		[]string{
			"--reverse", "--border", "--height=40%",
			"--header", "Press Enter to select, then choose what to copy\n\n",
		},
	)

	if err != nil {
		return "", fmt.Errorf("fzf parse options: %w", err)
	}

	options.Input = inputChan
	options.Output = outputChan
	code, err := fzf.Run(options)
	close(outputChan)
	wg.Wait()

	if code == 130 {
		return "", fmt.Errorf("Search cancelled")
	}

	if code != 0 || err != nil {
		return "", fmt.Errorf("fzf run: code %d, %w", code, err)
	}

	return selected, nil
}
