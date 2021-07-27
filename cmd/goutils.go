package cmd

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

const (
	GoImports string = "goimports"
)

func goExec(mycmd string) (err error) {
	cmd := exec.Command(mycmd)
	fmt.Println(cmd)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Start()
	if err != nil {
		return err
	}

	err = cmd.Wait()
	if err != nil {
		log.Printf("Error while running go imports. Error: %v\n", err)
		return err
	}

	return nil
}
