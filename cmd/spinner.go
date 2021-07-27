package cmd

import (
	"time"

	spinner "github.com/janeczku/go-spinner"
)

func startBlocksSpinner(title string) (s *spinner.Spinner) {
	s = spinner.NewSpinner(title)
	s.Start()
	return s
}

func incSpinner() {
	time.Sleep(1000 * time.Millisecond)
}
