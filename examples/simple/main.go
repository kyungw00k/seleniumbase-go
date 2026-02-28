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

		url := p.GetCurrentURL()
		fmt.Println("Current URL:", url)

		p.AssertTitle("Example Domain")
		p.AssertElement("h1")
		p.AssertText("Example Domain", "h1")

		fmt.Println("All assertions passed!")
		return nil
	}, sb.WithBrowser("chromium"), sb.WithHeadless(true))

	if err != nil {
		log.Fatal(err)
	}
}
