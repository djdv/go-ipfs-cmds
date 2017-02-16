package cmds

import (
	"fmt"
	"io"

	"github.com/ipfs/go-ipfs-cmds/cmdsutil"
)

func NewChanResponsePair(req Request) (ResponseEmitter, Response) {
	ch := make(chan interface{})

	r := &chanResponse{
		req: req,
		ch:  ch,
	}

	re := &chanResponseEmitter{
		ch:     ch,
		length: &r.length,
		err:    &r.err,
	}

	return re, r
}

type chanResponse struct {
	req Request

	err    *cmdsutil.Error
	length uint64

	ch <-chan interface{}
}

func (r *chanResponse) Request() Request {
	if r == nil {
		return nil
	}

	return r.req
}

func (r *chanResponse) Error() *cmdsutil.Error {
	if r == nil {
		return nil
	}

	return r.err
}

func (r *chanResponse) Length() uint64 {
	if r == nil {
		return 0
	}

	return r.length
}

func (r *chanResponse) Next() (interface{}, error) {
	if r == nil {
		return nil, io.EOF
	}

	v, ok := <-r.ch
	if ok {
		return v, nil
	}

	return nil, io.EOF
}

type chanResponseEmitter struct {
	ch chan<- interface{}

	length *uint64
	err    **cmdsutil.Error

	emitted bool

	tees []ResponseEmitter
}

func (re *chanResponseEmitter) SetError(err interface{}, t cmdsutil.ErrorType) {
	// don't change value after emitting
	/*
		if re.emitted {
			return
		}
	*/

	*re.err = &cmdsutil.Error{Message: fmt.Sprint(err), Code: t}

	for _, re_ := range re.tees {
		re_.SetError(err, t)
	}
}

func (re *chanResponseEmitter) SetLength(l uint64) {
	// don't change value after emitting
	if re.emitted {
		return
	}

	*re.length = l

	for _, re_ := range re.tees {
		re_.SetLength(l)
	}
}

func (re *chanResponseEmitter) Head() Head {
	return Head{
		Len: *re.length,
		Err: *re.err,
	}
}

func (re *chanResponseEmitter) Close() error {
	close(re.ch)
	re.ch = nil

	return nil
}

func (re *chanResponseEmitter) Emit(v interface{}) error {
	re.emitted = true

	if re.ch == nil {
		return fmt.Errorf("emitter closed")
	}

	re.ch <- v

	for _, re_ := range re.tees {
		go re_.Emit(v)
	}

	return nil
}

func (re *chanResponseEmitter) Tee(re_ ResponseEmitter) {
	if re_ == nil {
		return
	}

	re.tees = append(re.tees, re_)

	if re.emitted {
		re_.SetLength(*re.length)
	}

	if re.err != nil && *re.err != nil {
		re_.SetError((*re.err).Message, (*re.err).Code)
	}
}