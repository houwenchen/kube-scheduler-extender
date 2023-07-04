package signals

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

var (
	onlyOneSignalHandler = make(chan struct{})
	shutdownSignals      = []os.Signal{os.Interrupt, syscall.SIGTERM}
)

// 这里两个实现方法类似，及通过一个可以接收2个信号的channel来实现第一次调用ctx的cancel方法，一次调用os退出方法
func SetupSignalHandler() (stopCh <-chan struct{}) {
	close(onlyOneSignalHandler) // 确保此函数只会调用一次，因为关闭的 channel 不能再次关闭

	stop := make(chan struct{})
	c := make(chan os.Signal, 2)
	signal.Notify(c, shutdownSignals...)
	go func() {
		<-c
		close(stop)
		<-c
		os.Exit(1) // second signal. Exit directly.
	}()

	return stop
}

func SetupSignalContext() context.Context {
	close(onlyOneSignalHandler) // 确保此函数只会调用一次，因为关闭的 channel 不能再次关闭

	ctx, cancel := context.WithCancel(context.Background())
	c := make(chan os.Signal, 2)
	signal.Notify(c, shutdownSignals...)
	go func() {
		<-c
		cancel()
		<-c
		os.Exit(1)
	}()

	return ctx
}
