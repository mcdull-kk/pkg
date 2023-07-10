package ecode

import (
	"testing"
	"time"

	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/stretchr/testify/assert"
	spb "google.golang.org/genproto/googleapis/rpc/status"
)

func TestEqual(t *testing.T) {
	var (
		err1 = Error(BadRequest, "test")
		err2 = Errorf(BadRequest, "test")
	)
	assert.Equal(t, err1, err2)
	assert.True(t, Equal(nil, nil))
}

func TestDetail(t *testing.T) {
	m := &timestamp.Timestamp{Seconds: time.Now().Unix()}
	st, _ := Error(BadRequest, "BadRequest").WithDetails(m)

	assert.Equal(t, "BadRequest", st.Message())
	assert.Equal(t, int(BadRequest), st.Code())
	assert.IsType(t, m, st.Details()[0])
}

func TestFromCode(t *testing.T) {
	err := FromCode(BadRequest)

	assert.Equal(t, int(BadRequest), err.Code())
	assert.Equal(t, "400", err.Message())
}

func TestFromProto(t *testing.T) {
	msg := &spb.Status{Code: 2233, Message: "error"}
	err := FromProto(msg)

	assert.Equal(t, 2233, err.Code())
	assert.Equal(t, "error", err.Message())

	m := &timestamp.Timestamp{Seconds: time.Now().Unix()}
	err = FromProto(m)
	assert.Equal(t, 500, err.Code())
	assert.Contains(t, err.Message(), "invalid proto message get")
}
