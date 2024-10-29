package recver

import "CentralizedControl/ins_lite/proto/types"

type SessionMessage struct {
	LzmaDictSize           int32
	SessionId              int64
	Flags                  types.VarUInt32
	ClientId               int64 `ins:"(Flags&4)!= 0"`
	TimeOffset             int16
	UnUsed                 byte
	TransientToken         int32
	MsgNotLastLogin        int16
	ClientEncryptionSecret types.ListValue[byte, types.VarUInt32] `ins:"(Flags & 128)!=0"`
}
