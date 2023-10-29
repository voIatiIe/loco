package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"
)

var clientProtoFiles = map[string]string{
    "steammessages_base.proto":   "base.pb.go",
    "steammessages_auth.steamclient.proto":  "auth.pb.go",
    "enums.proto": "enums.pb.go",
    "steammessages_unified_base.steamclient.proto": "unified_base.pb.go",
    "steammessages_clientserver_login.proto": "clientserver_login.pb.go",
}

var dotaProtoFiles = map[string]string{
	"base_gcmessages.proto":   "base.pb.go",
	"econ_shared_enums.proto": "econ_shared_enum.pb.go",
	"econ_gcmessages.proto":   "econ.pb.go",
	"gcsdk_gcmessages.proto":  "gcsdk.pb.go",
	"gcsystemmsgs.proto":      "system.pb.go",
	"steammessages.proto":     "steam.pb.go",
}


func main() {
    if len(os.Args) != 2 {
        fmt.Println("Usage: generator <proto/steamlang/clean>")
        os.Exit(1)
    }

    switch os.Args[1] {
        case "proto":
            buildProtobufs()
        case "steamlang":
            fmt.Println("Not implemented yet")
        case "clean":
            fmt.Println("Not implemented yet")
        default:
            fmt.Println("Usage: generator <proto/steamlang/clean>")
    }
}


func buildProtobufs() {
    fmt.Println("# Compiling Protobufs...")

    buildProtobuf("generator/Protobufs/steam", clientProtoFiles, "protocol/steam/")
    buildProtobuf("generator/Protobufs/dota2", dotaProtoFiles, "protocol/dota2/")
}


func buildProtobuf(srcPath string, protoFilesMap map[string]string, dstPath string) {
    err := os.MkdirAll(dstPath, os.ModePerm)
    if err != nil {
        fmt.Println("Error creating directory:", err)
    }

    opt := []string{"-I=" + srcPath, "--go_out=" + dstPath}

    protos, _ := os.ReadDir(srcPath)
    for _, protoFile := range protos {
        opt = append(opt, "--go_opt=M"+protoFile.Name()+"="+"protobuf/")
    }

    wg := &sync.WaitGroup{}
    for protoFile, goFile := range protoFilesMap {
        wg.Add(1)
        go compileProto(protoFile, goFile, dstPath, opt, wg)
    }
    wg.Wait()
}


func compileProto(
    protoFile string,
    goFile string,
    dstPath string,
    opt []string,
    wg *sync.WaitGroup,
) {
    defer wg.Done()

    args := make([]string, len(opt))
    copy(args, opt)

    args = append(args, protoFile)

    execute("protoc", args)
    rename(protoFile, goFile, dstPath)
}


func execute(command string, args []string) {
	cmd := exec.Command(command, args...)
    cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
        fmt.Printf("Protoc error: %s\n", err)
		os.Exit(1)
	}
}


func rename(protoFile string, goFile string, dstPath string) {
    oldPath := dstPath + "protobuf/" + strings.Replace(protoFile, ".proto", ".pb.go", 1)
    newPath := dstPath + "protobuf/" + goFile

    err := os.Rename(oldPath, newPath)
    if err != nil {
        fmt.Println("Error renaming file:", err)
    }
}
