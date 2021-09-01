/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package grmanager

import (
	"context"
	"io"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/glory-go/glory/log"
)

var shutdownSignals = []os.Signal{
	syscall.SIGHUP,
	syscall.SIGINT,
	syscall.SIGQUIT,
	syscall.SIGILL,
	syscall.SIGTRAP,
	syscall.SIGABRT,
	syscall.SIGTERM,
	os.Interrupt,
	os.Kill,
	syscall.SIGKILL,
}

var (
	cancelPool []context.CancelFunc
	closers    []io.Closer
)

func init() {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, shutdownSignals...)
	go func() {
		for {
			select {
			case sig := <-signals:
				log.Infof("Got Interrupt signal %+v\n", sig)
				for _, v := range cancelPool {
					v()
				}
				for _, closer := range closers {
					go closer.Close()
				}
				time.Sleep(time.Second)
				os.Exit(0)
			}
		}
	}()
}

func RegisterCloser(c io.Closer) {
	closers = append(closers, c)
}

func NewCtx() context.Context {
	ctx, cancel := context.WithCancel(context.Background())
	cancelPool = append(cancelPool, cancel)
	return ctx
}
