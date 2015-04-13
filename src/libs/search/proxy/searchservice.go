// Autogenerated by Thrift Compiler (0.9.2)
// DO NOT EDIT UNLESS YOU ARE SURE THAT YOU KNOW WHAT YOU ARE DOING

package proxy

import (
	"bytes"
	"fmt"

	"github.com/thrift"
)

// (needed to ensure safety because of naive import list construction.)
var _ = thrift.ZERO
var _ = fmt.Printf
var _ = bytes.Equal

type SearchService interface {
	// Parameters:
	//  - Option
	//  - Words
	//  - IdxName
	Query(option *SearchOptions, words string, idxName string) (r *SearchResult_, err error)
}

type SearchServiceClient struct {
	Transport       thrift.TTransport
	ProtocolFactory thrift.TProtocolFactory
	InputProtocol   thrift.TProtocol
	OutputProtocol  thrift.TProtocol
	SeqId           int32
}

func NewSearchServiceClientFactory(t thrift.TTransport, f thrift.TProtocolFactory) *SearchServiceClient {
	return &SearchServiceClient{Transport: t,
		ProtocolFactory: f,
		InputProtocol:   f.GetProtocol(t),
		OutputProtocol:  f.GetProtocol(t),
		SeqId:           0,
	}
}

func NewSearchServiceClientProtocol(t thrift.TTransport, iprot thrift.TProtocol, oprot thrift.TProtocol) *SearchServiceClient {
	return &SearchServiceClient{Transport: t,
		ProtocolFactory: nil,
		InputProtocol:   iprot,
		OutputProtocol:  oprot,
		SeqId:           0,
	}
}

// Parameters:
//  - Option
//  - Words
//  - IdxName
func (p *SearchServiceClient) Query(option *SearchOptions, words string, idxName string) (r *SearchResult_, err error) {
	if err = p.sendQuery(option, words, idxName); err != nil {
		return
	}
	return p.recvQuery()
}

func (p *SearchServiceClient) sendQuery(option *SearchOptions, words string, idxName string) (err error) {
	oprot := p.OutputProtocol
	if oprot == nil {
		oprot = p.ProtocolFactory.GetProtocol(p.Transport)
		p.OutputProtocol = oprot
	}
	p.SeqId++
	if err = oprot.WriteMessageBegin("Query", thrift.CALL, p.SeqId); err != nil {
		return
	}
	args := QueryArgs{
		Option:  option,
		Words:   words,
		IdxName: idxName,
	}
	if err = args.Write(oprot); err != nil {
		return
	}
	if err = oprot.WriteMessageEnd(); err != nil {
		return
	}
	return oprot.Flush()
}

func (p *SearchServiceClient) recvQuery() (value *SearchResult_, err error) {
	iprot := p.InputProtocol
	if iprot == nil {
		iprot = p.ProtocolFactory.GetProtocol(p.Transport)
		p.InputProtocol = iprot
	}
	_, mTypeId, seqId, err := iprot.ReadMessageBegin()
	if err != nil {
		return
	}
	if mTypeId == thrift.EXCEPTION {
		error8 := thrift.NewTApplicationException(thrift.UNKNOWN_APPLICATION_EXCEPTION, "Unknown Exception")
		var error9 error
		error9, err = error8.Read(iprot)
		if err != nil {
			return
		}
		if err = iprot.ReadMessageEnd(); err != nil {
			return
		}
		err = error9
		return
	}
	if p.SeqId != seqId {
		err = thrift.NewTApplicationException(thrift.BAD_SEQUENCE_ID, "Query failed: out of sequence response")
		return
	}
	result := QueryResult{}
	if err = result.Read(iprot); err != nil {
		return
	}
	if err = iprot.ReadMessageEnd(); err != nil {
		return
	}
	value = result.GetSuccess()
	return
}

type SearchServiceProcessor struct {
	processorMap map[string]thrift.TProcessorFunction
	handler      SearchService
}

func (p *SearchServiceProcessor) AddToProcessorMap(key string, processor thrift.TProcessorFunction) {
	p.processorMap[key] = processor
}

func (p *SearchServiceProcessor) GetProcessorFunction(key string) (processor thrift.TProcessorFunction, ok bool) {
	processor, ok = p.processorMap[key]
	return processor, ok
}

func (p *SearchServiceProcessor) ProcessorMap() map[string]thrift.TProcessorFunction {
	return p.processorMap
}

func NewSearchServiceProcessor(handler SearchService) *SearchServiceProcessor {

	self10 := &SearchServiceProcessor{handler: handler, processorMap: make(map[string]thrift.TProcessorFunction)}
	self10.processorMap["Query"] = &searchServiceProcessorQuery{handler: handler}
	return self10
}

func (p *SearchServiceProcessor) Process(iprot, oprot thrift.TProtocol) (success bool, err thrift.TException) {
	name, _, seqId, err := iprot.ReadMessageBegin()
	if err != nil {
		return false, err
	}
	if processor, ok := p.GetProcessorFunction(name); ok {
		return processor.Process(seqId, iprot, oprot)
	}
	iprot.Skip(thrift.STRUCT)
	iprot.ReadMessageEnd()
	x11 := thrift.NewTApplicationException(thrift.UNKNOWN_METHOD, "Unknown function "+name)
	oprot.WriteMessageBegin(name, thrift.EXCEPTION, seqId)
	x11.Write(oprot)
	oprot.WriteMessageEnd()
	oprot.Flush()
	return false, x11

}

type searchServiceProcessorQuery struct {
	handler SearchService
}

func (p *searchServiceProcessorQuery) Process(seqId int32, iprot, oprot thrift.TProtocol) (success bool, err thrift.TException) {
	args := QueryArgs{}
	if err = args.Read(iprot); err != nil {
		iprot.ReadMessageEnd()
		x := thrift.NewTApplicationException(thrift.PROTOCOL_ERROR, err.Error())
		oprot.WriteMessageBegin("Query", thrift.EXCEPTION, seqId)
		x.Write(oprot)
		oprot.WriteMessageEnd()
		oprot.Flush()
		return false, err
	}

	iprot.ReadMessageEnd()
	result := QueryResult{}
	var retval *SearchResult_
	var err2 error
	if retval, err2 = p.handler.Query(args.Option, args.Words, args.IdxName); err2 != nil {
		x := thrift.NewTApplicationException(thrift.INTERNAL_ERROR, "Internal error processing Query: "+err2.Error())
		oprot.WriteMessageBegin("Query", thrift.EXCEPTION, seqId)
		x.Write(oprot)
		oprot.WriteMessageEnd()
		oprot.Flush()
		return true, err2
	} else {
		result.Success = retval
	}
	if err2 = oprot.WriteMessageBegin("Query", thrift.REPLY, seqId); err2 != nil {
		err = err2
	}
	if err2 = result.Write(oprot); err == nil && err2 != nil {
		err = err2
	}
	if err2 = oprot.WriteMessageEnd(); err == nil && err2 != nil {
		err = err2
	}
	if err2 = oprot.Flush(); err == nil && err2 != nil {
		err = err2
	}
	if err != nil {
		return
	}
	return true, err
}

// HELPER FUNCTIONS AND STRUCTURES

type QueryArgs struct {
	Option  *SearchOptions `thrift:"option,1,required" json:"option"`
	Words   string         `thrift:"words,2,required" json:"words"`
	IdxName string         `thrift:"idxName,3,required" json:"idxName"`
}

func NewQueryArgs() *QueryArgs {
	return &QueryArgs{}
}

var QueryArgs_Option_DEFAULT *SearchOptions

func (p *QueryArgs) GetOption() *SearchOptions {
	if !p.IsSetOption() {
		return QueryArgs_Option_DEFAULT
	}
	return p.Option
}

func (p *QueryArgs) GetWords() string {
	return p.Words
}

func (p *QueryArgs) GetIdxName() string {
	return p.IdxName
}
func (p *QueryArgs) IsSetOption() bool {
	return p.Option != nil
}

func (p *QueryArgs) Read(iprot thrift.TProtocol) error {
	if _, err := iprot.ReadStructBegin(); err != nil {
		return fmt.Errorf("%T read error: %s", p, err)
	}
	for {
		_, fieldTypeId, fieldId, err := iprot.ReadFieldBegin()
		if err != nil {
			return fmt.Errorf("%T field %d read error: %s", p, fieldId, err)
		}
		if fieldTypeId == thrift.STOP {
			break
		}
		switch fieldId {
		case 1:
			if err := p.ReadField1(iprot); err != nil {
				return err
			}
		case 2:
			if err := p.ReadField2(iprot); err != nil {
				return err
			}
		case 3:
			if err := p.ReadField3(iprot); err != nil {
				return err
			}
		default:
			if err := iprot.Skip(fieldTypeId); err != nil {
				return err
			}
		}
		if err := iprot.ReadFieldEnd(); err != nil {
			return err
		}
	}
	if err := iprot.ReadStructEnd(); err != nil {
		return fmt.Errorf("%T read struct end error: %s", p, err)
	}
	return nil
}

func (p *QueryArgs) ReadField1(iprot thrift.TProtocol) error {
	p.Option = &SearchOptions{}
	if err := p.Option.Read(iprot); err != nil {
		return fmt.Errorf("%T error reading struct: %s", p.Option, err)
	}
	return nil
}

func (p *QueryArgs) ReadField2(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadString(); err != nil {
		return fmt.Errorf("error reading field 2: %s", err)
	} else {
		p.Words = v
	}
	return nil
}

func (p *QueryArgs) ReadField3(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadString(); err != nil {
		return fmt.Errorf("error reading field 3: %s", err)
	} else {
		p.IdxName = v
	}
	return nil
}

func (p *QueryArgs) Write(oprot thrift.TProtocol) error {
	if err := oprot.WriteStructBegin("Query_args"); err != nil {
		return fmt.Errorf("%T write struct begin error: %s", p, err)
	}
	if err := p.writeField1(oprot); err != nil {
		return err
	}
	if err := p.writeField2(oprot); err != nil {
		return err
	}
	if err := p.writeField3(oprot); err != nil {
		return err
	}
	if err := oprot.WriteFieldStop(); err != nil {
		return fmt.Errorf("write field stop error: %s", err)
	}
	if err := oprot.WriteStructEnd(); err != nil {
		return fmt.Errorf("write struct stop error: %s", err)
	}
	return nil
}

func (p *QueryArgs) writeField1(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("option", thrift.STRUCT, 1); err != nil {
		return fmt.Errorf("%T write field begin error 1:option: %s", p, err)
	}
	if err := p.Option.Write(oprot); err != nil {
		return fmt.Errorf("%T error writing struct: %s", p.Option, err)
	}
	if err := oprot.WriteFieldEnd(); err != nil {
		return fmt.Errorf("%T write field end error 1:option: %s", p, err)
	}
	return err
}

func (p *QueryArgs) writeField2(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("words", thrift.STRING, 2); err != nil {
		return fmt.Errorf("%T write field begin error 2:words: %s", p, err)
	}
	if err := oprot.WriteString(string(p.Words)); err != nil {
		return fmt.Errorf("%T.words (2) field write error: %s", p, err)
	}
	if err := oprot.WriteFieldEnd(); err != nil {
		return fmt.Errorf("%T write field end error 2:words: %s", p, err)
	}
	return err
}

func (p *QueryArgs) writeField3(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("idxName", thrift.STRING, 3); err != nil {
		return fmt.Errorf("%T write field begin error 3:idxName: %s", p, err)
	}
	if err := oprot.WriteString(string(p.IdxName)); err != nil {
		return fmt.Errorf("%T.idxName (3) field write error: %s", p, err)
	}
	if err := oprot.WriteFieldEnd(); err != nil {
		return fmt.Errorf("%T write field end error 3:idxName: %s", p, err)
	}
	return err
}

func (p *QueryArgs) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("QueryArgs(%+v)", *p)
}

type QueryResult struct {
	Success *SearchResult_ `thrift:"success,0" json:"success"`
}

func NewQueryResult() *QueryResult {
	return &QueryResult{}
}

var QueryResult_Success_DEFAULT *SearchResult_

func (p *QueryResult) GetSuccess() *SearchResult_ {
	if !p.IsSetSuccess() {
		return QueryResult_Success_DEFAULT
	}
	return p.Success
}
func (p *QueryResult) IsSetSuccess() bool {
	return p.Success != nil
}

func (p *QueryResult) Read(iprot thrift.TProtocol) error {
	if _, err := iprot.ReadStructBegin(); err != nil {
		return fmt.Errorf("%T read error: %s", p, err)
	}
	for {
		_, fieldTypeId, fieldId, err := iprot.ReadFieldBegin()
		if err != nil {
			return fmt.Errorf("%T field %d read error: %s", p, fieldId, err)
		}
		if fieldTypeId == thrift.STOP {
			break
		}
		switch fieldId {
		case 0:
			if err := p.ReadField0(iprot); err != nil {
				return err
			}
		default:
			if err := iprot.Skip(fieldTypeId); err != nil {
				return err
			}
		}
		if err := iprot.ReadFieldEnd(); err != nil {
			return err
		}
	}
	if err := iprot.ReadStructEnd(); err != nil {
		return fmt.Errorf("%T read struct end error: %s", p, err)
	}
	return nil
}

func (p *QueryResult) ReadField0(iprot thrift.TProtocol) error {
	p.Success = &SearchResult_{}
	if err := p.Success.Read(iprot); err != nil {
		return fmt.Errorf("%T error reading struct: %s", p.Success, err)
	}
	return nil
}

func (p *QueryResult) Write(oprot thrift.TProtocol) error {
	if err := oprot.WriteStructBegin("Query_result"); err != nil {
		return fmt.Errorf("%T write struct begin error: %s", p, err)
	}
	if err := p.writeField0(oprot); err != nil {
		return err
	}
	if err := oprot.WriteFieldStop(); err != nil {
		return fmt.Errorf("write field stop error: %s", err)
	}
	if err := oprot.WriteStructEnd(); err != nil {
		return fmt.Errorf("write struct stop error: %s", err)
	}
	return nil
}

func (p *QueryResult) writeField0(oprot thrift.TProtocol) (err error) {
	if p.IsSetSuccess() {
		if err := oprot.WriteFieldBegin("success", thrift.STRUCT, 0); err != nil {
			return fmt.Errorf("%T write field begin error 0:success: %s", p, err)
		}
		if err := p.Success.Write(oprot); err != nil {
			return fmt.Errorf("%T error writing struct: %s", p.Success, err)
		}
		if err := oprot.WriteFieldEnd(); err != nil {
			return fmt.Errorf("%T write field end error 0:success: %s", p, err)
		}
	}
	return err
}

func (p *QueryResult) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("QueryResult(%+v)", *p)
}
