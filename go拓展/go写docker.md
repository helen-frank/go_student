# 1. 基础技术

## 1.1 Linux Namespace

Docker 是一个使用了 `Linux Namespace` 和 `Cgroups` 的虚拟化工具

### 1.1.1 概念

`Linux Namespace`是 Kernel 的一个功能，它可以隔离一系列的系统资源，比如`PID`（Process ID）, `User ID` , `Network`等等，还可以在一些资源上，将进程隔离起来，这些资源包括进程树，网络接口，挂载点等。

chroot可以把当前目录变成根目录（被隔离开的）

![image-20210819173210323](go%E5%86%99docker.assets/image-20210819173210323.png)

当前 Linux 一共实现了 6 种不同类型的 Namespace 。

| Namespace类型     | 系统调用参数  | 内核版本 |
| ----------------- | ------------- | -------- |
| Mount Namespace   | CLONE_NEWNS   | 2.4.19   |
| UTS Namespace     | CLONE_NEWUTS  | 2.6.19   |
| IPC Namespace     | CLONE_NEWIPC  | 2.6.19   |
| PID Namespace     | CLONE_NEWPID  | 2.6.24   |
| Network Namespace | CLONE_NEWNET  | 2.6.29   |
| User Namespace    | CLONE_NEWUSER | 3.8      |

`Namespace`的API主要使用如下3个系统调用

- `clone()`创建新进程。根据系统调用参数来判断哪些类型的Namespace被创建，而且它们的子进程也会被包含到这些Namespace中
- `unshare()`将进程移出某个 Namespace
- `setns()`将进程加入到Namespace中



### 1.1.2 UTS Namespace

`UTS Namespace`主要用来隔离`nodename`和`domainname`两个系统标识。在`UTS Namespace`里，每个`Namespace`允许有自己的`hostname`

```go
/*
	@file	main.go
	@author	helenfrank(helenfrank@protonmail.com)
	@date	2021-08-19 17:42:00
*/

package main

import (
	"log"
	"os"
	"os/exec"
	"syscall"
)

func main() {
	cmd := exec.Command("sh")   // 用来指定被fork出来的新进程内的初始命令，默认使用sh来执行。
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS,
	}
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
}

```

`pstree -pl`查看系统中线程之间的关系

`echo $$` 输出当前的 PID

`readlink /proc/$(当前PID)/ns/uts` 

`readlink /proc/$(父PID)/ns/uts` 

验证父进程和子进程是否不在同一个`UTS Namespace`中



`UTS Namespace`对`hostname`做了隔离，在这里修改`hostname`不影响外部主机



### 1.1.3 IPC Namespace





