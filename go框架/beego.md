# 1.基本入门

## 1.1 安装

```bash
go get github.com/astaxie/beego
```

## 1.2 bee工具的使用

```bash
go get github.com/beego/bee
```

安装完之后，`bee` 可执行文件默认存放在 `$GOPATH/bin` 里面，所以需要把 `$GOPATH/bin` 添加到您的环境变量中，才可以进行下一步

如果本机设置了 `GOBIN`，那么上面的命令就会安装到 `GOBIN` 下，请添加 `GOBIN `到环境变量中

```bash
[helen@hf ~]$ bee
2021/06/27 20:32:08 INFO     ▶ 0001 Getting bee latest version...
2021/06/27 20:32:10 WARN     ▶ 0002 Update available 1.12.0 ==> 2.0.2
2021/06/27 20:32:10 WARN     ▶ 0003 Run `bee update` to update
2021/06/27 20:32:10 INFO     ▶ 0004 Your bee are up to date
Bee is a Fast and Flexible tool for managing your Beego Web Application.

USAGE
    bee command [arguments]

AVAILABLE COMMANDS

    version     Prints the current Bee version
    migrate     Runs database migrations
    api         Creates a Beego API application
    bale        Transforms non-Go files to Go source files
    fix         Fixes your application by making it compatible with newer versions of Beego
    pro         Source code generator
    dlv         Start a debugging session using Delve
    dockerize   Generates a Dockerfile for your Beego application
    generate    Source code generator
    hprose      Creates an RPC application based on Hprose and Beego frameworks
    new         Creates a Beego application
    pack        Compresses a Beego application into a single file
    rs          Run customized scripts
    run         Run the application by starting a local development server
    server      serving static content over HTTP on port
    update      Update Bee

Use bee help [command] for more information about a command.

ADDITIONAL HELP TOPICS


Use bee help [topic] for more information about that topic.

```

### 1.2.1 new命令

`new` 命令是新建一个 Web 项目，我们在命令行下执行 `bee new <项目名>` 就可以创建一个新的项目。但是注意该命令必须在 `$GOPATH/src` 下执行。最后会在 `$GOPATH/src` 相应目录下生成如下目录结构的项目：

```
myproject
├── conf
│   └── app.conf
├── controllers
│   └── default.go
├── main.go
├── models
├── routers
│   └── router.go
├── static
│   ├── css
│   ├── img
│   └── js
├── tests
│   └── default_test.go
└── views
    └── index.tpl
```

### 1.2.2 api命令

上面的 `new` 命令是用来新建 Web 项目，不过很多用户使用 beego 来开发 API 应用。所以这个 `api` 命令就是用来创建 API 应用的，执行命令之后如下所示：

```bash
apiproject
├── conf
│   └── app.conf
├── controllers
│   └── object.go
│   └── user.go
├── docs
│   └── doc.go
├── main.go
├── models
│   └── object.go
│   └── user.go
├── routers
│   └── router.go
└── tests
    └── default_test.go
```

可以看到和 Web 项目相比，少了 `static` 和 `views` 目录，多了一个 test 模块，用来做单元测试的。

同时，该命令还支持一些自定义参数自动连接数据库创建相关 model 和 controller: `bee api [appname] [-tables=""] [-driver=mysql] [-conn="root:<password>@tcp(127.0.0.1:3306)/test"]` 如果 conn 参数为空则创建一个示例项目，否则将基于链接信息链接数据库创建项目。

### 1.2.3 run命令

开发 Go 项目的时候最大的问题是经常需要自己手动去编译再运行，`bee run` 命令是监控 beego 的项目，通过 [fsnotify](https://github.com/howeyc/fsnotify)监控文件系统。但是注意该命令必须在 `$GOPATH/src/appname` 下执行。 这样我们在开发过程中就可以实时的看到项目修改之后的效果

### 1.2.4 pack命令

`pack` 目录用来发布应用的时候打包，会把项目打包成 zip 包，这样我们部署的时候直接把打包之后的项目上传，解压就可以部署了

### 1.2.5 bale命令

这个命令目前仅限内部使用，具体实现方案未完善，主要用来压缩所有的静态文件变成一个变量申明文件，全部编译到二进制文件里面，用户发布的时候携带静态文件，包括 js、css、img 和 views。最后在启动运行时进行非覆盖式的自解压。

### 1.2.6 version命令

这个命令是动态获取 bee、beego 和 Go 的版本，这样一旦用户出现错误，可以通过该命令来查看当前的版本

### 1.2.7 generate命令

这个命令是用来自动化的生成代码的，包含了从数据库一键生成 model，还包含了 scaffold 的，通过这个命令，让大家开发代码不再慢

```bash
bee generate scaffold [scaffoldname] [-fields=""] [-driver=mysql] [-conn="root:@tcp(127.0.0.1:3306)/test"]
    The generate scaffold command will do a number of things for you.
    -fields: a list of table fields. Format: field:type, ...
    -driver: [mysql | postgres | sqlite], the default is mysql
    -conn:   the connection string used by the driver, the default is root:@tcp(127.0.0.1:3306)/test
    example: bee generate scaffold post -fields="title:string,body:text"

bee generate model [modelname] [-fields=""]
    generate RESTful model based on fields
    -fields: a list of table fields. Format: field:type, ...

bee generate controller [controllerfile]
    generate RESTful controllers

bee generate view [viewpath]
    generate CRUD view in viewpath

bee generate migration [migrationfile] [-fields=""]
    generate migration file for making database schema update
    -fields: a list of table fields. Format: field:type, ...

bee generate docs
    generate swagger doc file

bee generate test [routerfile]
    generate testcase

bee generate appcode [-tables=""] [-driver=mysql] [-conn="root:@tcp(127.0.0.1:3306)/test"] [-level=3]
    generate appcode based on an existing database
    -tables: a list of table names separated by ',', default is empty, indicating all tables
    -driver: [mysql | postgres | sqlite], the default is mysql
    -conn:   the connection string used by the driver.
             default for mysql:    root:@tcp(127.0.0.1:3306)/test
             default for postgres: postgres://postgres:postgres@127.0.0.1:5432/postgres
    -level:  [1 | 2 | 3], 1 = models; 2 = models,controllers; 3 = models,controllers,router
```

### 1.2.8 migrate命令

这个命令是应用的数据库迁移命令，主要是用来每次应用升级，降级的SQL管理。

```bash
bee migrate [-driver=mysql] [-conn="root:@tcp(127.0.0.1:3306)/test"]
    run all outstanding migrations
    -driver: [mysql | postgresql | sqlite], the default is mysql
    -conn:   the connection string used by the driver, the default is root:@tcp(127.0.0.1:3306)/test

bee migrate rollback [-driver=mysql] [-conn="root:@tcp(127.0.0.1:3306)/test"]
    rollback the last migration operation
    -driver: [mysql | postgresql | sqlite], the default is mysql
    -conn:   the connection string used by the driver, the default is root:@tcp(127.0.0.1:3306)/test

bee migrate reset [-driver=mysql] [-conn="root:@tcp(127.0.0.1:3306)/test"]
    rollback all migrations
    -driver: [mysql | postgresql | sqlite], the default is mysql
    -conn:   the connection string used by the driver, the default is root:@tcp(127.0.0.1:3306)/test

bee migrate refresh [-driver=mysql] [-conn="root:@tcp(127.0.0.1:3306)/test"]
    rollback all migrations and run them all again
    -driver: [mysql | postgresql | sqlite], the default is mysql
    -conn:   the connection string used by the driver, the default is root:@tcp(127.0.0.1:3306)/test
```

### 1.2.9 dockerize命令

这个命令可以通过生成Dockerfile文件来实现docker化你的应用。

例子:
生成一个以1.6.4版本Go环境为基础镜像的Dockerfile,并暴露9000端口

```bash
$ bee dockerize -image="library/golang:1.6.4" -expose=9000
______
| ___ \
| |_/ /  ___   ___
| ___ \ / _ \ / _ \
| |_/ /|  __/|  __/
\____/  \___| \___| v1.6.2
2016/12/26 22:34:54 INFO     ▶ 0001 Generating Dockerfile...
2016/12/26 22:34:54 SUCCESS  ▶ 0002 Dockerfile generated.
```



## 1.3 bee工具配置文件

在 bee 工具的源码目录下有一个 `bee.json` 文件，这个文件是针对 bee 工具的一些行为进行配置。该功能还未完全开发完成，不过其中的一些选项已经可以使用：

- `"version": 0`：配置文件版本，用于对比是否发生不兼容的配置格式版本。
- `"go_install": false`：如果您的包均使用完整的导入路径（例如：`github.com/user/repo/subpkg`）,则可以启用该选项来进行 `go install` 操作，加快构建操作。
- `"watch_ext": []`：用于监控其它类型的文件（默认只监控后缀为 `.go` 的文件）。
- `"dir_structure":{}`：如果您的目录名与默认的 MVC 架构的不同，则可以使用该选项进行修改。
- `"cmd_args": []`：如果您需要在每次启动时加入启动参数，则可以使用该选项。
- `"envs": []`：如果您需要在每次启动时设置临时环境变量参数，则可以使用该选项。

## 1.4 新建项目

### 1.4.1 创建项目

```bash
bee new quickstart
--------------------------------------------

quickstart
|-- conf
|   `-- app.conf
|-- controllers
|   `-- default.go
|-- main.go
|-- models
|-- routers
|   `-- router.go
|-- static
|   |-- css
|   |-- img
|   `-- js
|-- tests
|   `-- default_test.go
`-- views
    `-- index.tpl
```

### 1.4.2 运行项目

使用 `bee run` 来运行该项目，这样就可以做到热编译的效果



## 1.5 路由设置

```go
package main

import (
    _ "quickstart/routers"
    "github.com/astaxie/beego"
)

func main() {
    beego.Run()
}
```

引入了一个包 `_ "quickstart/routers"`,这个包只引入执行了里面的 init 函数

```go
package routers

import (
    "quickstart/controllers"
    "github.com/astaxie/beego"
)

func init() {
    beego.Router("/", &controllers.MainController{})
}
```

路由包里面我们看到执行了路由注册 `beego.Router`, 这个函数的功能是映射 URL 到 controller，第一个参数是 URL (用户请求的地址)，这里我们注册的是 `/`，也就是我们访问的不带任何参数的 URL，第二个参数是对应的 Controller，也就是我们即将把请求分发到那个控制器来执行相应的逻辑，我们可以执行类似的方式注册如下路由：

```go
beego.Router("/user", &controllers.UserController{})
```

这样用户就可以通过访问 `/user` 去执行 `UserController` 的逻辑。这就是我们所谓的路由，更多更复杂的路由规则请查询beego 的路由设置

再回来看看 main 函数里面的 `beego.Run`， `beego.Run` 执行之后，我们看到的效果好像只是监听服务端口这个过程，但是它内部做了很多事情：

- 解析配置文件

  beego 会自动解析在 conf 目录下面的配置文件 `app.conf`，通过修改配置文件相关的属性，我们可以定义：开启的端口，是否开启 session，应用名称等信息。

- 执行用户的 hookfunc

  beego 会执行用户注册的 hookfunc，默认的已经存在了注册 mime，用户可以通过函数 `AddAPPStartHook` 注册自己的启动函数。

- 是否开启 session

  会根据上面配置文件的分析之后判断是否开启 session，如果开启的话就初始化全局的 session。

- 是否编译模板

  beego 会在启动的时候根据配置把 views 目录下的所有模板进行预编译，然后存在 map 里面，这样可以有效的提高模板运行的效率，无需进行多次编译。

- 是否开启文档功能

  根据 EnableDocs 配置判断是否开启内置的文档路由功能

- 是否启动管理模块

  beego 目前做了一个很酷的模块，应用内监控模块，会在 8088 端口做一个内部监听，我们可以通过这个端口查询到 QPS、CPU、内存、GC、goroutine、thread 等统计信息。

- 监听服务端口

  这是最后一步也就是我们看到的访问 8080 看到的网页端口，内部其实调用了 `ListenAndServe`，充分利用了 goroutine 的优势

一旦 run 起来之后，我们的服务就监听在两个端口了，一个服务端口 8080 作为对外服务，另一个 8088 端口实行对内监控。

通过这个代码的分析我们了解了 beego 运行起来的过程，以及内部的一些机制。接下来让我们去剥离 Controller 如何来处理逻辑的。

## 1.6 Contronller 运行机制

`controllers/default.go`

```go
package controllers

import (
        "github.com/astaxie/beego"
)

type MainController struct {
        beego.Controller
}

func (this *MainController) Get() {
        this.Data["Website"] = "beego.me"
        this.Data["Email"] = "astaxie@gmail.com"
        this.TplName = "index.tpl"
}
```

上面的代码显示首先我们声明了一个控制器 `MainController`，这个控制器里面内嵌了 `beego.Controller`，这就是 Go 的嵌入方式，也就是 `MainController` 自动拥有了所有 `beego.Controller` 的方法。

而 `beego.Controller` 拥有很多方法，其中包括 `Init`、`Prepare`、`Post`、`Get`、`Delete`、`Head` 等方法。我们可以通过重写的方式来实现这些方法，而我们上面的代码就是重写了 `Get` 方法。

我们先前介绍过 beego 是一个 RESTful 的框架，所以我们的请求默认是执行对应 `req.Method` 的方法。例如浏览器的是 `GET` 请求，那么默认就会执行 `MainController` 下的 `Get` 方法。这样我们上面的 Get 方法就会被执行到，这样就进入了我们的逻辑处理。（用户可以改变这个行为，通过注册自定义的函数名

里面的代码是需要执行的逻辑，这里只是简单的输出数据，我们可以通过各种方式获取数据，然后赋值到 `this.Data` 中，这是一个用来存储输出数据的 map，可以赋值任意类型的值，这里我们只是简单举例输出两个字符串。

最后一个就是需要去渲染的模板，`this.TplName` 就是需要渲染的模板，这里指定了 `index.tpl`，如果用户不设置该参数，那么默认会去到模板目录的 `Controller/<方法名>.tpl` 查找，例如上面的方法会去 `maincontroller/get.tpl` ***(文件、文件夹必须小写)\***。

用户设置了模板之后系统会自动的调用 `Render` 函数（这个函数是在 beego.Controller 中实现的），所以无需用户自己来调用渲染。

当然也可以不使用模版，直接用 `this.Ctx.WriteString` 输出字符串，如：

```
func (this *MainController) Get() {
        this.Ctx.WriteString("hello")
}
```



## 1.7 Model逻辑

model 层一般用来做这些操作，我们的 `bee new` 例子不存在 Model 的演示，但是 `bee api` 应用中存在 model 的应用。说的简单一点，如果您的应用足够简单，那么 Controller 可以处理一切的逻辑，如果您的逻辑里面存在着可以复用的东西，那么就抽取出来变成一个模块。因此 Model 就是逐步抽象的过程，一般我们会在 Model 里面处理一些数据读取，如下是一个日志分析应用中的代码片段：

```go
package models

import (
    "loggo/utils"
    "path/filepath"
    "strconv"
    "strings"
)

var (
    NotPV []string = []string{"css", "js", "class", "gif", "jpg", "jpeg", "png", "bmp", "ico", "rss", "xml", "swf"}
)

const big = 0xFFFFFF

func LogPV(urls string) bool {
    ext := filepath.Ext(urls)
    if ext == "" {
        return true
    }
    for _, v := range NotPV {
        if v == strings.ToLower(ext) {
            return false
        }
    }
    return true
}
```

所以如果您的应用足够简单，那么就不需要 Model 了；如果你的模块开始多了，需要复用，需要逻辑分离了，那么 Model 是必不可少的。接下来我们将分析如何编写 View 层的东西。



## 1.8 view编写

在前面编写 Controller 的时候，我们在 Get 里面写过这样的语句 `this.TplName = "index.tpl"`，设置显示的模板文件，默认支持 `tpl` 和 `html` 的后缀名，如果想设置其他后缀你可以调用 `beego.AddTemplateExt` 接口设置，那么模板如何来显示相应的数据呢？beego 采用了 Go 语言默认的模板引擎，所以和 Go 的模板语法一样，Go 模板的详细使用方法请参考《Go Web 编程》模板使用指南

我们看看快速入门里面的代码（去掉了 css 样式）：

```html
<!DOCTYPE html>

<html>
    <head>
        <title>Beego</title>
        <meta http-equiv="Content-Type" content="text/html; charset=utf-8">
    </head>
    <body>
        <header class="hero-unit" style="background-color:#A9F16C">
            <div class="container">
                <div class="row">
                    <div class="hero-text">
                        <h1>Welcome to Beego!</h1>
                        <p class="description">
                            Beego is a simple & powerful Go web framework which is inspired by tornado and sinatra.
                            <br />
                            Official website: <a href="http://{{.Website}}">{{.Website}}</a>
                            <br />
                            Contact me: {{.Email}}
                        </p>
                    </div>
                </div>
            </div>
        </header>
    </body>
</html>
```

我们在 Controller 里面把数据赋值给了 data（map 类型），然后我们在模板中就直接通过 key 访问 `.Website` 和 `.Email` 。这样就做到了数据的输出。接下来我们讲解如何让静态文件输出。



## 1.9 静态文件处理

beego 默认注册了 static 目录为静态处理的目录，注册样式：URL 前缀和映射的目录（在/main.go文件中beego.Run()之前加入）：

```go
StaticDir["/static"] = "static"
```

用户可以设置多个静态文件处理目录，例如你有多个文件下载目录 download1、download2，你可以这样映射（在 /main.go 文件中 beego.Run() 之前加入）：

```go
beego.SetStaticPath("/down1", "download1")
beego.SetStaticPath("/down2", "download2")
```

这样用户访问 URL `http://localhost:8080/down1/123.txt` 则会请求 download1 目录下的 123.txt 文件。



# 2.beego的MVC架构介绍

## 2.1 controller设计

### 2.1.1 参数配置

beego 目前支持 INI、XML、JSON、YAML 格式的配置文件解析，但是默认采用了 INI 格式解析，用户可以通过简单的配置就可以获得很大的灵活性。

#### 2.1.1.1 默认配置解析

beego 默认会解析当前应用下的 `conf/app.conf` 文件。

通过这个文件你可以初始化很多 beego 的默认参数：

```
appname = beepkg
httpaddr = "127.0.0.1"
httpport = 9090
runmode ="dev"
autorender = false
recoverpanic = false
viewspath = "myview"
```

上面这些参数会替换 beego 默认的一些参数, beego 的参数主要有哪些呢？请参考https://godoc.org/github.com/astaxie/beego#pkg-constants 。 BConfig 就是 beego 里面的默认的配置，你也可以直接通过`beego.BConfig.AppName="beepkg"`这样来修改，和上面的配置效果一样，只是一个在代码里面写死了， 而配置文件就会显得更加灵活。

你也可以在配置文件中配置应用需要用的一些配置信息，例如下面所示的数据库信息：

```
mysqluser = "root"
mysqlpass = "rootpass"
mysqlurls = "127.0.0.1"
mysqldb   = "beego"
```

那么你就可以通过如下的方式获取设置的配置信息:

```
beego.AppConfig.String("mysqluser")
beego.AppConfig.String("mysqlpass")
beego.AppConfig.String("mysqlurls")
beego.AppConfig.String("mysqldb")
```

AppConfig 的方法如下：

- Set(key, val string) error
- String(key string) string
- Strings(key string) []string
- Int(key string) (int, error)
- Int64(key string) (int64, error)
- Bool(key string) (bool, error)
- Float(key string) (float64, error)
- DefaultString(key string, defaultVal string) string
- DefaultStrings(key string, defaultVal []string)
- DefaultInt(key string, defaultVal int) int
- DefaultInt64(key string, defaultVal int64) int64
- DefaultBool(key string, defaultVal bool) bool
- DefaultFloat(key string, defaultVal float64) float64
- DIY(key string) (interface{}, error)
- GetSection(section string) (map[string]string, error)
- SaveConfigFile(filename string) error

在使用 ini 类型的配置文件中, key 支持 section::key 模式.

你可以用 Default* 方法返回默认值.

##### 2.1.1.1.1 不同级别的配置

在配置文件里面支持 section，可以有不同的 Runmode 的配置，默认优先读取 runmode 下的配置信息，例如下面的配置文件：

```
appname = beepkg
httpaddr = "127.0.0.1"
httpport = 9090
runmode ="dev"
autorender = false
recoverpanic = false
viewspath = "myview"

[dev]
httpport = 8080
[prod]
httpport = 8088
[test]
httpport = 8888
```

上面的配置文件就是在不同的 runmode 下解析不同的配置，例如在 dev 模式下，httpport 是 8080，在 prod 模式下是 8088，在 test 模式下是 8888。其他配置文件同理。解析的时候优先解析 runmode 下的配置，然后解析默认的配置。

读取不同模式下配置参数的方法是“模式::配置参数名”，比如：beego.AppConfig.String("dev::mysqluser")。

对于自定义的参数，需使用 beego.GetConfig(typ, key string, defaultVal interface{}) 来获取指定 runmode 下的配置（需 1.4.0 以上版本），typ 为参数类型，key 为参数名, defaultVal 为默认值。



##### 2.1.1.1.2 多个配置文件

INI 格式配置支持 `include` 方式，引用多个配置文件，例如下面的两个配置文件效果同上：

`app.conf`

```
appname = beepkg
httpaddr = "127.0.0.1"
httpport = 9090

include "app2.conf"
```

`app2.conf`

```
runmode ="dev"
autorender = false
recoverpanic = false
viewspath = "myview"

[dev]
httpport = 8080
[prod]
httpport = 8088
[test]
httpport = 8888
```



##### 2.1.1.1.3. 支持环境变量配置

配置文件解析支持从环境变量中获取配置项，配置项格式：`${环境变量}`。例如下面的配置中优先使用环境变量中配置的 runmode 和 httpport，如果有配置环境变量 ProRunMode 则优先使用该环境变量值。如果不存在或者为空，则使用 "dev" 作为 runmode。

app.conf

```
runmode  = "${ProRunMode||dev}"
httpport = "${ProPort||9090}"
```



##### 2.1.1.1.4 系统默认参数

beego 中带有很多可配置的参数，我们来一一认识一下它们，这样有利于我们在接下来的 beego 开发中可以充分的发挥他们的作用(你可以通过在 `conf/app.conf` 中设置对应的值，不区分大小写)：

1. 基础配置

   - BConfig 保存了所有 beego 里面的系统默认参数，你可以通过 `beego.BConfig` 来访问和修改底下的所有配置信息.

   > > 配置文件路径，默认是应用程序对应的目录下的 `conf/app.conf`，用户可以在程序代码中加载自己的配置文件 `beego.LoadAppConfig("ini", "conf/app2.conf")` 也可以加载多个文件，只要你调用多次就可以了，如果后面的文件和前面的 key 冲突，那么以最新加载的为最新值

2. App配置

   - AppName

     应用名称，默认是 beego。通过 `bee new` 创建的是创建的项目名。

     `beego.BConfig.AppName = "beego"`

   - RunMode

     应用的运行模式，可选值为 `prod`, `dev` 或者 `test`. 默认是 `dev`, 为开发模式，在开发模式下出错会提示友好的出错页面，如前面错误描述中所述。

     `beego.BConfig.RunMode = "dev"`

   - RouterCaseSensitive

     是否路由忽略大小写匹配，默认是 true，区分大小写

     `beego.BConfig.RouterCaseSensitive = true`

   - ServerName

     beego 服务器默认在请求的时候输出 server 为 beego。

     `beego.BConfig.ServerName = "beego"`

   - RecoverPanic

     是否异常恢复，默认值为 true，即当应用出现异常的情况，通过 recover 恢复回来，而不会导致应用异常退出。

     `beego.BConfig.RecoverPanic = true`

   - CopyRequestBody

     是否允许在 HTTP 请求时，返回原始请求体数据字节，默认为 false （GET or HEAD or 上传文件请求除外）。

     `beego.BConfig.CopyRequestBody = false`

   - EnableGzip

     是否开启 gzip 支持，默认为 false 不支持 gzip，一旦开启了 gzip，那么在模板输出的内容会进行 gzip 或者 zlib 压缩，根据用户的 Accept-Encoding 来判断。

     `beego.BConfig.EnableGzip = false`

     Gzip允许用户自定义压缩级别、压缩长度阈值和针对请求类型压缩:

     1. 压缩级别, `gzipCompressLevel = 9`,取值为 1~9,如果不设置为 1(最快压缩)
     2. 压缩长度阈值, `gzipMinLength = 256`,当原始内容长度大于此阈值时才开启压缩,默认为 20B(ngnix默认长度)
     3. 请求类型, `includedMethods = get;post`,针对哪些请求类型进行压缩,默认只针对 GET 请求压缩

   - MaxMemory

     文件上传默认内存缓存大小，默认值是 `1 << 26`(64M)。

     `beego.BConfig.MaxMemory = 1 << 26`

   - EnableErrorsShow

     是否显示系统错误信息，默认为 true。

     `beego.BConfig.EnableErrorsShow = true`

   - EnableErrorsRender

     是否将错误信息进行渲染，默认值为 true，即出错会提示友好的出错页面，对于 API 类型的应用可能需要将该选项设置为 false 以阻止在 `dev` 模式下不必要的模板渲染信息返回。

3. Web配置

   - AutoRender

     是否模板自动渲染，默认值为 true，对于 API 类型的应用，应用需要把该选项设置为 false，不需要渲染模板。

     `beego.BConfig.WebConfig.AutoRender = true`

   - EnableDocs

     是否开启文档内置功能，默认是 false

     `beego.BConfig.WebConfig.EnableDocs = true`

   - FlashName

     Flash 数据设置时 Cookie 的名称，默认是 BEEGO_FLASH

     `beego.BConfig.WebConfig.FlashName = "BEEGO_FLASH"`

   - FlashSeperator

     Flash 数据的分隔符，默认是 BEEGOFLASH

     `beego.BConfig.WebConfig.FlashSeparator = "BEEGOFLASH"`

   - DirectoryIndex

     是否开启静态目录的列表显示，默认不显示目录，返回 403 错误。

     `beego.BConfig.WebConfig.DirectoryIndex = false`

   - StaticDir

     静态文件目录设置，默认是static

     可配置单个或多个目录:

     1. 单个目录, `StaticDir = download`. 相当于 `beego.SetStaticPath("/download","download")`

     2. 多个目录, `StaticDir = download:down download2:down2`. 相当于 `beego.SetStaticPath("/download","down")` 和 `beego.SetStaticPath("/download2","down2")`

        `beego.BConfig.WebConfig.StaticDir`

   - StaticExtensionsToGzip

     允许哪些后缀名的静态文件进行 gzip 压缩，默认支持 .css 和 .js

     `beego.BConfig.WebConfig.StaticExtensionsToGzip = []string{".css", ".js"}`

     等价 config 文件中

     ```
       StaticExtensionsToGzip = .css, .js
     ```

   - TemplateLeft

     模板左标签，默认值是`{{`。

     `beego.BConfig.WebConfig.TemplateLeft="{{"`

   - TemplateRight

     模板右标签，默认值是`}}`。

     `beego.BConfig.WebConfig.TemplateRight="}}"`

   - ViewsPath

     模板路径，默认值是 views。

     `beego.BConfig.WebConfig.ViewsPath="views"`

   - EnableXSRF

     是否开启 XSRF，默认为 false，不开启。

     `beego.BConfig.WebConfig.EnableXSRF = false`

   - XSRFKEY

     XSRF 的 key 信息，默认值是 beegoxsrf。 EnableXSRF＝true 才有效

     `beego.BConfig.WebConfig.XSRFKEY = "beegoxsrf"`

   - XSRFExpire

     XSRF 过期时间，默认值是 0，不过期。

     `beego.BConfig.WebConfig.XSRFExpire = 0`

4. 监听配置

   - Graceful

     是否开启热升级，默认是 false，关闭热升级。

     `beego.BConfig.Listen.Graceful=false`

   - ServerTimeOut

     设置 HTTP 的超时时间，默认是 0，不超时。

     `beego.BConfig.Listen.ServerTimeOut=0`

   - ListenTCP4

     监听本地网络地址类型，默认是TCP6，可以通过设置为true设置为TCP4。

     `beego.BConfig.Listen.ListenTCP4 = true`

   - EnableHTTP

     是否启用 HTTP 监听，默认是 true。

     `beego.BConfig.Listen.EnableHTTP = true`

   - HTTPAddr

     应用监听地址，默认为空，监听所有的网卡 IP。

     `beego.BConfig.Listen.HTTPAddr = ""`

   - HTTPPort

     应用监听端口，默认为 8080。

     `beego.BConfig.Listen.HTTPPort = 8080`

   - EnableHTTPS

     是否启用 HTTPS，默认是 false 关闭。当需要启用时，先设置 EnableHTTPS = true，并设置 `HTTPSCertFile` 和 `HTTPSKeyFile`

     `beego.BConfig.Listen.EnableHTTPS = false`

   - HTTPSAddr

     应用监听地址，默认为空，监听所有的网卡 IP。

     `beego.BConfig.Listen.HTTPSAddr = ""`

   - HTTPSPort

     应用监听端口，默认为 10443

     `beego.BConfig.Listen.HTTPSPort = 10443`

   - HTTPSCertFile

     开启 HTTPS 后，ssl 证书路径，默认为空。

     `beego.BConfig.Listen.HTTPSCertFile = "conf/ssl.crt"`

   - HTTPSKeyFile

     开启 HTTPS 之后，SSL 证书 keyfile 的路径。

     `beego.BConfig.Listen.HTTPSKeyFile = "conf/ssl.key"`

   - EnableAdmin

     是否开启进程内监控模块，默认 false 关闭。

     `beego.BConfig.Listen.EnableAdmin = false`

   - AdminAddr

     监控程序监听的地址，默认值是 localhost 。

     `beego.BConfig.Listen.AdminAddr = "localhost"`

   - AdminPort

     监控程序监听的地址，默认值是 8088 。

     `beego.BConfig.Listen.AdminPort = 8088`

   - EnableFcgi

     是否启用 fastcgi ， 默认是 false。

     `beego.BConfig.Listen.EnableFcgi = false`

   - EnableStdIo

     通过fastcgi 标准I/O，启用 fastcgi 后才生效，默认 false。

     `beego.BConfig.Listen.EnableStdIo = false`

5. Session配置

   - SessionOn

     session 是否开启，默认是 false。

     `beego.BConfig.WebConfig.Session.SessionOn = false`

   - SessionProvider

     session 的引擎，默认是 memory，详细参见 `session 模块`。

     `beego.BConfig.WebConfig.Session.SessionProvider = ""`

   - SessionName

     存在客户端的 cookie 名称，默认值是 beegosessionID。

     `beego.BConfig.WebConfig.Session.SessionName = "beegosessionID"`

   - SessionGCMaxLifetime

     session 过期时间，默认值是 3600 秒。

     `beego.BConfig.WebConfig.Session.SessionGCMaxLifetime = 3600`

   - SessionProviderConfig

     配置信息，根据不同的引擎设置不同的配置信息，详细的配置请看下面的引擎设置，详细参见 session 模块

   - SessionCookieLifeTime

     session 默认存在客户端的 cookie 的时间，默认值是 3600 秒。

     `beego.BConfig.WebConfig.Session.SessionCookieLifeTime = 3600`

   - SessionAutoSetCookie

     是否开启SetCookie, 默认值 true 开启。

     `beego.BConfig.WebConfig.Session.SessionAutoSetCookie = true`

   - SessionDomain

     session cookie 存储域名, 默认空。

     `beego.BConfig.WebConfig.Session.SessionDomain = ""`

6. Log配置

   ```
   log详细配置，请参见 `logs 模块`。
   ```

   - AccessLogs

     是否输出日志到 Log，默认在 prod 模式下不会输出日志，默认为 false 不输出日志。此参数不支持配置文件配置。

     `beego.BConfig.Log.AccessLogs = false`

   - FileLineNum

     是否在日志里面显示文件名和输出日志行号，默认 true。此参数不支持配置文件配置。

     `beego.BConfig.Log.FileLineNum = true`

   - Outputs

     日志输出配置，参考 logs 模块，console file 等配置，此参数不支持配置文件配置。

     `beego.BConfig.Log.Outputs = map[string]string{"console": ""}`

     or

     `beego.BConfig.Log.Outputs["console"] = ""`

### 2.1.2 路由设置

什么是路由设置呢？前面介绍的 MVC 结构执行时，介绍过 beego 存在三种方式的路由:固定路由、正则路由、自动路由，接下来详细的讲解如何使用这三种路由。

#### 2.1.2.1 基础路由

从 beego 1.2 版本开始支持了基本的 RESTful 函数式路由,应用中的大多数路由都会定义在 `routers/router.go` 文件中。最简单的 beego 路由由 URI 和闭包函数组成。

1. GET

```go
beego.Get("/",func(ctx *context.Context){
     ctx.Output.Body([]byte("hello world"))
})
```

2. 基本POST路由

```
beego.Post("/alice",func(ctx *context.Context){
     ctx.Output.Body([]byte("bob"))
})
```

3. 注册一个可以响应任何HTTP的路由

```go
beego.Any("/foo",func(ctx *context.Context){
     ctx.Output.Body([]byte("bar"))
})
```

4. 所有的支持的基本函数

- beego.Get(router, beego.FilterFunc)
- beego.Post(router, beego.FilterFunc)
- beego.Put(router, beego.FilterFunc)
- beego.Patch(router, beego.FilterFunc)
- beego.Head(router, beego.FilterFunc)
- beego.Options(router, beego.FilterFunc)
- beego.Delete(router, beego.FilterFunc)
- beego.Any(router, beego.FilterFunc)

5. 支持自定义的handler实现

有些时候我们已经实现了一些 rpc 的应用,但是想要集成到 beego 中,或者其他的 httpserver 应用,集成到 beego 中来.现在可以很方便的集成:

```go
s := rpc.NewServer()
s.RegisterCodec(json.NewCodec(), "application/json")
s.RegisterService(new(HelloService), "")
beego.Handler("/rpc", s)
```

`beego.Handler(router, http.Handler)` 这个函数是关键,第一个参数表示路由 URI, 第二个就是你自己实现的 `http.Handler`, 注册之后就会把所有 rpc 作为前缀的请求分发到 `http.Handler` 中进行处理.

这个函数其实还有第三个参数就是是否是前缀匹配,默认是 false, 如果设置了 true, 那么就会在路由匹配的时候前缀匹配,即 `/rpc/user` 这样的也会匹配去运行

6. 路由参数

后面会讲到固定路由,正则路由,这些参数一样适用于上面的这些函数

#### 2.1.2.2 RESTful Controller 路由

在介绍这三种 beego 的路由实现之前先介绍 RESTful，我们知道 RESTful 是一种目前 API 开发中广泛采用的形式，beego 默认就是支持这样的请求方法，也就是用户 Get 请求就执行 Get 方法，Post 请求就执行 Post 方法。因此默认的路由是这样 RESTful 的请求方式。

#### 2.1.2.3 固定路由

固定路由也就是全匹配的路由，如下所示：

```go
beego.Router("/", &controllers.MainController{})
beego.Router("/admin", &admin.UserController{})
beego.Router("/admin/index", &admin.ArticleController{})
beego.Router("/admin/addpkg", &admin.AddController{})
```

如上所示的路由就是我们最常用的路由方式，一个固定的路由，一个控制器，然后根据用户请求方法不同请求控制器中对应的方法，典型的 RESTful 方式。

#### 2.1.2.4 正则路由

为了用户更加方便的路由设置，beego 参考了 sinatra 的路由实现，支持多种方式的路由：

- beego.Router("/api/?:id", &controllers.RController{})

  默认匹配 //例如对于URL"/api/123"可以匹配成功，此时变量":id"值为"123"

- beego.Router("/api/:id", &controllers.RController{})

  默认匹配 //例如对于URL"/api/123"可以匹配成功，此时变量":id"值为"123"，但URL"/api/"匹配失败

- beego.Router("/api/:id([0-9]+)", &controllers.RController{})

  自定义正则匹配 //例如对于URL"/api/123"可以匹配成功，此时变量":id"值为"123"

- beego.Router("/user/:username([\\w]+)", &controllers.RController{})

  正则字符串匹配 //例如对于URL"/user/astaxie"可以匹配成功，此时变量":username"值为"astaxie"

- beego.Router("/download/*.*", &controllers.RController{})

  *匹配方式 //例如对于URL"/download/file/api.xml"可以匹配成功，此时变量":path"值为"file/api"， ":ext"值为"xml"

- beego.Router("/download/ceshi/*", &controllers.RController{})

  *全匹配方式 //例如对于URL"/download/ceshi/file/api.json"可以匹配成功，此时变量":splat"值为"file/api.json"

- beego.Router("/:id:int", &controllers.RController{})

  int 类型设置方式，匹配 :id为int 类型，框架帮你实现了正则 ([0-9]+)

- beego.Router("/:hi:string", &controllers.RController{})

  string 类型设置方式，匹配 :hi 为 string 类型。框架帮你实现了正则 ([\w]+)

- beego.Router("/cms_:id([0-9]+).html", &controllers.CmsController{})

  带有前缀的自定义正则 //匹配 :id 为正则类型。匹配 cms_123.html 这样的 url :id = 123

可以在 Controller 中通过如下方式获取上面的变量：

```
this.Ctx.Input.Param(":id")
this.Ctx.Input.Param(":username")
this.Ctx.Input.Param(":splat")
this.Ctx.Input.Param(":path")
this.Ctx.Input.Param(":ext")
```



#### 2.1.2.5 自定义方法及RESTful规则

上面列举的是默认的请求方法名（请求的 method 和函数名一致，例如 `GET` 请求执行 `Get` 函数，`POST` 请求执行 `Post` 函数），如果用户期望自定义函数名，那么可以使用如下方式：

```
beego.Router("/",&IndexController{},"*:Index")
```

使用第三个参数，第三个参数就是用来设置对应 method 到函数名，定义如下

- `*`表示任意的 method 都执行该函数
- 使用 httpmethod:funcname 格式来展示
- 多个不同的格式使用 `;` 分割
- 多个 method 对应同一个 funcname，method 之间通过 `,` 来分割

以下是一个 RESTful 的设计示例：

```
beego.Router("/api/list",&RestController{},"*:ListFood")
beego.Router("/api/create",&RestController{},"post:CreateFood")
beego.Router("/api/update",&RestController{},"put:UpdateFood")
beego.Router("/api/delete",&RestController{},"delete:DeleteFood")
```

以下是多个 HTTP Method 指向同一个函数的示例：

```
beego.Router("/api",&RestController{},"get,post:ApiFunc")
```

以下是不同的 method 对应不同的函数，通过 ; 进行分割的示例：

```
beego.Router("/simple",&SimpleController{},"get:GetFunc;post:PostFunc")
```

可用的 HTTP Method：

包含以下所有的函数

```
* get: GET 请求
* post: POST 请求
* put: PUT 请求
* delete: DELETE 请求
* patch: PATCH 请求
* options: OPTIONS 请求
* head: HEAD 请求
```

如果同时存在 * 和对应的 HTTP Method，那么优先执行 HTTP Method 的方法，例如同时注册了如下所示的路由：

```
beego.Router("/simple",&SimpleController{},"`*:AllFunc;post:PostFunc`")
```

那么执行 `POST` 请求的时候，执行 `PostFunc` 而不执行 `AllFunc`。

> > > 自定义函数的路由默认不支持 RESTful 的方法，也就是如果你设置了 `beego.Router("/api",&RestController{},"post:ApiFunc")` 这样的路由，如果请求的方法是 `POST`，那么不会默认去执行 `Post` 函数。



#### 2.1.2.6 自动匹配

用户首先需要把需要路由的控制器注册到自动路由中：

```
beego.AutoRouter(&controllers.ObjectController{})
```

那么 beego 就会通过反射获取该结构体中所有的实现方法，你就可以通过如下的方式访问到对应的方法中：

```
/object/login   调用 ObjectController 中的 Login 方法
/object/logout  调用 ObjectController 中的 Logout 方法
```

除了前缀两个 `/:controller/:method` 的匹配之外，剩下的 url beego 会帮你自动化解析为参数，保存在 `this.Ctx.Input.Params` 当中：

```
/object/blog/2013/09/12  调用 ObjectController 中的 Blog 方法，参数如下：map[0:2013 1:09 2:12]
```

方法名在内部是保存了用户设置的，例如 Login，url 匹配的时候都会转化为小写，所以，`/object/LOGIN` 这样的 `url` 也一样可以路由到用户定义的 `Login` 方法中。

现在已经可以通过自动识别出来下面类似的所有 url，都会把请求分发到 `controller` 的 `simple` 方法：

```
/controller/simple
/controller/simple.html
/controller/simple.json
/controller/simple.xml
```

可以通过 `this.Ctx.Input.Param(":ext")` 获取后缀名。
