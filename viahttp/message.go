package viahttp

import "errors"

type message struct {
	T   string      `json:"topic"`
	Err string      `json:"err"`
	Msg interface{} `json:"msg"`
}

func (m *message) Error() error {
	if m.Err != "" {
		return errors.New(m.Err)
	}
	return nil
}

func (m *message) Topic() string {
	return m.T
}

func FromMessage(topic string, msg interface{}) *message {
	return &message{T: topic, Msg: msg}
}

func (m *message) Message() interface{} {
	return m.Msg
}

func (m *message) Unmarshal(v interface{}) error {
	panic(v)
}

type Response struct {
	Topic string `json:"topic"`
	Err   string `json:"err"`
}
