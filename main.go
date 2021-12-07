package main

import (
	"errors"
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	log "github.com/sirupsen/logrus"
	ssh0 "golang.org/x/crypto/ssh"
	"os"
	"path/filepath"
)

// https://github.com/src-d/go-git/issues/637

// ssh-keygen -t ed25519 -C  "qiuzhanghua@icloud.com"  -P "" -f gitee

func main() {
	url := "ssh://gitea@git.taiji.com.cn:8888/TDP/tdp_go.git"
	//url = "git@github.com:qiuzhanghua/tdp_go.git"
	// 	url = "git@gitee.com:qiuzhanghua/Aug2021.git"  // gitee 不行
	//url = "https://github.com/qiuzhanghua/rust-docker.git"
	password := ""
	privateKeyFile, _ := ExpandHome("~/.ssh/id_rsa")
	_, err := os.Stat(privateKeyFile)
	if err != nil {
		log.Warnf("read file %s failed %s\n", privateKeyFile, err.Error())
		return
	}

	// Clone the given repository to the given directory
	log.Infof("git clone %s ", url)
	publicKeys, err := ssh.NewPublicKeysFromFile("gitea", privateKeyFile, password)
	publicKeys.SetHostKeyCallback(&ssh0.ClientConfig{
		HostKeyCallback: ssh0.InsecureIgnoreHostKey(),
	})
	fmt.Println(publicKeys)
	_, err = git.PlainClone("/tmp/foo", false, &git.CloneOptions{
		URL:      url,
		Auth:     publicKeys,
		Progress: os.Stdout,
	})

	fmt.Println(err)
}

// ExpandHome expend dir start with ~
func ExpandHome(path string) (string, error) {
	if len(path) == 0 {
		return path, nil
	}

	if path[0] != '~' {
		return path, nil
	}

	if len(path) > 1 && path[1] != '/' && path[1] != '\\' {
		return "", errors.New("cannot expand home dir")
	}

	dir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(dir, path[1:]), nil
}

func AbsPath(p string) (string, error) {
	absPath, err := filepath.EvalSymlinks(p)
	if err != nil {
		return ".", err
	}
	absPath, err = filepath.Abs(absPath)
	if err != nil {
		return ".", err
	}
	//absPath, err = filepath.Abs(path.Dir(absPath) + string(os.PathSeparator) + path.Base(absPath))
	//if err != nil {
	//	return ".", err
	//}
	return absPath, nil
}
