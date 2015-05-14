package driver

import (
    "bytes"
    capn "github.com/glycerine/go-capnproto"
)

func ReadJob(payload []byte) (job Job, err error) {
    var buf = bytes.NewBuffer(payload)
    var s *capn.Segment
    if s, err = capn.ReadFromStream(buf, nil); err != nil {
        return
    }
    job = ReadRootJob(s)
    return
}

func (j Job) Clone () (job Job) {
    var s = capn.NewBuffer(nil)
    job = NewRootJob(s)
    job.SetName(j.Name())
    job.SetId(j.Id())
    job.SetFunc(j.Func())
    job.SetArgs(j.Args())
    job.SetTimeout(j.Timeout())
    job.SetSchedAt(j.SchedAt())
    job.SetRunAt(j.RunAt())
    job.SetStatus(j.Status())
    return
}
