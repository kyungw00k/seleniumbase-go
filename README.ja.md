[English](README.md) | [한국어](README.ko.md) | [中文](README.zh.md) | **日本語**

# seleniumbase-go

[![Go Reference](https://pkg.go.dev/badge/github.com/kyungw00k/seleniumbase-go.svg)](https://pkg.go.dev/github.com/kyungw00k/seleniumbase-go)
[![CI](https://github.com/kyungw00k/seleniumbase-go/actions/workflows/ci.yml/badge.svg)](https://github.com/kyungw00k/seleniumbase-go/actions)
[![Go Version](https://img.shields.io/badge/go-1.22+-blue)](https://go.dev)
[![License](https://img.shields.io/badge/license-MIT-green)](LICENSE)

[SeleniumBase](https://github.com/seleniumbase/SeleniumBase) の Go 移植版 — [playwright-go](https://github.com/playwright-community/playwright-go) を基盤とした、シンプルで強力なブラウザ自動化ライブラリです。

---

## インストール

```bash
go get github.com/kyungw00k/seleniumbase-go
```

Playwright のブラウザバイナリは初回実行時に**自動インストール**されます。手動でインストールする場合は以下を実行してください:

```bash
go run github.com/playwright-community/playwright-go/cmd/playwright@v0.5700.1 install --with-deps
```

---

## クイックスタート

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

関連項目: `testing.T` との統合については [`sb.RunTest`](docs/features.md#sbruntest--testingt-integration)、手動ライフサイクル制御については [`sb.NewPage`](docs/features.md#sbnewpage--manual-lifecycle) を参照してください。

---

## 使用パターン

seleniumbase-go はブラウザを使用する 3 つの方法を提供します:

| パターン | ユースケース | ライフサイクル |
|---------|------------|--------------|
| `sb.Run(func(p *sb.Page) error { ... })` | スクリプト・自動化 | 自動 |
| `sb.RunTest(t, func(p *sb.Page) { ... })` | `testing.T` を使った Go テスト | 自動、失敗は `t.Fatal` へ |
| `sb.NewPage(opts...)` → `page, cleanup, err` | 完全な制御 | `defer cleanup()` で手動管理 |

---

## 機能

- **SeleniumBase 互換 API** — Python ライブラリと同じメソッド名
- **セレクター変換** — `id=`、`name=`、`link=`、`css=`、`xpath=` プレフィックスまたは素の CSS を使用可能 ([詳細](docs/selectors.md))
- **自動待機アサーション** — Playwright のリトライエンジンによる安定したテスト
- **90 以上の Page メソッド** — ナビゲーション、操作、アサーション、クエリ、Cookie、ストレージ、JS 評価、スクロール、ネットワークなど ([全 API](docs/api.md))
- **25 以上の設定オプション** — ブラウザ、ヘッドレス、プロキシ、ビューポート、タイムアウト、ロケールなど ([全オプション](docs/configuration.md))
- **CDP ステルスモード** — アンチ検出フラグで Chrome を起動 ([ガイド](docs/features.md#stealth-mode))
- **Recorder** — ブラウザ操作をキャプチャして Go テストコードを生成 ([ガイド](docs/features.md#recorder))
- **ビジュアルテスト** — ピクセルレベルのスクリーンショット比較 ([ガイド](docs/features.md#visual-testing))
- **並列テストランナー** — 独立したブラウザで並行実行 ([ガイド](docs/features.md#parallel-test-runner))
- **レポート生成** — JUnit XML およびスタイル付き HTML レポート ([ガイド](docs/features.md#report-generation))
- **リモートブラウザ** — CDP または WebSocket で BrowserStack、LambdaTest に接続 ([ガイド](docs/features.md#remote-browser))
- **MasterQA** — 自動化と手動検証のハイブリッドモード ([ガイド](docs/features.md#masterqa))
- **ブラウザ自動検出** — システムの Chrome を検索するか、Playwright Chromium をインストール ([ガイド](docs/features.md#auto-browser-detection))
- **ネットワークインターセプト** — API レスポンスのルーティング、中断、モック
- **ハイライト / デモモード** — デバッグ用のビジュアル要素ハイライト
- **遅延アサーション** — 失敗をキューに積んでまとめて処理
- **MFA / TOTP サポート** — 時間ベースのワンタイムパスワードの生成と入力
- **Cookie 永続化** — ブラウザ状態の保存と読み込み
- **タブ管理、フレーム、ダイアログ** — フルブラウザ制御
- **エスケープハッチ** — いつでも `playwright.Page`、`playwright.Locator`、`playwright.BrowserContext` に直接アクセス可能

---

## ドキュメント

| ドキュメント | 説明 |
|------------|------|
| [API リファレンス](docs/api.md) | 90 以上の Page メソッドとシグネチャ一覧 |
| [設定](docs/configuration.md) | 25 以上のオプションとタイムアウト定数 |
| [セレクター構文](docs/selectors.md) | SeleniumBase セレクター変換ルール |
| [機能ガイド](docs/features.md) | ステルス、Recorder、ビジュアルテストなどの詳細ガイド |

---

## サンプル

| サンプル | 説明 |
|---------|------|
| [`examples/simple/main.go`](examples/simple/main.go) | `sb.Run` — タイトルとアサーションの確認 |
| [`examples/basic_test.go`](examples/basic_test.go) | `sb.RunTest` — ログイン・カート・ログアウトのフロー |
| [`examples/demo_site_test.go`](examples/demo_site_test.go) | ハイライト付きデモモード |
| [`examples/visual_test.go`](examples/visual_test.go) | ビジュアルリグレッションテスト |
| [`examples/stealth_coupang_test.go`](examples/stealth_coupang_test.go) | CDP ステルスモード |

インテグレーションテストの実行:

```bash
go test -tags integration ./examples/...
```

---

## 国際化 (i18n)

SeleniumBase Python は、静的サブクラス化によって 10 言語のメソッド名翻訳をサポートしています (例: 韓国語での `self.클릭()`)。これは Go 移植版では**サポートされていません** — Go のメソッド名はコンパイル時に解決されるため、ランタイムやエイリアスベースの翻訳は実用的ではありません。

---

## インスピレーション

このプロジェクトは、包括的な Python ブラウザ自動化フレームワークである [SeleniumBase](https://github.com/seleniumbase/SeleniumBase) にインスパイアされています。seleniumbase-go は、[playwright-go](https://github.com/playwright-community/playwright-go) を基盤エンジンとして使用し、同じシンプルで強力な API 設計を Go エコシステムにもたらします。

---

## ライセンス

MIT。詳細は [LICENSE](LICENSE) を参照してください。
