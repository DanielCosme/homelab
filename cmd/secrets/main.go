package main

import (
	"fmt"
	"os"

	"danicos.dev/daniel/homelab/pkg/root"
	"danicos.dev/daniel/homelab/pkg/secrets"
)

func main() {
	secrets_stack := secrets.Stack()
	err := secrets_stack.MarshalYaml(root.TMP_FOLDER)
	assertNoErr(err)
}

func assertNoErr(err error) {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
