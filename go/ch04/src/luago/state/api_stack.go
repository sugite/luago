package state

func (state *luaState) GetTop() int {
	return state.stack.top
}

func (state *luaState) AbsIndex(idx int) int {
	return state.stack.absIndex(idx)
}

func (state *luaState) CheckStack(n int) bool {
	state.stack.check(n)
	return true
}

func (state *luaState) Pop(n int) {
	for i := 0; i < n; i++ {
		state.stack.pop()
	}
}

func (state *luaState) Copy(fromIdx, toIdx int) {
	val := state.stack.get(fromIdx)
	state.stack.set(toIdx, val)
}

func (state *luaState) PushValue(idx int) {
	val := state.stack.get(idx)
	state.stack.push(val)
}

func (state *luaState) Replace(idx int) {
	val := state.stack.pop()
	state.stack.set(idx, val)
}

func (state *luaState) Insert(idx int) {
	state.Rotate(idx, 1)
}

func (state *luaState) Remove(idx int) {
	state.Rotate(idx, -1)
	state.Pop(1)
}

func (state *luaState) Rotate(idx, n int) {
	t := state.stack.top - 1
	p := state.stack.absIndex(idx) - 1
	var m int
	if n >= 0 {
		m = t - n
	} else {
		m = p - n - 1
	}
	state.stack.reverse(p, m)
	state.stack.reverse(m+1, t)
	state.stack.reverse(p, t)
}

func (state *luaState) SetTop(idx int) {
	newTop := state.stack.absIndex(idx)
	if newTop < 0 {
		panic("stack underflow!")
	}

	n := state.stack.top - newTop
	if n > 0 {
		for i := 0; i < n; i++ {
			state.stack.pop()
		}
	} else if n < 0 {
		for i := 0; i > n; i-- {
			state.stack.push(nil)
		}
	}
}
