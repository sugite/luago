package binchunk

import (
	"encoding/binary"
	"math"
)

type reader struct {
	data []byte
}

func (rdr *reader) readByte() byte {
	b := rdr.data[0]
	rdr.data = rdr.data[1:]
	return b
}

func (rdr *reader) readUint32() uint32 {
	i := binary.LittleEndian.Uint32(rdr.data)
	rdr.data = rdr.data[4:]
	return i
}

func (rdr *reader) readUint64() uint64 {
	i := binary.LittleEndian.Uint64(rdr.data)
	rdr.data = rdr.data[8:]
	return i
}

func (rdr *reader) readLuaInteger() int64 {
	return int64(rdr.readUint64())
}

func (rdr *reader) readLuaNumber() float64 {
	return math.Float64frombits(rdr.readUint64())
}

func (rdr *reader) readString() string {
	size := uint(rdr.readByte())
	if size == 0 {
		return ""
	}

	if size == 0xFF {
		size = uint(rdr.readUint64())
	}
	bytes := rdr.readBytes(size - 1)
	return string(bytes)
}

func (rdr *reader) readBytes(n uint) []byte {
	bytes := rdr.data[:n]
	rdr.data = rdr.data[n:]
	return bytes
}

func (rdr *reader) checkHeader() {
	if string(rdr.readBytes(4)) != LUA_SIGNATURE {
		panic("not a precompiled chunk!")
	} else if rdr.readByte() != LUAC_VERSION {
		panic("version mismatch!")
	} else if rdr.readByte() != LUAC_FORMAT {
		panic("format mismatch!")
	} else if string(rdr.readBytes(6)) != LUAC_DATA {
		panic("corrupted!")
	} else if rdr.readByte() != CINT_SIZE {
		panic("int size mismatch!")
	} else if rdr.readByte() != CSIZET_SIZE {
		panic("size_t size mismatch!")
	} else if rdr.readByte() != INSTRUCTION_SIZE {
		panic("instruction size mismatch!")
	} else if rdr.readByte() != LUA_INTEGER_SIZE {
		panic("lua_integer size mismatch!")
	} else if rdr.readByte() != LUA_NUMBER_SIZE {
		panic("lua_number size mismatch!")
	} else if rdr.readLuaInteger() != LUAC_INT {
		panic("endianness mismatch!")
	} else if rdr.readLuaNumber() != LUAC_NUM {
		panic("float format mismatch!")
	}
}

func (rdr *reader) readProtos(parentSource string) []*Prototype {
	protos := make([]*Prototype, rdr.readUint32())
	for i := range protos {
		protos[i] = rdr.readProto(parentSource)
	}
	return protos
}

func (rdr *reader) readProto(parentSource string) *Prototype {
	source := rdr.readString()
	if source == "" {
		source = parentSource
	}

	return &Prototype{
		Source:          source,
		LineDefined:     rdr.readUint32(),
		LastLineDefined: rdr.readUint32(),
		NumParams:       rdr.readByte(),
		IsVararg:        rdr.readByte(),
		MaxStackSize:    rdr.readByte(),
		Code:            rdr.readCode(),
		Constants:       rdr.readConstants(),
		Upvalues:        rdr.readUpvalues(),
		Protos:          rdr.readProtos(source),
		LineInfo:        rdr.readLineInfo(),
		LocVars:         rdr.readLocVars(),
		UpvalueNames:    rdr.readUpvalueNames(),
	}
}

func (rdr *reader) readCode() []uint32 {
	code := make([]uint32, rdr.readUint32())
	for i := range code {
		code[i] = rdr.readUint32()
	}
	return code
}

func (rdr *reader) readConstants() []interface{} {
	constants := make([]interface{}, rdr.readUint32())
	for i := range constants {
		constants[i] = rdr.readConstant()
	}
	return constants
}

func (rdr *reader) readConstant() interface{} {
	switch rdr.readByte() {
	case TAG_NIL:
		return nil
	case TAG_BOOLEAN:
		return rdr.readByte() != 0
	case TAG_INTEGER:
		return rdr.readLuaInteger()
	case TAG_NUMBER:
		return rdr.readLuaNumber()
	case TAG_SHORT_STR:
		return rdr.readString()
	case TAG_LONG_STR:
		return rdr.readString()
	default:
		panic("corrupted!")
	}
}

func (rdr *reader) readUpvalues() []Upvalue {
	upvalues := make([]Upvalue, rdr.readUint32())
	for i := range upvalues {
		upvalues[i] = Upvalue{
			Instack: rdr.readByte(),
			Idx:     rdr.readByte(),
		}
	}
	return upvalues
}

func (rdr *reader) readLineInfo() []uint32 {
	lineInfo := make([]uint32, rdr.readUint32())
	for i := range lineInfo {
		lineInfo[i] = rdr.readUint32()
	}
	return lineInfo
}

func (rdr *reader) readLocVars() []LocVar {
	locVars := make([]LocVar, rdr.readUint32())
	for i := range locVars {
		locVars[i] = LocVar{
			VarName: rdr.readString(),
			StartPC: rdr.readUint32(),
			EndPC:   rdr.readUint32(),
		}
	}

	return locVars
}

func (rdr *reader) readUpvalueNames() []string {
	names := make([]string, rdr.readUint32())
	for i := range names {
		names[i] = rdr.readString()
	}
	return names
}
