package shell

import (
    "fmt"
    "os"
    "strings"
    "runtime"
    "os/exec"
    "bufio"
    "unicode/utf8"

    "github.com/fatih/color"
    "github.com/logeshkannan96/dbcli/internal/database"
)

var queryHistory []string

var (
    promptColor    = color.New(color.FgCyan, color.Bold)
    errorColor     = color.New(color.FgRed, color.Bold)
    successColor   = color.New(color.FgGreen)
    columnColor    = color.New(color.FgYellow, color.Bold)
    resultColor    = color.New(color.FgWhite)
    historyColor   = color.New(color.FgMagenta)
    welcomeColor   = color.New(color.FgHiCyan, color.Bold)
)

func ClearScreen() {
    switch runtime.GOOS {
    case "linux", "darwin":
        fmt.Print("\033[2J\033[H")
    case "windows":
        cmd := exec.Command("cmd", "/c", "cls")
        cmd.Stdout = os.Stdout
        cmd.Run()
    }
}

func AddToHistory(query string) {
    queryHistory = append(queryHistory, query)
    if len(queryHistory) > 50 {
        queryHistory = queryHistory[1:]
    }
}

func PrintHistory() {
    historyColor.Println("")
    historyColor.Println("Query History:")
    for i, query := range queryHistory {
        historyColor.Printf("%d: %s\n", i+1, query)
    }
}

func StartShell() {
    welcomeColor.Println("Welcome to the MySQL shell. Type 'exit' to quit.")
    
    reader := bufio.NewReader(os.Stdin)

    for {
        promptColor.Print("dbcli> ")
        input, err := reader.ReadString('\n')
        if err != nil {
            errorColor.Printf("Error reading input: %v\n", err)
            continue
        }

        input = strings.TrimSpace(input)

        if input == "" {
            continue
        }

        if input == "exit" {
            successColor.Println("Goodbye!")
            return
        }

        if input == "clear" {
            ClearScreen()
            continue
        }

        if input == "history" {
            PrintHistory()
            successColor.Println("\nPress Enter to continue...")
            reader.ReadString('\n')
            ClearScreen()
            continue
        }

        results, err := database.ExecuteQuery(input)
        if err != nil {
            errorColor.Printf("Error: %v\n", err)
            continue
        }

        if len(results) == 0 {
            successColor.Println("Query executed successfully. No results returned.")
            continue
        }

        AddToHistory(input)
        PrintResults(results)
    }
}

func PrintResults(results []map[string]interface{}) {
    var columns []string
    columnWidths := make(map[string]int)

    for col := range results[0] {
        columns = append(columns, col)
        columnWidths[col] = utf8.RuneCountInString(col)
    }

    for _, row := range results {
        for col, val := range row {
            width := 0
            switch v := val.(type) {
            case []byte:
                width = utf8.RuneCountInString(string(v))
            default:
                width = utf8.RuneCountInString(fmt.Sprintf("%v", v))
            }
            if width > columnWidths[col] {
                columnWidths[col] = width
            }
        }
    }

    // Print column names
    for _, col := range columns {
        format := fmt.Sprintf("%%-%ds", columnWidths[col]+2)
        columnColor.Printf(format, col)
    }
    fmt.Println()

    // Print separator
    for _, col := range columns {
        fmt.Print(strings.Repeat("-", columnWidths[col]+2))
    }
    fmt.Println()

    // Print rows
    for _, row := range results {
        for _, col := range columns {
            val := row[col]
            format := fmt.Sprintf("%%-%dv", columnWidths[col]+2)
            switch v := val.(type) {
            case []byte:
                resultColor.Printf(format, string(v))
            default:
                resultColor.Printf(format, v)
            }
        }
        fmt.Println()
    }
    fmt.Println()
}