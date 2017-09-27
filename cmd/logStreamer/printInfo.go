package main

import "fmt"

func printLogo() {
	fmt.Println("")
	fmt.Println("  _                   _____ _")
	fmt.Println(" | |                 /  ___| |                                     ")
	fmt.Println(" | |     ___   __ _  \\ `--.| |_ _ __ ___  __ _ _ __ ___   ___ _ __ ")
	fmt.Println(" | |    / _ \\ / _` |  `--. \\ __| '__/ _ \\/ _` | '_ ` _ \\ / _ \\ '__|")
	fmt.Println(" | |___| (_) | (_| | /\\__/ / |_| | |  __/ (_| | | | | | |  __/ |   ")
	fmt.Println(" \\_____/\\___/ \\__, | \\____/ \\__|_|  \\___|\\__,_|_| |_| |_|\\___|_|   ")
	fmt.Println("               __/ |                                               ")
	fmt.Println("              |___/                                                ")
	fmt.Println()
}

func infoBlock(port, pport int, maxSize int64, path string, pprofInfo bool) {
	// info block
	fmt.Println("##########################################################")
	fmt.Println("\t==>\tPort number:\t", port)
	fmt.Println("\t==>\tLog path:\t", path)
	fmt.Println("\t==>\tLog size:\t", maxSize)
	fmt.Println("\t==>\tPprof server:\t", pprofInfo)
	fmt.Println("\t==>\tPprof port:\t", pport)
	fmt.Println("##########################################################")
}
