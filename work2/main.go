package main

import (
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"
	"net/http"
	"os"
	"os/signal"
)

//1. 基于 errgroup 实现一个 http server 的启动和关闭 ，以及 linux signal 信号的注册和处理，要保证能够一个退出，全部注销退出。
func main() {
	ctx := context.Background()
	// 定义 withCancel -> cancel() 方法 去取消下游的 Context
	ctx, cancel := context.WithCancel(ctx)
	// 使用 errgroup 进行 goroutine 取消
	group, errCtx := errgroup.WithContext(ctx)
	//http server
	server := &http.Server{Addr: ":8080"}

	group.Go(func() error {
		return StartHttpServer(server)
	})

	group.Go(func() error {
		fmt.Println("blocking befor close server....")
		<-errCtx.Done()
		fmt.Println(" closing server...")
		// 关闭 TestServer
		return server.Shutdown(errCtx)
	})

	// 监听 linux signal 信号
	chanel := make(chan os.Signal, 1)
	signal.Notify(chanel)

	group.Go(func() error {
		for {
			select {
			case <-errCtx.Done():
				fmt.Println("异常关闭....")
				return errCtx.Err()
			case <-chanel:
				// 终止 程序
				fmt.Println("监听到终止信号...")
				cancel()
			}
		}
		return nil
	})

	// 阻塞等待
	if err := group.Wait(); err != nil {
		fmt.Println("group  error: ", err)
	}
	fmt.Println("all group done!")

}
func StartHttpServer(server *http.Server) error {
	http.HandleFunc("/hello", TestServer)
	fmt.Println("test server start...")
	err := server.ListenAndServe()
	return err
}

func TestServer(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "hello world")
}
