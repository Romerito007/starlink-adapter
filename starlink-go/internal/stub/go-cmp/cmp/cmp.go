package cmp

type Option interface{}

type Options []Option

func Diff(_, _ any, _ ...Option) string { return "" }

func Equal(x, y any, _ ...Option) bool { return x == y }
