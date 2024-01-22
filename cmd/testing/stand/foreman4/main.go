package main

import (
	"context"
	"fmt"
	"github.com/ischenkx/kantoku/cmd/testing/stand/common"
	"github.com/ischenkx/kantoku/pkg/common/broker"
	"github.com/ischenkx/kantoku/pkg/core/event"
	"time"
)

func main() {
	common.InitLogger()
	ctx := context.Background()
	sys := common.NewSystem(ctx, "foreman-0")

	channel1, err := sys.Events().Consume(ctx, broker.TopicsInfo{
		Group:  "group1",
		Topics: []string{"a", "b"},
	})
	if err != nil {
		fmt.Println("failed:", err)
		return
	}

	channel2, err := sys.Events().Consume(ctx, broker.TopicsInfo{
		Group:  "group2",
		Topics: []string{"a", "b"},
	})
	if err != nil {
		fmt.Println("failed:", err)
		return
	}

	go func() {
		i := 1
		for mes := range channel1 {
			fmt.Println(i, "group1w1", "item:", mes.Item())
			i++
			mes.Ack()
		}
	}()

	go func() {
		i := 1
		for mes := range channel1 {
			fmt.Println(i, "group1w2", "item:", mes.Item())
			i++
			mes.Ack()
		}
	}()

	go func() {
		i := 1
		for mes := range channel2 {
			//fmt.Println(i, "group2", "item:", mes.Item())
			i++
			mes.Ack()
		}
	}()

	for i := 0; i < 20; i++ {
		if err := sys.Events().Send(ctx, event.New("a", []byte("hello"))); err != nil {
			fmt.Println("failed to send:", err)
		} else {
			fmt.Println("sent!")
		}
	}

	time.Sleep(time.Minute)

	//list, err := recutil.List(ctx, sys.
	//	Tasks().
	//	Filter(record.R{"inputs": ops.Contains[string]("8f46a98185394216bb1f94bbab325bd1")}).
	//	Cursor().
	//	Iter(),
	//)
	//
	//if err != nil {
	//	fmt.Println("failed to list:", err)
	//	return
	//}
	//
	//fmt.Println("total:", len(list))
	//
	//for _, t := range list {
	//	fmt.Println("-", t.ID)
	//}
}
