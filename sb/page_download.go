package sb

// WaitForDownload triggers the given action and waits for a download to start.
// Returns the downloaded file path.
func (p *Page) WaitForDownload(action func() error) (string, error) {
	download, err := p.pw.ExpectDownload(func() error {
		return action()
	})
	if err != nil {
		return "", err
	}
	path, err := download.Path()
	if err != nil {
		return "", err
	}
	return path, nil
}

// SaveDownload triggers the given action, waits for download, and saves to the specified path.
func (p *Page) SaveDownload(path string, action func() error) error {
	download, err := p.pw.ExpectDownload(func() error {
		return action()
	})
	if err != nil {
		return err
	}
	return download.SaveAs(path)
}
