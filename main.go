package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	m := new(sync.Mutex)
	c := sync.NewCond(m)

	wg := new(sync.WaitGroup)          // workerの実行が終わるまで待つ用
	countDownWg := new(sync.WaitGroup) // workerが全員wait終わるまで待つ用
	for i := 0; i < 10; i++ {
		wg.Add(1)
		countDownWg.Add(1)
		go func(i int) {
			fmt.Printf("worker %d wait ...\n", i)
			countDownWg.Done()
			c.L.Lock()
			defer c.L.Unlock()
			defer wg.Done()

			// c.Broadcast()されるまで待つ (Signalだと、順番に再開される)
			c.Wait()

			fmt.Printf("worker %d start\n", i)
		}(i)
	}

	// このwaitが終わると全部のworkerが wait 状態
	countDownWg.Wait()

	// 5秒カウントダウンして
	for i := 5; i > 0; i-- {
		fmt.Printf("%d\n", i)
		time.Sleep(time.Second * 1)
	}

	//------------------------------
	// ここで一気にwait解除させるならBroadcastで、一つづつwait解除するならsignalを実行する

	// // 一斉に実行
	// c.Broadcast()

	// 順番に実行
	for i := 0; i < 10; i++ {
		fmt.Printf("次開始\n")
		c.Signal()
		time.Sleep(time.Millisecond * 500)
	}

	// 全部のworkerの処理が終わるまで待つ
	wg.Wait()

}
