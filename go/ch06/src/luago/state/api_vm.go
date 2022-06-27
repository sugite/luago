package state

func (state *luaState) PC() int {
	return state.pc
}

func (state *luaState) AddPC(n int) {
	state.pc += n
}

func (state *luaState) Fetch() uint32 {
	i := state.proto.Code[state.pc]
	state.pc++
	return i
}

func (state *luaState) GetConst(idx int) {
	c := state.proto.Constants[idx]
	state.stack.push(c)
}

func (state *luaState) GetRK(rk int) {
	if rk > 0xFF {
		state.GetConst(rk & 0xFF)
	} else {
		state.PushValue(rk + 1)
	}
}
