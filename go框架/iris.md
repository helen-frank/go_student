# 1.关于iris

**`Iris`** 是一个通过 Go 编写的快速的，简单的，但是功能齐全和非常有效率的 web 框架。

**`Iris`** 为你下一个网站或者 API 提供了一个精美的、使用简单的基础。

**`Iris`** 为它的使用者提供了一个完整且体面的支持。

## 1.1 iris的哲学

**`Iris`** 的哲学是为 HTTP 提供强大的工具，使其成为单页应用、网站或者公共 HTTP API的好的解决方案。记住，目前为止，就实际性能而言，**`Iris`** 是至今为止最快的 web 框架。

**`Iris`** 不会强制你使用任何特定的 **`ORM`** 或者**`模板引擎`**。支持最强大和快速的模板引擎，你可以快速开发出完美的应用程序。

## 1.2. 为什么还要使用其他 Web 框架呢?

Go 是一个伟大的技术栈，可用来为 Web 应用构建可扩展的、基于网络的后端系统。

当你考虑用 Go 构建 web 应用程序和 web API，或者简单构建 HTTP 服务器，你是否考虑过标准库 **`net/http`** ? 然后你不得不解决一些常见的情况，例如静态路由、安全和用户认证，实时通信和许多其他问题，而这些问题 **`net/http`** 无法解决。

**`net/http`** 还不够完整，无法用来快速构建设计良好的 web 后端系统。当你意识到这个的时候，你可以思考下下面的话：

- **`net/http`** 不适合我，但是有许多框架，我选择哪个呢？
- 每个框架都告诉我它是最好的。我不知道该怎么选择。

# 2.安装

```bash
go get github.com/kataras/iris/v12@latest
```

 

# 3.使用

**`Iris`** 具有用于路由的表达语法，感觉就像在家一样。路由算法是强大的，通过 [muxie](https://github.com/kataras/muxie) 来处理请求和比其替代品(httprouter、gin、echo) 更快地匹配路由。

`main.go`

```go
package main

import "github.com/kataras/iris/v12"

func main() {
    app := iris.Default()
    app.Use(myMiddleware)

    app.Handle("GET", "/ping", func(ctx iris.Context) {
        ctx.JSON(iris.Map{"message": "pong"})
    })

    // Listens and serves incoming http requests
    // on http://localhost:8080.
    app.Run(iris.Addr(":8080"))
}

func myMiddleware(ctx iris.Context) {
    ctx.Application().Logger().Infof("Runs before %s", ctx.Path())
    ctx.Next()
}
```



## 3.1 简单实例

`main.go`

```go
package main

import "github.com/kataras/iris/v12"

func main() {
    app := iris.New()
    // Load all templates from the "./views" folder
    // where extension is ".html" and parse them
    // using the standard `html/template` package.
    app.RegisterView(iris.HTML("./views", ".html"))

    // Method:    GET
    // Resource:  http://localhost:8080
    app.Get("/", func(ctx iris.Context) {
        // Bind: {{.message}} with "Hello world!"
        ctx.ViewData("message", "Hello world!")
        // Render template file: ./views/hello.html
        ctx.View("hello.html")
    })

    // Method:    GET
    // Resource:  http://localhost:8080/user/42
    //
    // Need to use a custom regexp instead?
    // Easy;
    // Just mark the parameter's type to 'string'
    // which accepts anything and make use of
    // its `regexp` macro function, i.e:
    // app.Get("/user/{id:string regexp(^[0-9]+$)}")
    app.Get("/user/{id:uint64}", func(ctx iris.Context) {
        userID, _ := ctx.Params().GetUint64("id")
        ctx.Writef("User ID: %d", userID)
    })

    // Start the server using a network address.
    app.Run(iris.Addr(":8080"))
}
```

`./views/hello.html`

```html
<html>
<head>
    <title>Hello Page</title>
</head>
<body>
    <h1>{{ .message }}</h1>
</body>
</html>
```

## 3.2 iris热更新

想要当你的代码发生改变时自动重启你的程序？安装 [rizla](https://github.com/kataras/rizla) 工具，执行 `rizla main.go` 来替代 `go run main.go`。

唯一要求，go version > 1.7

```bash
go get -u github.com/kataras/rizla
```



# 4.Host

你可以开启服务监听任何 `net.Listener` 或者 `http.Server` 类型的实例。初始化服务器的方法应该在最后传递给 `Run` 函数。

Go 开发者最常用的方法是通过传递一个形如 `hostname:ip` 形式的网络地址来开启一个服务。**`Iris`** 中我们使用`iris.Addr`，它是一个 `iris.Runner` 类型。

```go
// Listening on tcp with network address 0.0.0.0:8080
app.Run(iris.Addr(":8080"))
```

有时候你在你的应用程序的其他地方创建一个标准库 `het/http` 服务器，并且你想使用它作为你的 Iris web 程序提供服务。

```go
// Same as before but using a custom http.Server which may being used somewhere else too
app.Run(iris.Server(&http.Server{Addr:":8080"}))
```

最高级的用法是创建一个自定义的或者标准的 `net/Listener`，然后传递给 `app.Run`。

```go
// Using a custom net.Listener
l, err := net.Listen("tcp4", ":8080")
if err != nil {
    panic(err)
}
app.Run(iris.Listener(l))
```

一个更加完整的示例，使用的是仅Unix套接字文件特性

```go
package main

import (
    "os"
    "net"

    "github.com/kataras/iris/v12"
)

func main() {
    app := iris.New()

    // UNIX socket
    if errOs := os.Remove(socketFile); errOs != nil && !os.IsNotExist(errOs) {
        app.Logger().Fatal(errOs)
    }

    l, err := net.Listen("unix", socketFile)

    if err != nil {
        app.Logger().Fatal(err)
    }

    if err = os.Chmod(socketFile, mode); err != nil {
        app.Logger().Fatal(err)
    }

    app.Run(iris.Listener(l))
}
```

UNIX 和 BSD 主机可以使用重用端口的功能。

```go
package main

import (
    // Package tcplisten provides customizable TCP net.Listener with various
    // performance-related options:
    //
    //   - SO_REUSEPORT. This option allows linear scaling server performance
    //     on multi-CPU servers.
    //     See https://www.nginx.com/blog/socket-sharding-nginx-release-1-9-1/ for details.
    //
    //   - TCP_DEFER_ACCEPT. This option expects the server reads from the accepted
    //     connection before writing to them.
    //
    //   - TCP_FASTOPEN. See https://lwn.net/Articles/508865/ for details.
    "github.com/valyala/tcplisten"

    "github.com/kataras/iris/v12"
)

// go get github.com/valyala/tcplisten
// go run main.go

func main() {
    app := iris.New()

    app.Get("/", func(ctx iris.Context) {
        ctx.HTML("<h1>Hello World!</h1>")
    })

    listenerCfg := tcplisten.Config{
        ReusePort:   true,
        DeferAccept: true,
        FastOpen:    true,
    }

    l, err := listenerCfg.NewListener("tcp", ":8080")
    if err != nil {
        app.Logger().Fatal(err)
    }

    app.Run(iris.Listener(l))
}
```



## 4.1 HTTP/2和安全

如果你有签名文件密钥，你可以使用 `iris.TLS` 基于这些验证密钥开启 `https` 服务。

```go
// TLS using files
app.Run(iris.TLS("127.0.0.1:443", "mycert.cert", "mykey.key"))
```

当你的应用准备部署**生产**时，你可以使用 `iris.AutoTLS` 方法，它通过[ https://letsencrypt.org](https://letsencrypt.org/) 免费提供的证书来开启一个安全的服务。

```go
// Automatic TLS
app.Run(iris.AutoTLS(":443", "example.com", "admin@example.com"))
```



## 4.2 任意iris.Runner

有时你想要监听一些特定的东西，并且这些东西不是 `net.Listener` 类型的。你能够通过 `iris.Raw` 方法做到，但是你得对此方法负责。

```go
// Using any func() error,
// the responsibility of starting up a listener is up to you with this way,
// for the sake of simplicity we will use the
// ListenAndServe function of the `net/http` package.
app.Run(iris.Raw(&http.Server{Addr:":8080"}).ListenAndServe)
```



## 4.3 host配置

形如上面所示的监听方式都可以在最后接受一个 **`func(\*iris.Supervisor)`** 的可变参数。通过函数的传递用来为特定 host 添加配置器。

例如，我们想要当服务器关闭的时候触发的回调函数：

```go
app.Run(iris.Addr(":8080", func(h *iris.Supervisor) {
    h.RegisterOnShutdown(func() {
        println("server terminated")
    })
}))
```

你甚至可以在再 `app.Run` 之前配置，但是不同的是，这个 host 配置器将会在所有的主机上执行(我们将在稍后看到 `app.NewHost` )

```go
app := iris.New()
app.ConfigureHost(func(h *iris.Supervisor) {
    h.RegisterOnShutdown(func() {
        println("server terminated")
    })
})
app.Run(iris.Addr(":8080"))
```

当 `Run` 方法运行之后，通过 `Application#Hosts` 字段提供的所有 hosts 你的应用服务都可以访问。

但是最常用的场景是你可能需要在运行 `app.Run` 之前访问 hosts，这里有2种方法来获得访问 hosts 的监管，阅读下面。

我们已经看到通过 `app.Run` 的第二个参数或者 `app.ConfigureHost` 方法来配置所有的应用程序 hosts 。还有一种更加适合简单场景的方法，那就是使用 `app.NewHost` 来创建一个新的 host，然后使用它的 `Serve` 或者 `Listen` 函数， 通过 `iris#Raw` 来启动服务。

记住这个方法需要额外导入 `net/http` 包。

示例代码：

```go
h := app.NewHost(&http.Server{Addr:":8080"})
h.RegisterOnShutdown(func(){
    println("server terminated")
})

app.Run(iris.Raw(h.ListenAndServe))
```



## 4.4 多个主机

你可以使用多个 hosts 来启动你的 iris 程序，`iris.Router` 兼容 `net/http/Handler` 函数，因此我们可以理解为，它可以被适用于任何 `net/http` 服务器，然而，通过使用 `app.NewHost` 是一个更加简单的方法，它也会复制所有的 host 配置器，并在 `app.Shutdown` 时关闭所有依附在特定web 服务的主机 host 。 `app := iris.New() app.Get("/", indexHandler)`

```go
// run in different goroutine in order to not block the main "goroutine".
go app.Run(iris.Addr(":8080"))
// start a second server which is listening on tcp 0.0.0.0:9090,
// without "go" keyword because we want to block at the last server-run.
app.NewHost(&http.Server{Addr:":9090"}).ListenAndServe()
```



## 4.5 优雅地关闭

让我们继续学习怎么接受 `CONTROL + C` / `COMMAND + C` 或者 unix `kill` 命令，优雅地关闭服务器。(默认是启用的)

为了手动地管理app被中断时需要做的事情，我们需要通过使用 `WithoutInterruptHandler` 选项禁用默认的行为，然后注册一个新的中断处理器(在所有可能的hosts上)。

示例代码：

```go
package main

import (
    "context"
    "time"

    "github.com/kataras/iris/v12"
)

func main() {
    app := iris.New()

    iris.RegisterOnInterrupt(func() {
        timeout := 5 * time.Second
        ctx, cancel := context.WithTimeout(context.Background(), timeout)
        defer cancel()
        // close all hosts
        app.Shutdown(ctx)
    })

    app.Get("/", func(ctx iris.Context) {
        ctx.HTML(" <h1>hi, I just exist in order to see if the server is closed</h1>")
    })

    app.Run(iris.Addr(":8080"), iris.WithoutInterruptHandler)
}
```



# 5.配置

`iris.New` 函数返回一个 `iris.Application` 实例。这个实例可以通过它的 `Configure(...iris.Configurator)` 和 `Run` 方法进行配置。

`app.Run` 方法的第二个参数是可选的、可变长的，接受一个或者多个 `iris.Configurator`。一个 `iris.Configurator` 是 `func(app *iris.Application)` 类型的函数。自定义的 `iris.Configurator` 能够修改你的 `*iris.Application`。

每个核心的配置字段都有一个内建的 `iris.Configurator`。例如， `iris.WithoutStartupLog`，`iris.WithCharset("UTF-8")`，`iris.WithOptimizations`，`iris.WithConfiguration(iris.Congiguration{...})` 函数。

每个模块，例如视图引擎，websockets，sessions和每个中间件都有它们各自配置器和选项，它们大多数都与核心的配置分离。

## 5.1 使用配置

唯一的配置结构体是 `iris.Configuration`。让我们从通过 `iris.WithConfiguration` 函数来创建 `iris.Configurator` 开始。

所有的 `iris.Configuration` 字段的默认值都是最常用的。 **`Iris`** 在 `app.Run` 运行之前不需要任何的配置。但是你想要在运行服务之前使用自定义的 `iris.Configurator`，你可以把你的配置器传递给`app.Configure` 方法。

```o
 config := iris.WithConfiguration(iris.Configuration {
  DisableStartupLog: true,
  Optimizations: true,
  Charset: "UTF-8",
})

app.Run(iris.Addr(":8080"), config)
```



## 5.2 从YAML加载

使用 `iris.YAML("path")`

File：iris.yml

```yaml
FireMethodNotAllowed: true
DisableBodyConsumptionOnUnmarshal: true
TimeFormat: Mon, 01 Jan 2006 15:04:05 GMT
Charset: UTF-8
```

File：main.go

```go
config := iris.WithConfiguration(iris.YAML("./iris.yml"))
app.Run(iris.Addr(":8080"), config)
```



## 5.3 从TOML加载

使用 `iris.TOML("path")`

File：iris.tml

```toml
FireMethodNotAllowed = true
DisableBodyConsumptionOnUnmarshal = false
TimeFormat = "Mon, 01 Jan 2006 15:04:05 GMT"
Charset = "UTF-8"

[Other]
    ServerName = "my fancy iris server"
    ServerOwner = "admin@example.com"
```

File：main.go

```go
config := iris.WithConfiguration(iris.TOML("./iris.tml"))
app.Run(iris.Addr(":8080"), config)
```



## 5.4 使用函数方式

我们已经提到，你可以传递任何数量的 `iris.Configurator` 到 `app.Run` 的第二个参数。I**`Iris`** 为每个`iris.Configuration` 的字段提供了一个选项。

```go
app.Run(iris.Addr(":8080"), iris.WithoutInterruptHandler,
    iris.WithoutServerError(iris.ErrServerClosed),
    iris.WithoutBodyConsumptionOnUnmarshal,
    iris.WithoutAutoFireStatusCode,
    iris.WithOptimizations,
    iris.WithTimeFormat("Mon, 01 Jan 2006 15:04:05 GMT"),
)
```

当你想要改变一些 `iris.Configuration` 的字段的时候，这是一个很好的做法。通过前缀 :`With` 或者 `Without`，代码编辑器能够帮助你浏览所有的配置选项，甚至你都不需要翻阅文档。



## 5.5 自定义值

`iris.Configuration` 包含一个名为 `Other map[string]interface` 的字段，它可以接受任何自定义的 `key:value` 选项，因此你可以依据需求使用这个字段来传递程序需要的指定的值。

```go
app.Run(iris.Addr(":8080"), 
    iris.WithOtherValue("ServerName", "my amazing iris server"),
    iris.WithOtherValue("ServerOwner", "admin@example.com"),
)
```

你可以通过 `app.ConfigurationReadOnly` 来访问这些字段。

```go
serverName := app.ConfigurationReadOnly().Other["MyServerName"]
serverOwner := app.ConfigurationReadOnly().Other["ServerOwner"]
```



## 5.6 从Context中配置

在一个处理器中，通过下面的方式访问这些字段。

```go
ctx.Application().ConfigurationReadOnly()
```



# 6.路由

### 6.0.1 处理器类

处理器，物如其名，就是处理请求的东西。

一个处理器响应一个 HTTP 请求。它写入响应头和数据到 `Context.ResponseWriter()`，然后再返回。返回信号表明请求已经完成；在处理完成调用当时或者之后使用 `Context` 是无效的。

由于 HTTP 客户端软件，HTTP 协议版本，和任何客户端和 iris 服务器中间媒介等因素，在向 `context.ResponseWriter()` 写入后可能无法从 `context.Request().Body` 中读取数据。注意应该先从`context.Request().Body` 中读取数据后，然后再响应它。

除了读取请求体，处理器不应该改变提供的 `context`

如果处理器出现 `panic`，服务器(处理器的调用者)会假定这个 panic 的影响与存活的请求无关。它会 recover 这个 panic，记录栈追踪日志到服务器错误日志中，并中断连接。

```go
type Handler func(iris.Context)
```

一旦处理器被注册，我们可以给返回的 `路由` 实例指定一个名字，以便更加容易地调试和在视图中匹配相对路径。更多信息，查看 `反向查询(Reverse Lookups)` 章节。

### 6.0.2 行为

**`Iris`** 默认接受和注册形如 `/api/user` 这样的路径的路由，且尾部不带斜杠。如果客户端尝试访问 `$your_host/api/user/`，**`Iris`** 路由会自动永久重定向(301)到 `$your_host/api/user`，以便由注册的路由进行处理。这是设计 APIs的现代化的方式。

然而如果你想禁用请求的 `路径更正` 的功能的话，你可以在 `app.Run` 传递 `iris.WithoutPathCorrection` 配置选项。例如：

```go
 // [app := iris.New...]
// [...]

app.Run(iris.Addr(":8080"), iris.WithoutPathCorrection)
```

如果你想 `/api/user` 和 `/api/user/` 在不重定向的情况下(常用场景)拥有相同的处理器，只需要 `iris.WithoutPathCorrectionRedirection` 选项即可：

```go
app.Run(iris.Addr(":8080"), iris.WithoutPathCorrectionRedirection)
```



### 6.0.3 API

支持所有的 HTTP 方法，开发者也可以在相同路劲的不同的方法注册处理器(比如 `/user` 的 GET 和 POST)。

第一个参数是 HTTP 方法，第二个参数是请求的路径，第三个可变参数应该包含一个或者多个`iris.Handler`，当客户端请求到特定的资源路径时，这些处理器将会按照注册的顺序依次执行。

示例代码：

```html
app := iris.New()

app.Handle("GET", "/contact", func(ctx iris.Context) {
    ctx.HTML("<h1> Hello from /contact </h1>")
})
```

为了让后端开发者做事更加容易， **`Iris`** 为所有的 HTTP 方法提供了"帮手"。第一个参数是路由的请求路径，第二个可变参数是一个或者多个 `iris.Handler`，也会按照注册的顺序依次执行。

示例代码:

```go
app := iris.New()

// Method: "GET"
app.Get("/", handler)

// Method: "POST"
app.Post("/", handler)

// Method: "PUT"
app.Put("/", handler)

// Method: "DELETE"
app.Delete("/", handler)

// Method: "OPTIONS"
app.Options("/", handler)

// Method: "TRACE"
app.Trace("/", handler)

// Method: "CONNECT"
app.Connect("/", handler)

// Method: "HEAD"
app.Head("/", handler)

// Method: "PATCH"
app.Patch("/", handler)

// register the route for all HTTP Methods
app.Any("/", handler)

func handler(ctx iris.Context){
    ctx.Writef("Hello from method: %s and path: %s\n", ctx.Method(), ctx.Path())
}
```



### 6.0.4 离线路由

在 **`Iris`** 中有一个特殊的方法你可以使用。它被成为 `None`，你可以使用它向外部隐藏一条路由，但仍然可以从其他路由处理中通过 `Context.Exec` 方法调用它。每个 API 处理方法返回 Route 值。一个 Route 的 `IsOnline` 方法报告那个路由的当前状态。你可以通过它的 `Route.Method` 字段的值来改变路由 `离线` 状态为 `在线` 状态，反之亦然。当然每次状态的改变需要调用 `app.RefreshRouter()` 方法，这个使用是安全的。看看下面一个完整的例子：

```go
// file: main.go
package main

import (
    "github.com/kataras/iris/v12"
)

func main() {
    app := iris.New()

    none := app.None("/invisible/{username}", func(ctx iris.Context) {
        ctx.Writef("Hello %s with method: %s", ctx.Params().Get("username"), ctx.Method())

        if from := ctx.Values().GetString("from"); from != "" {
            ctx.Writef("\nI see that you're coming from %s", from)
        }
    })

    app.Get("/change", func(ctx iris.Context) {

        if none.IsOnline() {
            none.Method = iris.MethodNone
        } else {
            none.Method = iris.MethodGet
        }

        // refresh re-builds the router at serve-time in order to
        // be notified for its new routes.
        app.RefreshRouter()
    })

    app.Get("/execute", func(ctx iris.Context) {
        if !none.IsOnline() {
            ctx.Values().Set("from", "/execute with offline access")
            ctx.Exec("NONE", "/invisible/iris")
            return
        }

        // same as navigating to "http://localhost:8080/invisible/iris"
        // when /change has being invoked and route state changed
        // from "offline" to "online"
        ctx.Values().Set("from", "/execute")
        // values and session can be
        // shared when calling Exec from a "foreign" context.
        //     ctx.Exec("NONE", "/invisible/iris")
        // or after "/change":
        ctx.Exec("GET", "/invisible/iris")
    })

    app.Run(iris.Addr(":8080"))
}
```



**怎么运行**

1. 运行 `go run main.go`
2. 打开浏览器访问 `http://localhost:8080/invisible/iris`，你将会看到 `404 not found` 的错误。
3. 然而 `http://localhost:8080/execute` 将会执行这个路由。
4. 现在，如果你导航至 `http://localhost:8080/change` ，然后刷新`/invisible/iris` 选项卡，你将会看到它。



### 6.0.5 路由组

一些列路由可以通过路径的前缀(可选的)进行分组，共享相同的中间件处理器和模板布局。一个组也可以有一个内嵌的组。

`.Party` 用来路由分组，开发者可以声明不限数量的组。

示例代码：

```go
app := iris.New()

users := app.Party("/users", myAuthMiddlewareHandler)

// http://localhost:8080/users/42/profile
users.Get("/{id:uint64}/profile", userProfileHandler)
// http://localhost:8080/users/messages/1
users.Get("/messages/{id:uint64}", userMessageHandler)
```

你可以使用 `PartyFunc` 方法编写相同的内容，它接受子路由器或者Party。

```go
app := iris.New()

app.PartyFunc("/users", func(users iris.Party) {
    users.Use(myAuthMiddlewareHandler)

    // http://localhost:8080/users/42/profile
    users.Get("/{id:uint64}/profile", userProfileHandler)
    // http://localhost:8080/users/messages/1
    users.Get("/messages/{id:uint64}", userMessageHandler)
})
```



### 6.0.6 路径参数

与你见到的其他路由器不同， Iris 的路由器可以处理各种路由路径而不会发生冲突。

只匹配 `GET "/"`

```go
app.Get("/", indexHandler)
```

下面所示能匹配所有以 `/assets/**/*` 前缀的 GET 请求，它是一个通配符，`ctx.Params().Get("asset")`获取 `/assets/` 后面的所有路径。

```go
app.Get("/assets/{asset:path}", assetsWildcardHandler)
```

下面所示的能匹配所有以 `/profile/` 前缀的 GET 请求，但是获取的只是单个路径的部分。

```go
app.Get("/profile/{username:string}", userHandler)
```

下面所示的只能匹配 `/profile/me` 的 GET 请求，它不与 `/profile/{username:string}` 或者 `/{root:path}` 冲突。

```go
app.Get("/profile/me", userHandler)
```

下面所示的能匹配所有以 `/users/` 前缀的 GET 请求，并且后面的是一个数字，且数字要大于等于1。

```go
app.Get("/user/{userid:int min(1)}", getUserHandler)
```

下面所示的能匹配所有以 `/users/` 前缀的 DELETE 请求，并且后面的是一个数字，且数字要大于等于1。

```go
app.DELETE("/user/{userid:int min(1)}", getUserHandler)
```

下面所示的能匹配所有除了被其他路由器处理的 GET 请求。例如在这种情况下，上面的路由 ：`/`，`/assets/{asset:path}`，`/profile/{username}`，`/profile/me`，`/user/{userid:int ...}`，它将不会与余下的路由冲突。

```go
app.Get("{root:path}", rootWildcardHandler)
```

匹配所有的请求：

1. `/u/adcd` 映射为 `:alphabetical` (如果 `:alphabetical` 注册， 否则 `:string`)

2. `/u/42` 映射为 `:uint` (如果 `:uint` 注册， 否则 `:int`)

3. `/u/-1` 映射为 `:int` (如果 `:int` 注册， 否则 `:string`)

4. `/u/adcd123` 映射为 `:string`

   ```go
   app.Get("/u/{username:string}", func(ctx iris.Context) {
       ctx.Writef("username (string): %s", ctx.Params().Get("username"))
   })    
   app.Get("/u/{id:int}", func(ctx iris.Context) {
       ctx.Writef("id (int): %d", ctx.Params().GetIntDefault("id", 0))
   })
   app.Get("/u/{uid:uint}", func(ctx iris.Context) {
       ctx.Writef("uid (uint): %d", ctx.Params().GetUintDefault("uid", 0))
   })    
   app.Get("/u/{firstname:alphabetical}", func(ctx iris.Context) {
       ctx.Writef("firstname (alphabetical): %s", ctx.Params().Get("firstname"))
   })
   ```

分别匹配 `/abctenchars.xml` 和 `/abcdtenchars` 的所有 GET 请求

```go
app.Get("/{alias:string regexp(^[a-z0-9]{1,10}\\.xml$)}", PanoXML)
app.Get("/{alias:string regexp(^[a-z0-9]{1,10}$)}", Tour)
```

你可能知道 `{id:uint64}`、`:path`、`min(1)` 是什么。它们是在注册时可以键入的动态参数和函数。了解更多请阅读 `路径参数类型(Path Parameter Types)`



## 6.1 路由参数类型

**`Iris`** 拥有你见过的的最简单和强大路由处理。

**`Iris`** 自己拥有用于路由路径语法解析和判定的解释器(就像一门编程语言)。

它是快速的。它计算它的需求，如果没有特殊的正则需要，它仅仅会使用低级的路径语法来注册路由，除此之外，它预编译正则，然后加入到必需的中间件中。这就意味着与其他的路由器或者 web 框架相比，你的性能成本为零。

### 6.1.1 参数

一个路径参数的名字应该仅仅包含 `字母`。数字和类似 `"_"` 这样的符号是 `不允许的`。

不要迷惑于 `ctx.Params()` 和 `ctx.Values()`

- 路径的参数值可以通过 `ctx.Params()` 取出。
- `ctx` 中用于处理器与中间件之间通信的本地存储可以存储在 `ctx.Values()` 中。

下表是内建可用的参数类型：

| 参数类型        |  golang类型  |                           取值范围                           | 取值方式                 |
| :-------------- | :----------: | :----------------------------------------------------------: | ------------------------ |
| **`:string`**   | **`string`** |                     任何值(单个字段路径)                     | **`Params().Get`**       |
| **`:int`**      |  **`int`**   | -9223372036854775808 - 9223372036854775807 (x64) </br>-2147483648 - 2147483647 (x32) | **`Params().GetInt`**    |
| **`:int8`**     |  **`int8`**  |                          -128 - 127                          | **`Params().GetInt8`**   |
| **`:int16`**    | **`int16`**  |                        -32768 - 32767                        | **`Params().GetInt16`**  |
| **`:int32`**    | **`int32`**  |                   -2147483648 - 2147483647                   | **`Params().GetInt32`**  |
| **`:int64`**    | **`int64`**  |          -9223372036854775808 - 9223372036854775807          | **`Params().GetInt64`**  |
| **`:uint8`**    | **`uint8`**  |                           0 - 255                            | **`Params().GetUint8`**  |
| **`:uint16`**   | **`uint16`** |                          0 - 65535                           | **`Params().GetUint16`** |
| **`:uint32`**   | **`uint32`** |                        0 - 4294967295                        | **`Params().GetUint32`** |
| **`:uint64`**   | **`uint64`** |                   0 - 18446744073709551615                   | **`Params().GetUint64`** |
| **`:bool`**     |  **`bool`**  | "1"，"t"，"T"，"TRUE"，"true"，"True"，"0"，"f"， "F"， "FALSE",，"false"，"False" | **`Params().GetBool`**   |
| `:alphabetical` | **`string`** |                        小写或大写字母                        | **`Params().Get`**       |
| **`:file`**     | **`string`** | 大小写字母，数字，下划线(_)，横线(-)，点(.)，以及没有空格或者其他对文件名无效的特殊字符 | **`Params().Get`**       |
| **`:path`**     | **`string`** |  任何可以被斜线(/)分隔的路径段，但是应该为路由的最后一部分   | **`Params().Get`**       |

**示例：**

```go
app.Get("/users/{id:uint64}", func(ctx iris.Context){
    id := ctx.Params().GetUint64Default("id", 0)
    // [...]
})
```



#### 6.1.1.1 内建参数

|                           内建函数                           |                           参数类型                           |
| :----------------------------------------------------------: | :----------------------------------------------------------: |
|                  **`regexp(expr string)`**                   |                        **`:string`**                         |
|                 **`prefix(prefix string)`**                  |                        **`:string`**                         |
|                 **`suffix(suffix string)`**                  |                        **`:string`**                         |
|                   **`contains(s string)`**                   |                        **`:string`**                         |
| **`min(最小值)，接收：int，int8，int16，int32，int64，uint8`**， **`uint16，uint32，uint64，float32，float64)`** | **`:string(字符长度)，:int，:int16，:int32，:int64`**，</br>**`:uint，:uint16，:uint32，:uint64`** |
| **`max(最大值)，接收：int，int8，int16，int32，int64，uint8`**， **`uint16，uint32，uint64，float32，float64)`** | **`:string(字符长度)，:int，:int16，:int32，:int64`**，</br>**`:uint，:uint16，:uint32，:uint64`** |
| **`range(最小值，最大值)，接收：int，int8，int16，int32，int64，uint8`**， **`uint16，uint32，uint64，float32，float64)`** | **`:int，:int16，:int32，:int64`**，</br>**`:uint，:uint16，:uint32，:uint64`** |

**示例**：

```go
app.Get("/profile/{name:alphabetical max(255)}", func(ctx iris.Context){
    name := ctx.Params().Get("name")
    // len(name) <=255 otherwise this route will fire 404 Not Found
    // and this handler will not be executed at all.
})
```



#### 6.1.1.2 自己做

`RegisterFunc` 可以接受任何返回 `func(paramValue string) bool` 的函数。如果验证失败将会触发 `404` 或者任意 `else`关键字拥有的状态码。

```go
latLonExpr := "^-?[0-9]{1,3}(?:\\.[0-9]{1,10})?$"
latLonRegex, _ := regexp.Compile(latLonExpr)

// Register your custom argument-less macro function to the :string param type.
// MatchString is a type of func(string) bool, so we use it as it is.
app.Macros().Get("string").RegisterFunc("coordinate", latLonRegex.MatchString)

app.Get("/coordinates/{lat:string coordinate()}/{lon:string coordinate()}",
func(ctx iris.Context) {
    ctx.Writef("Lat: %s | Lon: %s", ctx.Params().Get("lat"), ctx.Params().Get("lon"))
})
```

注册接受两个 `int` 参数的自定义的宏函数。

```go
app.Macros().Get("string").RegisterFunc("range",
func(minLength, maxLength int) func(string) bool {
    return func(paramValue string) bool {
        return len(paramValue) >= minLength && len(paramValue) <= maxLength
    }
})

app.Get("/limitchar/{name:string range(1,200) else 400}", func(ctx iris.Context) {
    name := ctx.Params().Get("name")
    ctx.Writef(`Hello %s | the name should be between 1 and 200 characters length
    otherwise this handler will not be executed`, name)
})
```

注册接受一个 `[]string` 参数的自定义的宏函数。

```go
app.Macros().Get("string").RegisterFunc("has",
func(validNames []string) func(string) bool {
    return func(paramValue string) bool {
        for _, validName := range validNames {
            if validName == paramValue {
                return true
            }
        }

        return false
    }
})

app.Get("/static_validation/{name:string has([kataras,maropoulos])}",
func(ctx iris.Context) {
    name := ctx.Params().Get("name")
    ctx.Writef(`Hello %s | the name should be "kataras" or "maropoulos"
    otherwise this handler will not be executed`, name)
})
```

示例代码：

```go
func main() {
    app := iris.Default()

    // This handler will match /user/john but will not match neither /user/ or /user.
    app.Get("/user/{name}", func(ctx iris.Context) {
        name := ctx.Params().Get("name")
        ctx.Writef("Hello %s", name)
    })

    // This handler will match /users/42
    // but will not match /users/-1 because uint should be bigger than zero
    // neither /users or /users/.
    app.Get("/users/{id:uint64}", func(ctx iris.Context) {
        id := ctx.Params().GetUint64Default("id", 0)
        ctx.Writef("User with ID: %d", id)
    })

    // However, this one will match /user/john/send and also /user/john/everything/else/here
    // but will not match /user/john neither /user/john/.
    app.Post("/user/{name:string}/{action:path}", func(ctx iris.Context) {
        name := ctx.Params().Get("name")
        action := ctx.Params().Get("action")
        message := name + " is " + action
        ctx.WriteString(message)
    })

    app.Run(iris.Addr(":8080"))
}
```

当没有指定参数类型时，默认为 `string`，因此 `{name:string}` 和 `{name}` 是完全相同的。



## 6.2 反向查询

**`Iris`** 提供了一些处理器注册方法，每个方法都返回一个 `Route` 实例。

### 6.2.1 路由命名

路由命名非常简单，我们只需要调用返回的 `*Route` 的 `Name` 字段来定义名字。

```go
package main

import "github.com/kataras/iris/v12"

func main() {
    app := iris.New()
    // define a function
    h := func(ctx iris.Context) {
        ctx.HTML("<b>Hi</b1>")
    }

    // handler registration and naming
    home := app.Get("/", h)
    home.Name = "home"
    // or
    app.Get("/about", h).Name = "about"
    app.Get("/page/{id}", h).Name = "page"g

    app.Run(iris.Addr(":8080"))
}
```



### 6.2.2 从路由名字生成URL

当我们为特定的路径注册处理器时，我们可以根据传递给 Iris 的结构化数据创建 URLs。如上面的例子所示，我们命名了三个路由，其中之一甚至带有参数。如果我们使用默认的 `html/template` 视图引擎，我们可以使用一个简单的操作来反转路由(生成示例的 URLs)：

```go
Home: {{ urlpath "home" }}
About: {{ urlpath "about" }}
Page 17: {{ urlpath "page" "17" }}
```

上面的代理可以生成下面的输出：

```go
Home: http://localhost:8080/ 
About: http://localhost:8080/about
Page 17: http://localhost:8080/page/17
```



### 6.2.3 在代码中使用路由名字

我们可以使用以下方法/函数来处理命名路由（及其参数）：

- `GetRoutes` 函数获取所有注册的路由
- `GetRoute(routeName string)` 方法通过名字获得路由
- `URL(routeName string, paramValues ...interface{})` 方法通过提供的值来生成 URL 字符串
- `Path(routeName string, paramValues ...interface{})` 方法通过提供的值生成URL的路径部分(没有主机地址和协议)。

## 6.3 中间件

当我们谈论 Iris 中的中间件时，我们谈论的是一个 HTTP 请求的生命周期中主处理器代码运行前/后运行的代码。例如，日志中间件可能记录一个传入请求的详情到日志中，然后调用处理器代码，然后再编写有关响应的详细信息到日志中。关于中间件的一件很酷的事情是，这些单元非常灵活且可重复使用。

中间件仅是一个 `Handler` 格式的函数 `func(ctx iris.Context)`，当前一个中间件调用 `ctx.Next()` 方法时，此中间件被执行，这可以用作身份验证，即如果请求验证通过，就调用 `ctx.Next()` 来执行该请求剩下链上的处理器，否则触发一个错误响应。

### 6.3.1 编写一个中间件

```go
package main

import "github.com/kataras/iris/v12"

func main() {
    app := iris.New()
    // or app.Use(before) and app.Done(after).
    app.Get("/", before, mainHandler, after)
    app.Run(iris.Addr(":8080"))
}

func before(ctx iris.Context) {
    shareInformation := "this is a sharable information between handlers"

    requestPath := ctx.Path()
    println("Before the mainHandler: " + requestPath)

    ctx.Values().Set("info", shareInformation)
    ctx.Next() // execute the next handler, in this case the main one.
}

func after(ctx iris.Context) {
    println("After the mainHandler")
}

func mainHandler(ctx iris.Context) {
    println("Inside mainHandler")

    // take the info from the "before" handler.
    info := ctx.Values().GetString("info")

    // write something to the client as a response.
    ctx.HTML("<h1>Response</h1>")
    ctx.HTML("<br/> Info: " + info)

    ctx.Next() // execute the "after".
}
```



### 6.3.2 全局范围

```go
package main

import "github.com/kataras/iris/v12"

func main() {
    app := iris.New()

    // register our routes.
    app.Get("/", indexHandler)
    app.Get("/contact", contactHandler)

    // Order of those calls does not matter,
    // `UseGlobal` and `DoneGlobal` are applied to existing routes
    // and future routes also.
    //
    // Remember: the `Use` and `Done` are applied to the current party's and its children,
    // so if we used the `app.Use/Done before the routes registration
    // it would work like UseGlobal/DoneGlobal in this case,
    // because the `app` is the root "Party".
    app.UseGlobal(before)
    app.DoneGlobal(after)

    app.Run(iris.Addr(":8080"))
}

func before(ctx iris.Context) {
     // [...]
}

func after(ctx iris.Context) {
    // [...]
}

func indexHandler(ctx iris.Context) {
    // write something to the client as a response.
    ctx.HTML("<h1>Index</h1>")

    ctx.Next() // execute the "after" handler registered via `Done`.
}

func contactHandler(ctx iris.Context) {
    // write something to the client as a response.
    ctx.HTML("<h1>Contact</h1>")

    ctx.Next() // execute the "after" handler registered via `Done`.
}
```

你也可以使用 `ExecutionRules` 强制处理器在没有 `ctx.Next()` 的情况下完成执行。你可以这样做：

```go
app.SetExecutionRules(iris.ExecutionRules{
    // Begin: ...
    // Main:  ...
    Done: iris.ExecutionOptions{Force: true},
})
```
