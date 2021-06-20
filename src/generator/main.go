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

func main() {
    f, err := os.Create("build/main.asm")
    check(err)
    defer f.Close()

    w := bufio.NewWriter(f)
    w.WriteString(`
    global    start
    section   .text

start:`)

    // write to stdout
    message := "Hello NoGo"
    w.WriteString(fmt.Sprintf(`
    mov       rax, 0x02000004         ; system call for write
    mov       rdi, 1                  ; file handle 1 is stdout
    mov       rsi, message            ; address of string to output
    mov       rdx, %d                 ; number of bytes
    syscall                           ; invoke operating system to do the write
    `,
    len(message)))

    // exit
    exit_code := 0
    w.WriteString(fmt.Sprintf(`
    mov       rax, 0x02000001         ; system call for exit
    mov       rdi, %d                 ; exit code
    syscall                           ; invoke operating system to exit
    `, exit_code))
    w.WriteString(fmt.Sprintf(`
    section   .data
message:  db        "%s", 15      ; note the newline at the end
`, message))
    w.Flush()


	cmd := exec.Command("nasm", "-f", "macho64", "build/main.asm")
    err = cmd.Run()
    check(err)

    cmd = exec.Command("ld","-e", "start", "-static", "build/main.o", "-o", "build/main")
    err = cmd.Run()
    check(err)

    cmd = exec.Command("./build/main")
    var out bytes.Buffer
    var errs bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errs
    err = cmd.Run()
    check(err)
    fmt.Println("stdout:", out.String())
    fmt.Println("stderr:", errs.String())
}
