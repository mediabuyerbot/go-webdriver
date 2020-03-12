package protocol

type Options interface {
	Proxy() *Proxy
	FirstMatch() []O
	AlwaysMatch() O
}

type O map[string]interface{}

func MakeOptions() O {
	return make(O)
}

func (o O) Has(key string) bool {
	_, ok := o[key]
	return ok
}

func (o O) HasNot(key string) bool {
	_, ok := o[key]
	return !ok
}

func (o O) Set(k string, v interface{}) O {
	o[k] = v
	return o
}

func (o O) GetString(key string) (s string) {
	v, ok := o[key]
	if !ok {
		return
	}
	s, ok = v.(string)
	if !ok {
		return
	}
	return s
}

func (o O) GetBytes(key string) (s []byte) {
	v, ok := o[key]
	if !ok {
		return nil
	}
	s, ok = v.([]byte)
	if !ok {
		return nil
	}
	return s
}

func (o O) GetStringSlice(key string) (s []string) {
	v, ok := o[key]
	if !ok {
		return nil
	}
	s, ok = v.([]string)
	if !ok {
		return nil
	}
	return s
}

func (o O) GetOpts(key string) O {
	v, ok := o[key]
	if !ok {
		return nil
	}
	switch x := v.(type) {
	case map[string]interface{}:
		return O(x)
	case O:
		return x
	default:
		return nil
	}
}

func (o O) GetBool(key string) (b bool) {
	v, ok := o[key]
	if !ok {
		return
	}
	b, ok = v.(bool)
	if !ok {
		return
	}
	return b
}

func (o O) GetUint(key string) (i uint) {
	v, ok := o[key]
	if !ok {
		return
	}
	i, ok = v.(uint)
	if !ok {
		return
	}
	return i
}

func (o O) GetInt(key string) (i int) {
	v, ok := o[key]
	if !ok {
		return
	}
	i, ok = v.(int)
	if !ok {
		return
	}
	return i
}

func (o O) GetFloat(key string) (f float64) {
	v, ok := o[key]
	if !ok {
		return
	}
	f, ok = v.(float64)
	if !ok {
		return
	}
	return f
}
