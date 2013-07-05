/**
 * Date: 15.01.13
 * Time: 2:02
 *
 * @author Vladimir Matveev
 */
package peers

import (
    . "github.com/dpx-infinity/bridge-server/common"
    "log"
    "runtime"
    "sync"
)

type ChanAggregator struct {
    sinks       []SinkChan
    sources     []SourceChan
    sink        Chan
    source      Chan
    lock        sync.Mutex
    interrupted bool
}

func NewAggregator() *ChanAggregator {
    return &ChanAggregator{
        sink:        make(SinkChan),
        source:      make(SourceChan),
        interrupted: true,
    }
}

func (a *ChanAggregator) AddSink(sink SinkChan) {
    a.lock.Lock()
    defer a.lock.Unlock()

    a.sinks = append(a.sinks, sink)
}

func (a *ChanAggregator) AddSource(source SourceChan) {
    a.lock.Lock()
    defer a.lock.Unlock()

    a.sources = append(a.sources, source)
}

func (a *ChanAggregator) OwnSink() SinkChan {
    return a.sink
}

func (a *ChanAggregator) OwnSource() SourceChan {
    return a.source
}

func (a *ChanAggregator) Start() {
    if !a.interrupted {
        return
    }

    a.interrupted = false

    // Start sink loop
    go func() {
        for {
            msg, ok := <-a.sink
            if !ok {
                // TODO: handle error somehow
                log.Println("Sink chan has been closed")
                return
            }

            a.lock.Lock()
            for _, c := range a.sinks {
                c <- msg
            }
            a.lock.Unlock()

            if a.interrupted {
                return
            }
        }
    }()

    // Start source loop
    go func() {
        for {
            a.lock.Lock()
            for _, c := range a.sources {
                select {
                case msg, ok := <-c:
                    if ok {
                        a.source <- msg
                    } else {
                        // TODO: handle error somehow, or maybe remove closed channel
                        log.Println("One of source chans has been closed")
                        return
                    }

                default:
                    // Do nothing for now
                    runtime.Gosched()
                }
            }
            a.lock.Unlock()

            if a.interrupted {
                return
            }

            runtime.Gosched()
        }
    }()
}

func (a *ChanAggregator) Stop() {
    a.interrupted = true
    close(a.sink)
    // We don't close source channel
}
