package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRepoContext_Querier(t *testing.T) {
	querier := &mockQuerier{}
	rc := NewRepoContext(context.Background(), querier)
	assert.Equal(t, querier, rc.Querier())
}

func TestRepoContext_AfterCommit(t *testing.T) {
	rc := NewRepoContext(context.Background(), &mockQuerier{})
	var called bool
	rc.AfterCommit(func() { called = true })
	rc.ExecAfterCommit(rc)
	assert.True(t, called)
}

func TestRepoContext_AfterCommitDo(t *testing.T) {
	rc := NewRepoContext(context.Background(), &mockQuerier{})
	var called bool
	rc.AfterCommitDo(func(ctx Context) { called = true })
	rc.ExecAfterCommit(rc)
	assert.True(t, called)
}

func ExampleRepoContext_usage() {
	querier := &mockQuerier{}
	rc := NewRepoContext(context.Background(), querier)
	rc.AfterCommit(func() {})
	rc.AfterCommitDo(func(Context) {})
	rc.ExecAfterCommit(rc)
	_ = rc.Querier()
	// Output:
}
