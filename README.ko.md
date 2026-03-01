[English](README.md) | **한국어** | [中文](README.zh.md) | [日本語](README.ja.md)

# seleniumbase-go

[![Go Reference](https://pkg.go.dev/badge/github.com/kyungw00k/seleniumbase-go.svg)](https://pkg.go.dev/github.com/kyungw00k/seleniumbase-go)
[![CI](https://github.com/kyungw00k/seleniumbase-go/actions/workflows/ci.yml/badge.svg)](https://github.com/kyungw00k/seleniumbase-go/actions)
[![Go Version](https://img.shields.io/badge/go-1.22+-blue)](https://go.dev)
[![License](https://img.shields.io/badge/license-MIT-green)](LICENSE)

[SeleniumBase](https://github.com/seleniumbase/SeleniumBase)의 Go 포트 — [playwright-go](https://github.com/playwright-community/playwright-go) 위에서 동작하는 간단하고 강력한 브라우저 자동화 라이브러리입니다.

---

## 설치

```bash
go get github.com/kyungw00k/seleniumbase-go
```

Playwright 브라우저 바이너리는 첫 실행 시 **자동으로 설치**됩니다. 수동으로 설치하려면 다음 명령을 실행하세요:

```bash
go run github.com/playwright-community/playwright-go/cmd/playwright@v0.5700.1 install --with-deps
```

---

## 빠른 시작

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

`testing.T` 통합을 위해서는 [`sb.RunTest`](docs/features.md#sbruntest--testingt-integration)를, 라이프사이클을 직접 제어하려면 [`sb.NewPage`](docs/features.md#sbnewpage--manual-lifecycle)를 참고하세요.

---

## 사용 패턴

seleniumbase-go는 브라우저를 사용하는 세 가지 방식을 제공합니다:

| 패턴 | 사용 목적 | 라이프사이클 |
|------|----------|-------------|
| `sb.Run(func(p *sb.Page) error { ... })` | 스크립트 및 자동화 | 자동 |
| `sb.RunTest(t, func(p *sb.Page) { ... })` | `testing.T`를 사용하는 Go 테스트 | 자동, 실패 시 `t.Fatal` 호출 |
| `sb.NewPage(opts...)` → `page, cleanup, err` | 전체 직접 제어 | `defer cleanup()`으로 수동 관리 |

---

## 기능

- **SeleniumBase 호환 API** — Python 라이브러리의 친숙한 메서드명 사용
- **셀렉터 변환** — `id=`, `name=`, `link=`, `css=`, `xpath=` 접두사 또는 일반 CSS 사용 가능 ([상세 정보](docs/selectors.md))
- **자동 대기 어설션** — Playwright의 재시도 엔진을 활용한 안정적인 테스트
- **90개 이상의 Page 메서드** — 내비게이션, 인터랙션, 어설션, 쿼리, 쿠키, 스토리지, JS 실행, 스크롤, 네트워크 등 ([전체 API](docs/api.md))
- **25개 이상의 설정 옵션** — 브라우저, 헤드리스 모드, 프록시, 뷰포트, 타임아웃, 로케일 등 ([전체 옵션](docs/configuration.md))
- **CDP 스텔스 모드** — 감지 방지 플래그로 Chrome 실행 ([가이드](docs/features.md#stealth-mode))
- **Recorder** — 브라우저 인터랙션을 기록하여 Go 테스트 코드 생성 ([가이드](docs/features.md#recorder))
- **비주얼 테스팅** — 픽셀 수준의 스크린샷 비교 ([가이드](docs/features.md#visual-testing))
- **병렬 테스트 실행기** — 격리된 브라우저로 동시 실행 ([가이드](docs/features.md#parallel-test-runner))
- **리포트 생성** — JUnit XML 및 스타일링된 HTML 리포트 ([가이드](docs/features.md#report-generation))
- **원격 브라우저** — CDP 또는 WebSocket을 통해 BrowserStack, LambdaTest 연결 ([가이드](docs/features.md#remote-browser))
- **MasterQA** — 자동화와 수동 검증을 결합한 하이브리드 방식 ([가이드](docs/features.md#masterqa))
- **자동 브라우저 감지** — 시스템 Chrome을 찾거나 Playwright Chromium을 자동 설치 ([가이드](docs/features.md#auto-browser-detection))
- **네트워크 인터셉션** — API 응답을 라우팅, 차단, 모킹
- **하이라이트 / 데모 모드** — 디버깅을 위한 시각적 요소 강조
- **지연 어설션** — 실패를 큐에 쌓아 일괄 처리
- **MFA / TOTP 지원** — 시간 기반 일회용 비밀번호 생성 및 입력
- **쿠키 유지** — 브라우저 상태 저장 및 불러오기
- **탭 관리, 프레임, 다이얼로그** — 완전한 브라우저 제어
- **이스케이프 해치** — 언제든지 `playwright.Page`, `playwright.Locator`, `playwright.BrowserContext`로 직접 접근 가능

---

## 문서

| 문서 | 설명 |
|------|------|
| [API 레퍼런스](docs/api.md) | 90개 이상의 Page 메서드와 시그니처 |
| [설정](docs/configuration.md) | 25개 이상의 옵션 및 타임아웃 상수 |
| [셀렉터 문법](docs/selectors.md) | SeleniumBase 셀렉터 변환 규칙 |
| [기능 가이드](docs/features.md) | 스텔스 모드, Recorder, 비주얼 테스팅 등 상세 가이드 |

---

## 예시

| 예시 | 설명 |
|------|------|
| [`examples/simple/main.go`](examples/simple/main.go) | `sb.Run` — 타이틀 및 어설션 확인 |
| [`examples/basic_test.go`](examples/basic_test.go) | `sb.RunTest` — 로그인/장바구니/로그아웃 흐름 |
| [`examples/demo_site_test.go`](examples/demo_site_test.go) | 하이라이트를 포함한 데모 모드 |
| [`examples/visual_test.go`](examples/visual_test.go) | 비주얼 회귀 테스트 |
| [`examples/stealth_coupang_test.go`](examples/stealth_coupang_test.go) | CDP 스텔스 모드 |

통합 테스트 실행:

```bash
go test -tags integration ./examples/...
```

---

## 국제화 (i18n)

SeleniumBase Python은 정적 서브클래싱을 통해 10개 언어로 메서드 이름을 [번역하는 기능](https://github.com/seleniumbase/SeleniumBase/tree/master/seleniumbase/translate)을 지원합니다 (예: 한국어의 `self.클릭()`). 이 기능은 Go 포트에서는 **지원되지 않습니다** — Go의 메서드 이름은 컴파일 타임에 결정되므로, 런타임 또는 별칭 기반의 번역이 실용적이지 않습니다.

---

## 영감을 받은 프로젝트

이 프로젝트는 포괄적인 Python 브라우저 자동화 프레임워크인 [SeleniumBase](https://github.com/seleniumbase/SeleniumBase)에서 영감을 받았습니다. seleniumbase-go는 [playwright-go](https://github.com/playwright-community/playwright-go)를 기반 엔진으로 사용하여, 동일한 단순하고 강력한 API 설계를 Go 생태계에 가져옵니다.

---

## 라이선스

MIT. [LICENSE](LICENSE) 파일을 참고하세요.
