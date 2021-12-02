#  1. github/urfave/cli

## 1.1 install

```bash
go get github/urfave/cli
```

指定`v2`版本

```bash
go get github/urfave/cli.v2
```

## 1.2 基本用法

- 下面的代码没有添加 command，只有一个默认的动作。
- 所有命令后跟着的参数都会被放到 `c.Args` 中

```go
package main

import (
	"fmt"
	"log"
	"os"

	"github/urfave/cli.v1"
)

func main() {
	app := cli.NewApp()
	app.Name = "watch"
	app.Usage = "fight the loneliness!"
	app.Action = func(c *cli.Context) error {
		fmt.Printf("Hello %q!", c.Args().Get(0))
		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatalln(err)
	}
}
```

## 1.3 添加全局flag

- 下面的代码中添加了全局的__flag__: `--lang`
- __flag__的`Value`既可以通过`Destination`属性存储到自定义变量中，也可以使用`c.String("lang")`的方式来获取

```go
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	var language string

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "lang",
			Value:       "english",
			Usage:       "language for the greeting",
			Destination: &language,
		},
	}

	app.Action = func(c *cli.Context) error {
		name := "bwangel"
		if c.NArg() > 0 {
			name = c.Args().Get(0)
		}
		// if c.String("lang") == "chinese" {
		if language == "chinese" {
			fmt.Println("你好", name)
		} else {
			fmt.Println("Hello", name)
		}
		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatalln(err)
	}
}

```

