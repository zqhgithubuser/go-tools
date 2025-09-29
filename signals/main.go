package main

import (
    "context"
    "fmt"
    "log"
    "os"
    "os/signal"
    "path/filepath"
    "strconv"
    "syscall"
    "time"
)

func main() {
    tmpDir := createTempDir()
    defer cleanup(tmpDir)

    ctx, cancel := context.WithCancel(context.Background())
    // 表示 goroutine 成功结束
    done := make(chan struct{})

    go createFiles(ctx, tmpDir, done)
    handleSignals(cancel, done)

    fmt.Println("Program exited")
}

// 创建临时目录
func createTempDir() string {
    tmpDir, err := os.MkdirTemp("", "app_*")
    if err != nil {
        log.Fatal("Failed to create temporary directory:", err)
    }
    fmt.Println("Temporary directory created:", tmpDir)
    return tmpDir
}

// 清理临时目录
func cleanup(tmpDir string) {
    if err := os.RemoveAll(tmpDir); err != nil {
        fmt.Println("Cleanup failed:", err)
    }
    fmt.Println("Cleanup completed")
}

// 在临时目录下创建文件
func createFiles(ctx context.Context, tmpDir string, done chan struct{}) {
    defer close(done)
    for i := 0; i < 30; i++ {
        if ctx.Err() != nil {
            fmt.Println("File creation canceled, stopping")
            return
        }
        filePath := filepath.Join(tmpDir, strconv.Itoa(i))
        f, err := os.Create(filePath)
        if err != nil {
            panic("Failed to create file: " + err.Error())
        }
        f.Close()
        fmt.Println("Created file:", filePath)
        time.Sleep(1 * time.Second)
    }
}

// 处理信号
func handleSignals(cancel context.CancelFunc, done chan struct{}) {
    sigs := make(chan os.Signal, 1)
    signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

    for {
        sig := <-sigs // 收到信号
        switch sig {
        case syscall.SIGINT, syscall.SIGTERM:
            fmt.Println("Signal received:", sig, "- exiting gracefully")
            cancel()
            <-done
            return
        case syscall.SIGQUIT:
            fmt.Println("Signal received: SIGQUIT - panicking after cleanup")
            cancel()
            <-done
            panic("SIGQUIT received")
        default:
            fmt.Println("Unknown signal received:", sig)
        }
    }
}
