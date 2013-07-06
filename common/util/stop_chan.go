/**
 * Date: 05.07.13
 * Time: 22:18
 *
 * @author Vladimir Matveev
 */
package util

// StopChan represents a channel used to tell to some process that it should stop its work.
type StopChan chan interface{}

// Stopped checks whether the given StopChan received a request to stop.
func (ch StopChan) Stopped() bool {
    select {
    case _ = <-ch:
        return true
    default:
        // Still not stopped
    }
    return false
}

// Wait waits until a stop request will have been sent to the given StopChan.
func (ch StopChan) Wait() {
    _ = <-ch
}

// Stop sends a stop request to the given StopChan.
func (ch StopChan) Stop() {
    select {
    case ch <- nil:
    default:
    }
    close(ch)
}
