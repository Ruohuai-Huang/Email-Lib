### 安装
>  go get github.com/jordan-wright/email

1.实现简单的邮件发送：

package main

import (
    "log"
    "net/smtp"

    "github.com/jordan-wright/email"
)

func main() {
    e := email.NewEmail()
    //设置发送方的邮箱
    e.From = "dj <XXX@qq.com>"
    // 设置接收方的邮箱
    e.To = []string{"XXX@qq.com"}
    //设置主题
    e.Subject = "这是主题"
    //设置文件发送的内容
    e.Text = []byte("www.topgoer.com是个不错的go语言中文文档")
    //设置服务器相关的配置
    err := e.Send("smtp.qq.com:25", smtp.PlainAuth("", "你的邮箱账号", "这块是你的授权码", "smtp.qq.com"))
    if err != nil {
        log.Fatal(err)
    }
}
运行程序就会给你设置的邮箱发送一个邮件，有的邮箱会把邮件当成垃圾邮件发到垃圾箱里面，如果找不到邮件可以去垃圾箱看下。

2.实现抄送功能

该插件有两种抄送模式即 CC（Carbon Copy）和 BCC （Blind Carbon Copy）

抄送功能只需要添加两个参数就好了

    e.Cc = []string{"XXX@qq.com",XXX@qq.com}
    e.Bcc = []string{"XXX@qq.com"}
全部代码：

package main

import (
    "log"
    "net/smtp"

    "github.com/jordan-wright/email"
)

func main() {
    e := email.NewEmail()
    //设置发送方的邮箱
    e.From = "dj <XXX@qq.com>"
    // 设置接收方的邮箱
    e.To = []string{"XXX@qq.com"}
    //设置抄送如果抄送多人逗号隔开
    e.Cc = []string{"XXX@qq.com",XXX@qq.com}
    //设置秘密抄送
    e.Bcc = []string{"XXX@qq.com"}
    //设置主题
    e.Subject = "这是主题"
    //设置文件发送的内容
    e.Text = []byte("www.topgoer.com是个不错的go语言中文文档")
    //设置服务器相关的配置
    err := e.Send("smtp.qq.com:25", smtp.PlainAuth("", "你的邮箱账号", "这块是你的授权码", "smtp.qq.com"))
    if err != nil {
        log.Fatal(err)
    }
}
3.发送html代码的邮件

代码实现：

package main

import (
    "log"
    "net/smtp"

    "github.com/jordan-wright/email"
)

func main() {
    e := email.NewEmail()
    //设置发送方的邮箱
    e.From = "dj <XXX@qq.com>"
    // 设置接收方的邮箱
    e.To = []string{"XXX@qq.com"}
    //设置主题
    e.Subject = "这是主题"
    //设置文件发送的内容
    e.HTML = []byte(`
    <h1><a href="http://www.topgoer.com/">go语言中文网站</a></h1>    
    `)
    //设置服务器相关的配置
    err := e.Send("smtp.qq.com:25", smtp.PlainAuth("", "你的邮箱账号", "这块是你的授权码", "smtp.qq.com"))
    if err != nil {
        log.Fatal(err)
    }
}
4.实现邮件附件的发送

直接调用AttachFile即可

package main

import (
    "log"
    "net/smtp"

    "github.com/jordan-wright/email"
)

func main() {
    e := email.NewEmail()
    //设置发送方的邮箱
    e.From = "dj <XXX@qq.com>"
    // 设置接收方的邮箱
    e.To = []string{"XXX@qq.com"}
    //设置主题
    e.Subject = "这是主题"
    //设置文件发送的内容
    e.HTML = []byte(`
    <h1><a href="http://www.topgoer.com/">go语言中文网站</a></h1>    
    `)
    //这块是设置附件
    e.AttachFile("./test.txt")
    //设置服务器相关的配置
    err := e.Send("smtp.qq.com:25", smtp.PlainAuth("", "你的邮箱账号", "这块是你的授权码", "smtp.qq.com"))
    if err != nil {
        log.Fatal(err)
    }
}
5.连接池

实际上每次调用Send时都会和 SMTP 服务器建立一次连接，如果发送邮件很多很频繁的话可能会有性能问题。email提供了连接池，可以复用网络连接：

package main

import (
    "fmt"
    "log"
    "net/smtp"
    "os"
    "sync"
    "time"

    "github.com/jordan-wright/email"
)

func main() {
    ch := make(chan *email.Email, 10)
    p, err := email.NewPool(
        "smtp.qq.com:25",
        4,
        smtp.PlainAuth("", "XXX@qq.com", "你的授权码", "smtp.qq.com"),
    )

    if err != nil {
        log.Fatal("failed to create pool:", err)
    }

    var wg sync.WaitGroup
    wg.Add(4)
    for i := 0; i < 4; i++ {
        go func() {
            defer wg.Done()
            for e := range ch {
                err := p.Send(e, 10*time.Second)
                if err != nil {
                    fmt.Fprintf(os.Stderr, "email:%v sent error:%v\n", e, err)
                }
            }
        }()
    }

    for i := 0; i < 10; i++ {
        e := email.NewEmail()
        e.From = "dj <XXX@qq.com>"
        e.To = []string{"XXX@qq.com"}
        e.Subject = "Awesome web"
        e.Text = []byte(fmt.Sprintf("Awesome Web %d", i+1))
        ch <- e
    }

    close(ch)
    wg.Wait()
}
上面程序中，我们创建 4 goroutine 共用一个连接池发送邮件，发送 10 封邮件后程序退出。为了等邮件都发送完成或失败，程序才退出，我们使用了sync.WaitGroup。由于使用了 goroutine，邮件顺序不能保证。

参考：https://github.com/darjun/go-daily-lib