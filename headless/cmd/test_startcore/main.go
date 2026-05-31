package main

import (
	"fmt"
	"os"

	"flclash-headless/action"
	"flclash-headless/app"
)

func main() {
	a := app.New()
	if err := a.InitStorage(); err != nil {
		fmt.Fprintf(os.Stderr, "FAIL InitStorage: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("InitStorage OK")
	a.StateStore.Get()
	a.StateStore.Save()

	fmt.Println("Starting core...")
	if err := action.StartCore(a); err != nil {
		fmt.Fprintf(os.Stderr, "FAIL StartCore: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("StartCore OK!")

	proxies, err := a.CoreClient.GetProxies()
	if err != nil {
		fmt.Fprintf(os.Stderr, "FAIL GetProxies: %v\n", err)
	} else {
		fmt.Printf("Proxies: %d entries\n", len(proxies.Proxies))
		names := proxies.All
		fmt.Printf("Groups (%d):\n", len(names))
		for _, n := range names {
			fmt.Printf("  - %s\n", n)
		}
	}

	action.StopCore(a)
	fmt.Println("DONE")
}
