package clix

import (
	"flag"
	"time"
)

// @author valor.

// FlagBoolVar calls flag.BoolVar
func (c *Command) FlagBoolVar(p *bool, name string, value bool, usage string) {
	c.goFlagNum += 1
	c.flags().BoolVar(p, name, value, usage)
}

// FlagIntVar calls flag.IntVar
func (c *Command) FlagIntVar(p *int, name string, value int, usage string) {
	c.goFlagNum += 1
	c.flags().IntVar(p, name, value, usage)
}

// FlagInt64Var calls flag.Int64Var
func (c *Command) FlagInt64Var(p *int64, name string, value int64, usage string) {
	c.goFlagNum += 1
	c.flags().Int64Var(p, name, value, usage)
}

// FlagUintVar calls flag.UintVar
func (c *Command) FlagUintVar(p *uint, name string, value uint, usage string) {
	c.goFlagNum += 1
	c.flags().UintVar(p, name, value, usage)
}

// FlagUint64Var calls flag.Uint64Var
func (c *Command) FlagUint64Var(p *uint64, name string, value uint64, usage string) {
	c.goFlagNum += 1
	c.flags().Uint64Var(p, name, value, usage)
}

// FlagStringVar calls flag.StringVar
func (c *Command) FlagStringVar(p *string, name string, value string, usage string) {
	c.goFlagNum += 1
	c.flags().StringVar(p, name, value, usage)
}

// FlagFloat64Var calls flag.Float64Var
func (c *Command) FlagFloat64Var(p *float64, name string, value float64, usage string) {
	c.goFlagNum += 1
	c.flags().Float64Var(p, name, value, usage)
}

// FlagDurationVar calls flag.DurationVar
func (c *Command) FlagDurationVar(p *time.Duration, name string, value time.Duration, usage string) {
	c.goFlagNum += 1
	c.flags().DurationVar(p, name, value, usage)
}

// FlagFunc calls flag.Func
func (c *Command) FlagFunc(name, usage string, fn func(string) error) {
	c.goFlagNum += 1
	c.flags().Func(name, usage, fn)
}

// FlagVar calls flag.Var
func (c *Command) FlagVar(value flag.Value, name string, usage string) {
	c.goFlagNum += 1
	c.flags().Var(value, name, usage)
}
