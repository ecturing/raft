package messageHandler

import (
	"raft/pkg/rpc/Entries"
	"reflect"
	"testing"
)

func TestRpcMsgHandler_Vote(t *testing.T) {
	type fields struct {
		nodeList []*Entries.NetMeta
	}
	type args struct {
		arg Entries.RequestVoteArgs
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   Entries.RequestVoteReply
	}{
		// TODO: Add test cases.
		{"1", fields{initNetMetaList()}, args{Entries.RequestVoteArgs{0}}, Entries.RequestVoteReply{VoteGranted: false}},
		{"2", fields{initNetMetaList()}, args{Entries.RequestVoteArgs{1}}, Entries.RequestVoteReply{VoteGranted: false}},
		{"3", fields{initNetMetaList()}, args{Entries.RequestVoteArgs{2}}, Entries.RequestVoteReply{VoteGranted: false}},
		{"4", fields{initNetMetaList()}, args{Entries.RequestVoteArgs{3}}, Entries.RequestVoteReply{VoteGranted: false}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := RpcMsgHandler{
				nodeList: tt.fields.nodeList,
			}
			if got := r.SendVote(tt.args.arg); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SendVote() = %v, want %v", got, tt.want)
			}
		})
	}
}

func initNetMetaList() []*Entries.NetMeta {
	var nodelist []*Entries.NetMeta
	for i := 0; i < 10; i++ {
		nodelist = append(nodelist, &Entries.NetMeta{Ip: "hshuu", Port: "3123"})
	}
	return nodelist
}
