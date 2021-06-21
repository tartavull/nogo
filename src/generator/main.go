package main

import (
    "bufio"
    "os"
    "fmt"
    "os/exec"
    "bytes"
)

func check(e error) {
    if e != nil {
        panic(e)
    }
}

// https://github.com/opensource-apple/xnu/blob/master/bsd/kern/syscalls.master#L41
const (
    EXIT = 0x02000001
    FORK = 0x02000002
    READ = 0x02000003
    WRITE = 0x02000004
    OPEN = 0x02000005
    CLOSE = 0x02000006
    WAIT4 = 0x02000007
    ENOYSYS = 0x02000008
    LINK = 0x02000009
    UNLINK = 0x0200000A
)

func syscall_exit(code int) string {
    return fmt.Sprintf(`
        mov       rax, %#x        ; system call for exit
        mov       rdi, %d                 ; exit code
        syscall                           ; invoke operating system to exit
    `, EXIT, code)
}

func syscall_write(fd int, addr string, size int) string {
    /* fd 1 -> stdout
    *
    */
    return fmt.Sprintf(`
        mov       rax, %#x           ; system call for write
        mov       rdi, %d            ; file handle 1 is stdout
        mov       rsi, %s            ; address of string to output
        mov       rdx, %d            ; number of bytes
        syscall                      ; invoke operating system to do the write
    `, WRITE, fd, addr, size)
}

func add_data(id, value string, ) string {
    return fmt.Sprintf("%s: db \"%s\"\n", id, value)
}

func compile(path string) {
	cmd := exec.Command("nasm", "-f", "macho64", path+".asm")
    err := cmd.Run()
    check(err)
    cmd = exec.Command("ld","-e", "start", "-static", path+".o", "-o", path)
    err = cmd.Run()
    check(err)
}

func execute(path string) {
    cmd := exec.Command(path)
    var out bytes.Buffer
    var errs bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errs
    err := cmd.Run()
    check(err)
    fmt.Println("stdout:", out.String())
    fmt.Println("stderr:", errs.String())
}

func generate(w *bufio.Writer) {
    w.WriteString(`
    global    start
    section   .text

start:`)

    message := "Hello Wonderful World"
    w.WriteString(syscall_write(1, "d0", len(message)))
    w.WriteString(syscall_exit(0))
    w.WriteString("\nsection   .data\n")
    w.WriteString(add_data("d0", message))
    w.Flush()
}

func main() {
    path := "build/main"
    f, err := os.Create(path+".asm")
    check(err)
    defer f.Close()

    w := bufio.NewWriter(f)
    generate(w)
    compile(path)
    execute(path)
}
