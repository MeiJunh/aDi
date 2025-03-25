package util

import "aDi/log"

func GoSafe(fn func()) {
	go runSafe(fn)
}

func runSafe(fn func()) {
	defer func() {
		errP := recover()
		PanicDeal(errP, func(panic string) { log.Errorf("panic:%s", panic) })
	}()
	fn()
}
