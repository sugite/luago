package binchunk

type header struct {
	signature       [4]byte // 签名，一个魔数  Esc、L、u、a的ASCII码
	version         byte    // 版本号，5.3.4，大版本号5*16+小版本号
	format          byte    // 格式号 0
	luacData        [6]byte // 0x1993 \r\n\x1a\n ，进一步校验
	cintSize        byte    // 4
	sizetSize       byte    // 8
	instructionSize byte    // 4
	luaIntegerSize  byte    // 8
	luaNumberSize   byte    // 8
	luacInt         int64
	luacNum         float64
}

type Prototype struct {
	Source          string
	LineDefined     uint32
	LastLineDefined uint32
	NumParams       byte
	IsVararg        byte
	MaxStackSize    byte          // 虚拟寄存器数量
	Code            []uint32      // 指令表 一个指令4个字节
	Constants       []interface{} // 常量表 每个常量以1字节tag开头 包括 nil bool int float string
	Upvalues        []Upvalue
	Protos          []*Prototype
	LineInfo        []uint32
	LocVars         []LocVar
	UpvalueNames    []string
}

type Upvalue struct {
	Instack byte
	Idx     byte
}

type LocVar struct {
	VarName string
	StartPC uint32
	EndPC   uint32
}

type binaryChunk struct {
	header                  // 头部
	sizeUpvalues byte       // 主函数upvalue数量
	mainFunc     *Prototype // 主函数原型
}

const (
	LUA_SIGNATURE    = "\x1bLua"
	LUAC_VERSION     = 0x53
	LUAC_FORMAT      = 0
	LUAC_DATA        = "\x19\x93\r\n\x1a\n"
	CINT_SIZE        = 4
	CSIZET_SIZE      = 8
	INSTRUCTION_SIZE = 4
	LUA_INTEGER_SIZE = 8
	LUA_NUMBER_SIZE  = 8
	LUAC_INT         = 0x5678
	LUAC_NUM         = 370.5
)

const (
	TAG_NIL       = 0x00
	TAG_BOOLEAN   = 0x01
	TAG_NUMBER    = 0x03
	TAG_INTEGER   = 0x13
	TAG_SHORT_STR = 0x04
	TAG_LONG_STR  = 0x14
)

func Undump(data []byte) *Prototype {
	reader := &reader{data}
	reader.checkHeader()        // 校验头部
	reader.readByte()           // 跳过Upvalue数量
	return reader.readProto("") // 读取函数原型
}
