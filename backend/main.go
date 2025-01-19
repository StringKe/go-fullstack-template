package main

// 总入口，应用是一个 Cli 需要通过不同的参数启动不同的功能
func main() {
	err := RootCmd.Execute()
	if err != nil {
		panic(err)
	}
}
