package driver

// AUTO GENERATED - DO NOT EDIT

import (
	C "github.com/glycerine/go-capnproto"
)

type Z C.Struct

func NewZ(s *C.Segment) Z        { return Z(s.NewStruct(0, 1)) }
func NewRootZ(s *C.Segment) Z    { return Z(s.NewRootStruct(0, 1)) }
func AutoNewZ(s *C.Segment) Z    { return Z(s.NewStructAR(0, 1)) }
func ReadRootZ(s *C.Segment) Z   { return Z(s.Root(0).ToStruct()) }
func (s Z) JobVec() Job_List     { return Job_List(C.Struct(s).GetObject(0)) }
func (s Z) SetJobVec(v Job_List) { C.Struct(s).SetObject(0, C.Object(v)) }

// capn.JSON_enabled == false so we stub MarshallJSON().
func (s Z) MarshalJSON() (bs []byte, err error) { return }

type Z_List C.PointerList

func NewZList(s *C.Segment, sz int) Z_List { return Z_List(s.NewCompositeList(0, 1, sz)) }
func (s Z_List) Len() int                  { return C.PointerList(s).Len() }
func (s Z_List) At(i int) Z                { return Z(C.PointerList(s).At(i).ToStruct()) }
func (s Z_List) ToArray() []Z {
	n := s.Len()
	a := make([]Z, n)
	for i := 0; i < n; i++ {
		a[i] = s.At(i)
	}
	return a
}
func (s Z_List) Set(i int, item Z) { C.PointerList(s).Set(i, C.Object(item)) }

type Job C.Struct

func NewJob(s *C.Segment) Job      { return Job(s.NewStruct(32, 4)) }
func NewRootJob(s *C.Segment) Job  { return Job(s.NewRootStruct(32, 4)) }
func AutoNewJob(s *C.Segment) Job  { return Job(s.NewStructAR(32, 4)) }
func ReadRootJob(s *C.Segment) Job { return Job(s.Root(0).ToStruct()) }
func (s Job) Id() int64            { return int64(C.Struct(s).Get64(0)) }
func (s Job) SetId(v int64)        { C.Struct(s).Set64(0, uint64(v)) }
func (s Job) Name() string         { return C.Struct(s).GetObject(0).ToText() }
func (s Job) SetName(v string)     { C.Struct(s).SetObject(0, s.Segment.NewText(v)) }
func (s Job) Func() string         { return C.Struct(s).GetObject(1).ToText() }
func (s Job) SetFunc(v string)     { C.Struct(s).SetObject(1, s.Segment.NewText(v)) }
func (s Job) Args() string         { return C.Struct(s).GetObject(2).ToText() }
func (s Job) SetArgs(v string)     { C.Struct(s).SetObject(2, s.Segment.NewText(v)) }
func (s Job) Timeout() int64       { return int64(C.Struct(s).Get64(8)) }
func (s Job) SetTimeout(v int64)   { C.Struct(s).Set64(8, uint64(v)) }
func (s Job) SchedAt() int64       { return int64(C.Struct(s).Get64(16)) }
func (s Job) SetSchedAt(v int64)   { C.Struct(s).Set64(16, uint64(v)) }
func (s Job) RunAt() int64         { return int64(C.Struct(s).Get64(24)) }
func (s Job) SetRunAt(v int64)     { C.Struct(s).Set64(24, uint64(v)) }
func (s Job) Status() string       { return C.Struct(s).GetObject(3).ToText() }
func (s Job) SetStatus(v string)   { C.Struct(s).SetObject(3, s.Segment.NewText(v)) }

// capn.JSON_enabled == false so we stub MarshallJSON().
func (s Job) MarshalJSON() (bs []byte, err error) { return }

type Job_List C.PointerList

func NewJobList(s *C.Segment, sz int) Job_List { return Job_List(s.NewCompositeList(32, 4, sz)) }
func (s Job_List) Len() int                    { return C.PointerList(s).Len() }
func (s Job_List) At(i int) Job                { return Job(C.PointerList(s).At(i).ToStruct()) }
func (s Job_List) ToArray() []Job {
	n := s.Len()
	a := make([]Job, n)
	for i := 0; i < n; i++ {
		a[i] = s.At(i)
	}
	return a
}
func (s Job_List) Set(i int, item Job) { C.PointerList(s).Set(i, C.Object(item)) }
