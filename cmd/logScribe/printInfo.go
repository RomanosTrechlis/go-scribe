package main

import "fmt"

func printLogo() {
	fmt.Println("	██╗      ██████╗  ██████╗     ███████╗ ██████╗██████╗ ██╗██████╗ ███████╗")
	fmt.Println("	██║     ██╔═══██╗██╔════╝     ██╔════╝██╔════╝██╔══██╗██║██╔══██╗██╔════╝")
	fmt.Println("	██║     ██║   ██║██║  ███╗    ███████╗██║     ██████╔╝██║██████╔╝█████╗  ")
	fmt.Println("	██║     ██║   ██║██║   ██║    ╚════██║██║     ██╔══██╗██║██╔══██╗██╔══╝  ")
	fmt.Println("	███████╗╚██████╔╝╚██████╔╝    ███████║╚██████╗██║  ██║██║██████╔╝███████╗")
	fmt.Println("	╚══════╝ ╚═════╝  ╚═════╝     ╚══════╝ ╚═════╝╚═╝  ╚═╝╚═╝╚═════╝ ╚══════╝")
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
