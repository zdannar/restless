package restless

import ()

type Constructor interface {
    Single() interface{}
    Slice() interface{}
}
