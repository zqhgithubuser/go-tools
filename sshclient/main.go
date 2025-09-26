package main

import (
    "flag"
    "fmt"
    "golang.org/x/crypto/ssh"
    "golang.org/x/term"
    "log"
    "os"
    "time"
)

func main() {
    host := flag.String("host", "", "主机 (必填)")
    port := flag.Int("port", 22, "端口")
    user := flag.String("user", "", "用户名 (必填)")
    password := flag.String("password", "", "密码")
    cmd := flag.String("cmd", "uname -a", "执行的命令")
    flag.Parse()

    if *host == "" || *user == "" {
        flag.Usage()
        os.Exit(2)
    }

    pass := *password
    if pass == "" {
        pass = readPassword()
    }

    config := &ssh.ClientConfig{
        User:            *user,
        Auth:            []ssh.AuthMethod{ssh.Password(pass)},
        HostKeyCallback: ssh.InsecureIgnoreHostKey(),
        Timeout:         5 * time.Second,
    }

    addr := fmt.Sprintf("%s:%d", *host, *port)
    client, err := ssh.Dial("tcp", addr, config)
    if err != nil {
        log.Fatalf("连接失败: %v", err)
    }
    defer client.Close()

    sess, err := client.NewSession()
    if err != nil {
        log.Fatalf("创建会话失败: %v", err)
    }
    defer sess.Close()

    output, err := sess.CombinedOutput(*cmd)
    if err != nil {
        log.Fatalf("执行命令失败: %v\n输出: %s", err, output)
    }
    fmt.Print(string(output))
}

func readPassword() string {
    fmt.Print("密码: ")
    p, err := term.ReadPassword(int(os.Stdin.Fd()))
    fmt.Println()
    if err != nil {
        log.Fatalf("读取密码失败: %v", err)
    }
    return string(p)
}
