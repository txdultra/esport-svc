package periodic

import (
    "io"
    "log"
    "sync"
    "strconv"
    "bytes"
    "github.com/Lupino/periodic/driver"
    "github.com/Lupino/periodic/protocol"
)


type worker struct {
    jobQueue map[int64]driver.Job
    conn     protocol.Conn
    sched    *Sched
    alive    bool
    funcs    []string
    locker   *sync.Mutex
}


func newWorker(sched *Sched, conn protocol.Conn) (w *worker) {
    w = new(worker)
    w.conn = conn
    w.jobQueue = make(map[int64]driver.Job)
    w.sched = sched
    w.funcs = make([]string, 0)
    w.alive = true
    w.locker = new(sync.Mutex)
    return
}


func (w *worker) IsAlive() bool {
    return w.alive
}


func (w *worker) handleJobAssign(msgId []byte, job driver.Job) (err error){
    defer w.locker.Unlock()
    w.locker.Lock()
    w.jobQueue[job.Id()] = job
    buf := bytes.NewBuffer(nil)
    buf.Write(msgId)
    buf.Write(protocol.NULL_CHAR)
    buf.Write(protocol.JOB_ASSIGN.Bytes())
    buf.Write(protocol.NULL_CHAR)
    buf.WriteString(strconv.FormatInt(job.Id(), 10))
    buf.Write(protocol.NULL_CHAR)
    job.Clone().Segment.WriteTo(buf)
    err = w.conn.Send(buf.Bytes())
    return
}


func (w *worker) handleCanDo(Func string) error {
    for _, f := range w.funcs {
        if f == Func {
            return nil
        }
    }
    w.funcs = append(w.funcs, Func)
    w.sched.incrStatFunc(Func)
    return nil
}


func (w *worker) handleCanNoDo(Func string) error {
    newFuncs := make([]string, 0)
    for _, f := range w.funcs {
        if f == Func {
            continue
        }
        newFuncs = append(newFuncs, f)
    }
    w.funcs = newFuncs
    return nil
}


func (w *worker) handleDone(jobId int64) (err error) {
    w.sched.done(jobId)
    defer w.locker.Unlock()
    w.locker.Lock()
    if _, ok := w.jobQueue[jobId]; ok {
        delete(w.jobQueue, jobId)
    }
    return nil
}


func (w *worker) handleFail(jobId int64) (err error) {
    w.sched.fail(jobId)
    defer w.locker.Unlock()
    w.locker.Lock()
    if _, ok := w.jobQueue[jobId]; ok {
        delete(w.jobQueue, jobId)
    }
    return nil
}


func (w *worker) handleCommand(msgId []byte, cmd protocol.Command) (err error) {
    buf := bytes.NewBuffer(nil)
    buf.Write(msgId)
    buf.Write(protocol.NULL_CHAR)
    buf.Write(cmd.Bytes())
    err = w.conn.Send(buf.Bytes())
    return
}


func (w *worker) handleSchedLater(jobId, delay int64) (err error){
    w.sched.schedLater(jobId, delay)
    defer w.locker.Unlock()
    w.locker.Lock()
    if _, ok := w.jobQueue[jobId]; ok {
        delete(w.jobQueue, jobId)
    }
    return nil
}


func (w *worker) handleGrabJob(msgId []byte) (err error){
    item := grabItem{
        w: w,
        msgId: msgId,
    }
    w.sched.grabQueue.push(item)
    w.sched.notifyJobTimer()
    return nil
}


func (w *worker) handle() {
    var payload []byte
    var err error
    var conn = w.conn
    var msgId []byte
    var cmd protocol.Command
    defer func() {
        if x := recover(); x != nil {
            log.Printf("[worker] painc: %v\n", x)
        }
    } ()
    defer w.Close()
    for {
        payload, err = conn.Receive()
        if err != nil {
            if err != io.EOF {
                log.Printf("workerError: %s\n", err.Error())
            }
            break
        }

        msgId, cmd, payload = protocol.ParseCommand(payload)

        switch cmd {
        case protocol.GRAB_JOB:
            err = w.handleGrabJob(msgId)
            break
        case protocol.WORK_DONE:
            jobId, _ := strconv.ParseInt(string(payload), 10, 0)
            err = w.handleDone(jobId)
            break
        case protocol.WORK_FAIL:
            jobId, _ := strconv.ParseInt(string(payload), 10, 0)
            err = w.handleFail(jobId)
            break
        case protocol.SCHED_LATER:
            parts := bytes.SplitN(payload, protocol.NULL_CHAR, 2)
            if len(parts) != 2 {
                log.Printf("Error: invalid format.")
                break
            }
            jobId, _ := strconv.ParseInt(string(parts[0]), 10, 0)
            delay, _ := strconv.ParseInt(string(parts[1]), 10, 0)
            err = w.handleSchedLater(jobId, delay)
            break
        case protocol.SLEEP:
            err = w.handleCommand(msgId, protocol.NOOP)
            break
        case protocol.PING:
            err = w.handleCommand(msgId, protocol.PONG)
            break
        case protocol.CAN_DO:
            err = w.handleCanDo(string(payload))
            break
        case protocol.CANT_DO:
            err = w.handleCanNoDo(string(payload))
            break
        default:
            err = w.handleCommand(msgId, protocol.UNKNOWN)
            break
        }
        if err != nil {
            if err != io.EOF {
                log.Printf("workerError: %s\n", err.Error())
            }
            break
        }

        if !w.alive {
            break
        }
    }
}


func (w *worker) Close() {
    defer w.sched.notifyJobTimer()
    defer w.conn.Close()
    w.sched.grabQueue.removeWorker(w)
    w.alive = false
    for k, _ := range w.jobQueue {
        w.sched.fail(k)
    }
    w.jobQueue = nil
    for _, Func := range w.funcs {
        w.sched.decrStatFunc(Func)
    }
    w = nil
}
