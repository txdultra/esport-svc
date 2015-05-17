package libs

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/Lupino/periodic/driver"
	"github.com/Lupino/periodic/protocol"
	"github.com/glycerine/go-capnproto"
)

type TimerJobFunc func(jobName, args string, otherArgs ...interface{}) (schedLater int, err error)

type TimerJobSys struct {
	entryPoint string
}

func NewTimerJobSys(entryPoint string) *TimerJobSys {
	tts := &TimerJobSys{}
	tts.entryPoint = entryPoint
	return tts
}

func (t *TimerJobSys) SetJob(name, kind, args string, timeout int64, runTime time.Time) error {
	parts := strings.SplitN(t.entryPoint, "://", 2)
	c, err := net.Dial(parts[0], parts[1])
	if err != nil {
		log.Fatal(err)
		return err
	}
	conn := protocol.NewClientConn(c)
	defer conn.Close()
	err = conn.Send(protocol.TYPE_CLIENT.Bytes())
	if err != nil {
		log.Fatal(err)
		return err
	}
	//build job
	var s = capn.NewBuffer(nil)
	var job = driver.NewRootJob(s)
	job.SetName(name)
	job.SetFunc(kind)
	job.SetArgs(args)
	job.SetTimeout(timeout)
	job.SetSchedAt(runTime.Unix())

	var msgId = []byte("")
	buf := bytes.NewBuffer(nil)
	buf.Write(msgId)
	buf.Write(protocol.NULL_CHAR)
	buf.WriteByte(byte(protocol.SUBMIT_JOB))
	buf.Write(protocol.NULL_CHAR)
	job.Segment.WriteTo(buf)
	err = conn.Send(buf.Bytes())
	if err != nil {
		log.Fatal(err)
		return err
	}
	payload, err := conn.Receive()
	if err != nil {
		log.Fatal(err)
		return err
	}
	_, cmd, _ := protocol.ParseCommand(payload)
	rlt := cmd.String()
	if rlt == "SUCCESS" {
		return nil
	}
	return fmt.Errorf(rlt)
}

func (t *TimerJobSys) RemoveJob(name, kind string) error {
	parts := strings.SplitN(t.entryPoint, "://", 2)
	c, err := net.Dial(parts[0], parts[1])
	if err != nil {
		log.Fatal(err)
		return err
	}
	conn := protocol.NewClientConn(c)
	defer conn.Close()
	err = conn.Send(protocol.TYPE_CLIENT.Bytes())
	if err != nil {
		log.Fatal(err)
		return err
	}
	//build job
	var s = capn.NewBuffer(nil)
	var job = driver.NewRootJob(s)
	job.SetName(name)
	job.SetFunc(kind)

	buf := bytes.NewBuffer(nil)
	buf.Write([]byte(""))
	buf.Write(protocol.NULL_CHAR)
	buf.WriteByte(byte(protocol.REMOVE_JOB))
	buf.Write(protocol.NULL_CHAR)
	job.Segment.WriteTo(buf)
	err = conn.Send(buf.Bytes())
	if err != nil {
		log.Fatal(err)
		return err
	}
	payload, err := conn.Receive()
	if err != nil {
		log.Fatal(err)
		return err
	}
	_, cmd, _ := protocol.ParseCommand(payload)
	rlt := cmd.String()
	if rlt == "SUCCESS" {
		return nil
	}
	return fmt.Errorf(rlt)
}

func (t *TimerJobSys) RunGrabJob(kind string, jobFunc TimerJobFunc) error {
	parts := strings.SplitN(t.entryPoint, "://", 2)
	for {
		c, err := net.Dial(parts[0], parts[1])
		if err != nil {
			if err != io.EOF {
				log.Printf("Error: %s\n", err.Error())
			}
			log.Printf("Wait 5 second to reconnecting")
			time.Sleep(5 * time.Second)
			continue
		}
		conn := protocol.NewClientConn(c)
		err = t.handleWorker(conn, kind, jobFunc)
		if err != nil {
			if err != io.EOF {
				log.Printf("Error: %s\n", err.Error())
			}
		}
		conn.Close()
	}
}

func (t *TimerJobSys) handleWorker(conn protocol.Conn, kind string, jobFunc TimerJobFunc) (err error) {
	err = conn.Send(protocol.TYPE_WORKER.Bytes())
	if err != nil {
		log.Fatal(err)
		return
	}
	var msgId = []byte("")
	buf := bytes.NewBuffer(nil)
	buf.Write(msgId)
	buf.Write(protocol.NULL_CHAR)
	buf.WriteByte(byte(protocol.CAN_DO))
	buf.Write(protocol.NULL_CHAR)
	buf.WriteString(kind)
	err = conn.Send(buf.Bytes())
	if err != nil {
		log.Fatal(err)
		return
	}

	var payload []byte
	var job driver.Job
	var jobHandle []byte
	for {
		buf = bytes.NewBuffer(nil)
		buf.Write(msgId)
		buf.Write(protocol.NULL_CHAR)
		buf.Write(protocol.GRAB_JOB.Bytes())
		err = conn.Send(buf.Bytes())
		if err != nil {
			log.Fatal(err)
			return
		}
		payload, err = conn.Receive()
		if err != nil {
			log.Fatal(err)
			return
		}
		job, jobHandle, err = t.extraJob(payload)
		//运行任务
		schedLater, err := jobFunc(job.Name(), job.Args(), job.Status())

		buf = bytes.NewBuffer(nil)
		buf.Write(msgId)
		buf.Write(protocol.NULL_CHAR)
		if err != nil {
			buf.WriteByte(byte(protocol.WORK_FAIL))
		} else if schedLater > 0 {
			buf.WriteByte(byte(protocol.SCHED_LATER))
		} else {
			buf.WriteByte(byte(protocol.WORK_DONE))
		}
		buf.Write(protocol.NULL_CHAR)
		buf.Write(jobHandle)
		if schedLater > 0 {
			buf.Write(protocol.NULL_CHAR)
			buf.WriteString(strconv.Itoa(schedLater))
		}
		err = conn.Send(buf.Bytes())
		if err != nil {
			log.Fatal(err)
		}
	}
}

func (t *TimerJobSys) extraJob(payload []byte) (job driver.Job, jobHandle []byte, err error) {
	parts := bytes.SplitN(payload, protocol.NULL_CHAR, 4)
	if len(parts) != 4 {
		err = errors.New("Invalid payload " + string(payload))
		log.Fatal(err)
		return
	}
	job, err = driver.ReadJob(parts[3])
	jobHandle = parts[2]
	return
}
