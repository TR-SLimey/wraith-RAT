// +build debug

/*
This is a debug receiver for testing purposes. It only gets included in debug builds.
When included, it "receives" a print command every 2 seconds.
*/

package rx

import "time"

var Debug Rx

func init() {
	// Register handler for the debug:// URL scheme (which is never really used)
	RxList.Add("debug", &Debug)
	// Create a channel to trigger exit via the `Stop` method
	Debug.Data["exitTrigger"] = make(chan struct{})
	// On start, run a thread pushing a debug message every 2 seconds
	Debug.Start = func() {
		go func() {
			defer close(Debug.Data["exitTrigger"].(chan struct{}))
			for {
				select {
				case <-Debug.Data["exitTrigger"].(chan struct{}):
					return
				case <-time.After(2 * time.Second):
					Debug.Data["queue"].(RxQueue) <- RxQueueElement{Data: []byte{}} /*RxQueueElement{Data: map[string]interface{}{
						"w.cmd": `func wcmd() string {println("Message from debug receiver"); return ""}`,
					}}*/
				}
			}
		}()
	}
	// On stop
	Debug.Stop = func() {
		// Trigger exit
		Debug.Data["exitTrigger"].(chan struct{}) <- struct{}{}
		// Wait until channel closed (exit confirmed)
		<-Debug.Data["exitTrigger"].(chan struct{})
	}
}
