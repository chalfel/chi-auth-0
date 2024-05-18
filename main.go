package main

import "github.com/chalfel/chi-auth-0/cmd"

func main() {
	rootCmd := cmd.NewRootCmd()

	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}
