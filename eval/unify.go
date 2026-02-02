package eval

import (
	"fmt"
	"strings"

	"github.com/twolodzko/kanren/types"
)

// The variable-value mapping
type KeyVal struct {
	key types.Variable
	val any
}

// The alist holding key-value pairs for the unification results
// (see Byrd, 2009, p. 25)
type Stream struct {
	list []KeyVal
}

func NewStream() *Stream {
	return &Stream{make([]KeyVal, 0)}
}

// Unify two values, return status (see Byrd, 2009, p. 29)
func (s *Stream) unify(u, v any) bool {
	if Debug {
		fmt.Printf(" ↪ unify: (== %v %v)\n", types.ToString(u), types.ToString(v))
		fmt.Printf("   subst: %v\n", s)
	}
	u = s.walk(u)
	v = s.walk(v)
	if u == v {
		return true
	}
	if u, ok := u.(types.Variable); ok {
		return s.extend(u, v)
	}
	if v, ok := v.(types.Variable); ok {
		return s.extend(v, u)
	}
	if u, ok := u.(types.Pair); ok {
		if v, ok := v.(types.Pair); ok {
			if !s.unify(u.This, v.This) {
				return false
			}
			return s.unify(u.Next, v.Next)
		}
	}
	return false
}

func (s Stream) reify(v any) any {
	v = s.deepWalk(v)
	fresh := NewStream()
	fresh.reifyStream(v)
	return fresh.deepWalk(v)
}

func (s *Stream) reifyStream(v any) bool {
	switch v := s.walk(v).(type) {
	case types.Variable:
		free := types.Free(s.len())
		if Debug {
			fmt.Printf(" ↪ reify: %v = %v\n", types.ToString(v), types.ToString(free))
			fmt.Printf("   subst: %v\n", s)
		}
		return s.extend(v, free)
	case types.Pair:
		var head any = v
		for head != nil {
			switch p := head.(type) {
			case types.Pair:
				if !s.reifyStream(p.This) {
					return false
				}
				head = p.Next
			default:
				return s.reifyStream(head)
			}
		}
	}
	return true
}

// Recursively get the value for the key (see Byrd, 2009, p. 27)
func (s Stream) walk(v any) any {
	if v, ok := v.(types.Variable); ok {
		val, ok := s.get(v)
		if !ok {
			return v
		}
		if v == val {
			// encountered a birth record
			return v
		}
		return s.walk(val)
	}
	return v
}

func (s Stream) deepWalk(v any) any {
	v = s.walk(v)
	if v, ok := v.(types.Pair); ok {
		v := v.Map(func(x any) any {
			return s.deepWalk(x)
		})
		return types.Cons(v...)
	}
	return v
}

// Check for circular references between keys and values (see Byrd, 2009, p. 28)
func (s Stream) occurs(u any, v any) bool {
	v = s.walk(v)
	switch val := v.(type) {
	case types.Variable:
		return u == val
	case types.Pair:
		return val.Any(func(v any) bool {
			return s.occurs(u, v)
		})
	}
	return false
}

func (s *Stream) extend(u types.Variable, v any) bool {
	// occurs check
	// if s.occurs(u, v) {
	// 	return false
	// }
	s.list = append(s.list, KeyVal{u, v})
	return true
}

func (s Stream) get(v types.Variable) (any, bool) {
	for i := s.len() - 1; i >= 0; i-- {
		if s.list[i].key == v {
			return s.list[i].val, true
		}
	}
	return nil, false
}

func (s *Stream) birthRecord(key types.Variable) {
	s.list = append(s.list, KeyVal{key, key})
}

func (s *Stream) keep(size int) {
	s.list = s.list[:size]
}

func (s Stream) len() int {
	return len(s.list)
}

func (s Stream) String() string {
	var acc []string
	for _, kv := range s.list {
		acc = append(acc, kv.String())
	}
	return fmt.Sprintf("[%s]", strings.Join(acc, " "))
}

func (k KeyVal) String() string {
	return fmt.Sprintf("%s:%s", types.ToString(k.key), types.ToString(k.val))
}
