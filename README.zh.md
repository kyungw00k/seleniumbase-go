[English](README.md) | [한국어](README.ko.md) | **中文** | [日本語](README.ja.md)

# seleniumbase-go

[![Go Reference](https://pkg.go.dev/badge/github.com/kyungw00k/seleniumbase-go.svg)](https://pkg.go.dev/github.com/kyungw00k/seleniumbase-go)
[![CI](https://github.com/kyungw00k/seleniumbase-go/actions/workflows/ci.yml/badge.svg)](https://github.com/kyungw00k/seleniumbase-go/actions)
[![Go Version](https://img.shields.io/badge/go-1.22+-blue)](https://go.dev)
[![License](https://img.shields.io/badge/license-MIT-green)](LICENSE)

[SeleniumBase](https://github.com/seleniumbase/SeleniumBase) 的 Go 移植版 — 基于 [playwright-go](https://github.com/playwright-community/playwright-go) 构建的简单、强大的浏览器自动化工具。

---

## 安装

```bash
go get github.com/kyungw00k/seleniumbase-go
```

Playwright 浏览器二进制文件会在首次运行时**自动安装**。也可以手动安装：

```bash
go run github.com/playwright-community/playwright-go/cmd/playwright@v0.5700.1 install --with-deps
```

---

## 快速开始

```go
package main

import (
    "fmt"
    "log"

    "github.com/kyungw00k/seleniumbase-go/sb"
)

func main() {
    err := sb.Run(func(p *sb.Page) error {
        p.Open("https://example.com")

        title, err := p.GetTitle()
        if err != nil {
            return err
        }
        fmt.Println("Page title:", title)

        p.AssertTitle("Example Domain")
        p.AssertElement("h1")
        p.AssertText("Example Domain", "h1")
        return nil
    }, sb.WithBrowser("chromium"), sb.WithHeadless(true))

    if err != nil {
        log.Fatal(err)
    }
}
```

另请参阅：[`sb.RunTest`](docs/features.md#sbruntest--testingt-integration) 用于 `testing.T` 集成，[`sb.NewPage`](docs/features.md#sbnewpage--manual-lifecycle) 用于手动管理生命周期。

---

## 使用模式

seleniumbase-go 提供三种使用浏览器的方式：

| 模式 | 适用场景 | 生命周期 |
|------|----------|----------|
| `sb.Run(func(p *sb.Page) error { ... })` | 脚本与自动化 | 自动 |
| `sb.RunTest(t, func(p *sb.Page) { ... })` | 使用 `testing.T` 的 Go 测试 | 自动，失败会调用 `t.Fatal` |
| `sb.NewPage(opts...)` → `page, cleanup, err` | 完全控制 | 通过 `defer cleanup()` 手动管理 |

---

## 功能特性

- **SeleniumBase 兼容 API** — 沿用 Python 库中熟悉的方法名称
- **选择器转换** — 支持 `id=`、`name=`、`link=`、`css=`、`xpath=` 前缀或裸 CSS（[详情](docs/selectors.md)）
- **自动等待断言** — 基于 Playwright 的重试引擎，提供可靠的测试
- **90+ 页面方法** — 导航、交互、断言、查询、Cookie、存储、JS 执行、滚动、网络等（[完整 API](docs/api.md)）
- **25+ 配置选项** — 浏览器、无头模式、代理、视口、超时、语言环境等（[所有选项](docs/configuration.md)）
- **CDP 隐身模式** — 以反检测标志启动 Chrome（[指南](docs/features.md#stealth-mode)）
- **录制器** — 捕获浏览器交互并生成 Go 测试代码（[指南](docs/features.md#recorder)）
- **可视化测试** — 像素级截图对比（[指南](docs/features.md#visual-testing)）
- **并行测试运行器** — 使用隔离浏览器并发执行（[指南](docs/features.md#parallel-test-runner)）
- **报告生成** — JUnit XML 和样式化 HTML 报告（[指南](docs/features.md#report-generation)）
- **远程浏览器** — 通过 CDP 或 WebSocket 连接到 BrowserStack、LambdaTest（[指南](docs/features.md#remote-browser)）
- **MasterQA** — 自动化与手动验证的混合模式（[指南](docs/features.md#masterqa)）
- **自动浏览器检测** — 查找系统 Chrome 或安装 Playwright Chromium（[指南](docs/features.md#auto-browser-detection)）
- **网络拦截** — 路由、中止和模拟 API 响应
- **高亮 / 演示模式** — 用于调试的元素可视化高亮
- **延迟断言** — 将失败入队并统一处理
- **MFA / TOTP 支持** — 生成和输入基于时间的一次性密码
- **Cookie 持久化** — 保存和加载浏览器状态
- **标签页管理、框架、对话框** — 完整的浏览器控制
- **逃生通道** — 随时访问底层 `playwright.Page`、`playwright.Locator`、`playwright.BrowserContext`

---

## 文档

| 文档 | 说明 |
|------|------|
| [API 参考](docs/api.md) | 所有 90+ 页面方法及其签名 |
| [配置](docs/configuration.md) | 25+ 选项和超时常量 |
| [选择器语法](docs/selectors.md) | SeleniumBase 选择器转换规则 |
| [功能指南](docs/features.md) | 隐身模式、录制器、可视化测试等详细指南 |

---

## 示例

| 示例 | 说明 |
|------|------|
| [`examples/simple/main.go`](examples/simple/main.go) | `sb.Run` — 标题和断言检查 |
| [`examples/basic_test.go`](examples/basic_test.go) | `sb.RunTest` — 登录/购物车/退出登录流程 |
| [`examples/demo_site_test.go`](examples/demo_site_test.go) | 带高亮的演示模式 |
| [`examples/visual_test.go`](examples/visual_test.go) | 可视化回归测试 |
| [`examples/stealth_coupang_test.go`](examples/stealth_coupang_test.go) | CDP 隐身模式 |

运行集成测试：

```bash
go test -tags integration ./examples/...
```

---

## 国际化 (i18n)

SeleniumBase Python 通过静态子类化支持 10 种语言的[方法名翻译](https://github.com/seleniumbase/SeleniumBase/tree/master/seleniumbase/translate)（例如韩语中的 `self.클릭()`）。Go 移植版**不支持**此功能 — Go 的方法名在编译时解析，使得运行时或基于别名的翻译不切实际。

---

## 灵感来源

本项目的灵感来源于 [SeleniumBase](https://github.com/seleniumbase/SeleniumBase)，这是一个功能全面的 Python 浏览器自动化框架。seleniumbase-go 以 [playwright-go](https://github.com/playwright-community/playwright-go) 为底层引擎，将同样简洁、强大的 API 设计带入 Go 生态系统。

---

## 许可证

MIT。详见 [LICENSE](LICENSE)。
