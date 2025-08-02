package raft

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"raft-redis-cluster/store"

	"github.com/hashicorp/raft"
)

var _ raft.FSM = (*StateMachine)(nil)

type Op int

const (
	Put Op = iota
	Del
)

type KVCmd struct {
	Op  Op     `json:"op"`
	Key []byte `json:"key"`
	Val []byte `json:"val"`
}

func NewStateMachine(store store.Store) *StateMachine {
	return &StateMachine{store}
}

type StateMachine struct {
	store store.Store
}

// Apply　は、Raft ログエントリをキーバリューストアに適用します。
func (s *StateMachine) Apply(log *raft.Log) any {
	ctx := context.Background()
	c := KVCmd{}

	err := json.Unmarshal(log.Data, &c)
	if err != nil {
		return err
	}

	return s.handleRequest(ctx, c)
}

// Restore　は、Key-Value ストアの状態を以前の状態に保存します。
func (s *StateMachine) Restore(rc io.ReadCloser) error {
	return s.store.Restore(rc)
}

// Snapshot　は、Key-Value ストアのKVSnapshotを返します。
func (s *StateMachine) Snapshot() (raft.FSMSnapshot, error) {
	rc, err := s.store.Snapshot()
	if err != nil {
		return nil, err
	}
	return &KVSnapshot{rc}, nil
}

var ErrUnknownOp = errors.New("unknown op")

func (s *StateMachine) handleRequest(ctx context.Context, cmd KVCmd) error {
	switch cmd.Op {
	case Put:
		return s.store.Put(ctx, cmd.Key, cmd.Val)
	case Del:
		return s.store.Delete(ctx, cmd.Key)
	default:
		return ErrUnknownOp
	}
}
