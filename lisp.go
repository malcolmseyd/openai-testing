// Write a Lisp interpreter in Go

package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type LispVal struct {
	Type  string
	Value interface{}
}

var env map[string]LispVal

func main() {
	env = make(map[string]LispVal)
	env["+"] = LispVal{"func", func(args []LispVal) LispVal {
		var sum int
		for _, arg := range args {
			sum += arg.Value.(int)
		}
		return LispVal{"int", sum}
	}}
	env["-"] = LispVal{"func", func(args []LispVal) LispVal {
		var diff int
		for i, arg := range args {
			if i == 0 {
				diff = arg.Value.(int)
			} else {
				diff -= arg.Value.(int)
			}
		}
		return LispVal{"int", diff}
	}}
	env["*"] = LispVal{"func", func(args []LispVal) LispVal {
		var prod int
		for i, arg := range args {
			if i == 0 {
				prod = arg.Value.(int)
			} else {
				prod *= arg.Value.(int)
			}
		}
		return LispVal{"int", prod}
	}}
	env["/"] = LispVal{"func", func(args []LispVal) LispVal {
		var quot int
		for i, arg := range args {
			if i == 0 {
				quot = arg.Value.(int)
			} else {
				quot /= arg.Value.(int)
			}
		}
		return LispVal{"int", quot}
	}}
	env["="] = LispVal{"func", func(args []LispVal) LispVal {
		var eq bool
		var prev int
		for i, arg := range args {
			if i != 0 {
				eq = prev == arg.Value.(int)
			}
			prev = arg.Value.(int)
		}
		return LispVal{"bool", eq}
	}}
	env["<"] = LispVal{"func", func(args []LispVal) LispVal {
		var lt bool
		var prev int
		for i, arg := range args {
			if i != 0 {
				lt = prev < arg.Value.(int)
			}
			prev = arg.Value.(int)
		}
		return LispVal{"bool", lt}
	}}
	env[">"] = LispVal{"func", func(args []LispVal) LispVal {
		var gt bool
		var prev int
		for i, arg := range args {
			if i != 0 {
				gt = prev > arg.Value.(int)
			}
			prev = arg.Value.(int)
		}
		return LispVal{"bool", gt}
	}}
	env["<="] = LispVal{"func", func(args []LispVal) LispVal {
		var lte bool
		var prev int
		for i, arg := range args {
			if i != 0 {
				lte = prev <= arg.Value.(int)
			}
			prev = arg.Value.(int)
		}
		return LispVal{"bool", lte}
	}}
	env[">="] = LispVal{"func", func(args []LispVal) LispVal {
		var gte bool
		var prev int
		for i, arg := range args {
			if i != 0 {
				gte = prev >= arg.Value.(int)
			}
			prev = arg.Value.(int)
		}
		return LispVal{"bool", gte}
	}}
	env["cons"] = LispVal{"func", func(args []LispVal) LispVal {
		return LispVal{"list", append([]LispVal{args[0]}, args[1].Value.([]LispVal)...)}
	}}
	env["car"] = LispVal{"func", func(args []LispVal) LispVal {
		return args[0].Value.([]LispVal)[0]
	}}
	env["cdr"] = LispVal{"func", func(args []LispVal) LispVal {
		return LispVal{"list", args[0].Value.([]LispVal)[1:]}
	}}
	env["list"] = LispVal{"func", func(args []LispVal) LispVal {
		return LispVal{"list", args}
	}}
	env["null?"] = LispVal{"func", func(args []LispVal) LispVal {
		return LispVal{"bool", len(args[0].Value.([]LispVal)) == 0}
	}}
	env["if"] = LispVal{"func", func(args []LispVal) LispVal {
		if args[0].Value.(bool) {
			return args[1]
		} else {
			return args[2]
		}
	}}
	env["def"] = LispVal{"func", func(args []LispVal) LispVal {
		env[args[0].Value.(string)] = args[1]
		return args[1]
	}}
	env["lambda"] = LispVal{"func", func(args []LispVal) LispVal {
		return LispVal{"func", func(args2 []LispVal) LispVal {
			newEnv := make(map[string]LispVal)
			for k, v := range env {
				newEnv[k] = v
			}
			for i, arg := range args[0].Value.([]LispVal) {
				newEnv[arg.Value.(string)] = args2[i]
			}
			return eval(args[1], newEnv)
		}}
	}}
	env["print"] = LispVal{"func", func(args []LispVal) LispVal {
		fmt.Println(args[0].Value)
		return LispVal{"nil", nil}
	}}
	env["read"] = LispVal{"func", func(args []LispVal) LispVal {
		reader := bufio.NewReader(os.Stdin)
		text, _ := reader.ReadString('\n')
		return LispVal{"string", text}
	}}
	env["eval"] = LispVal{"func", func(args []LispVal) LispVal {
		return eval(args[0], env)
	}}
	env["exit"] = LispVal{"func", func(args []LispVal) LispVal {
		os.Exit(0)
		return LispVal{"nil", nil}
	}}
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print(">> ")
		text, _ := reader.ReadString('\n')
		fmt.Println(eval(read(text), env))
	}
}

func read(text string) LispVal {
	v, _ := parse(text)
	return v
}

func parse(text string) (LispVal, string) {
	text = strings.TrimSpace(text)
	if text[0] == '(' {
		return parseList(text)
	} else if text[0] == '"' {
		return parseString(text)
	} else if text[0] == '\'' {
		return parseQuote(text)
	} else if text[0] == '-' || (text[0] >= '0' && text[0] <= '9') {
		return parseInt(text)
	} else {
		return parseSymbol(text)
	}
}

func parseList(text string) (LispVal, string) {
	text = text[1:]
	var list []LispVal
	for len(text) > 0 {
		if text[0] == ')' {
			text = text[1:]
			break
		}
		var elem LispVal
		elem, text = parse(text)
		list = append(list, elem)
	}
	return LispVal{"list", list}, text
}

func parseString(text string) (LispVal, string) {
	text = text[1:]
	var str string
	for len(text) > 0 {
		if text[0] == '"' {
			text = text[1:]
			break
		}
		str += string(text[0])
		text = text[1:]
	}
	return LispVal{"string", str}, text
}

func parseQuote(text string) (LispVal, string) {
	text = text[1:]
	val, text := parse(text)
	return LispVal{"quote", val}, text
}

func parseInt(text string) (LispVal, string) {
	var num int
	var neg bool
	if text[0] == '-' {
		neg = true
		text = text[1:]
	}
	for len(text) > 0 {
		if text[0] == ' ' {
			text = text[1:]
			break
		}
		num *= 10
		num += int(text[0] - '0')
		text = text[1:]
	}
	if neg {
		num = -num
	}
	return LispVal{"int", num}, text
}

func parseSymbol(text string) (LispVal, string) {
	var sym string
	for len(text) > 0 {
		if text[0] == ' ' {
			text = text[1:]
			break
		}
		sym += string(text[0])
		text = text[1:]
	}
	return LispVal{"symbol", sym}, text
}

func eval(val LispVal, env map[string]LispVal) LispVal {
	if val.Type == "int" {
		return val
	} else if val.Type == "string" {
		return val
	} else if val.Type == "bool" {
		return val
	} else if val.Type == "quote" {
		return val
	} else if val.Type == "symbol" {
		return env[val.Value.(string)]
	} else if val.Type == "list" {
		if val.Value.([]LispVal)[0].Type == "symbol" {
			if val.Value.([]LispVal)[0].Value.(string) == "quote" {
				return val.Value.([]LispVal)[1]
			} else if val.Value.([]LispVal)[0].Value.(string) == "if" {
				if eval(val.Value.([]LispVal)[1], env).Value.(bool) {
					return eval(val.Value.([]LispVal)[2], env)
				} else {
					return eval(val.Value.([]LispVal)[3], env)
				}
			} else if val.Value.([]LispVal)[0].Value.(string) == "def" {
				env[val.Value.([]LispVal)[1].Value.(string)] = eval(val.Value.([]LispVal)[2], env)
				return env[val.Value.([]LispVal)[1].Value.(string)]
			} else if val.Value.([]LispVal)[0].Value.(string) == "lambda" {
				return LispVal{"func", func(args []LispVal) LispVal {
					newEnv := make(map[string]LispVal)
					for k, v := range env {
						newEnv[k] = v
					}
					for i, arg := range val.Value.([]LispVal)[1].Value.([]LispVal) {
						newEnv[arg.Value.(string)] = args[i]
					}
					return eval(val.Value.([]LispVal)[2], newEnv)
				}}
			} else {
				var args []LispVal
				for _, arg := range val.Value.([]LispVal)[1:] {
					args = append(args, eval(arg, env))
				}
				return env[val.Value.([]LispVal)[0].Value.(string)].Value.(func([]LispVal) LispVal)(args)
			}
		} else {
			var list []LispVal
			for _, arg := range val.Value.([]LispVal) {
				list = append(list, eval(arg, env))
			}
			return LispVal{"list", list}
		}
	} else if val.Type == "func" {
		return val
	} else {
		return LispVal{"nil", nil}
	}
}

func (val LispVal) String() string {
	if val.Type == "int" {
		return strconv.Itoa(val.Value.(int))
	} else if val.Type == "string" {
		return val.Value.(string)
	} else if val.Type == "bool" {
		if val.Value.(bool) {
			return "true"
		} else {
			return "false"
		}
	} else if val.Type == "quote" {
		return "'" + val.Value.(LispVal).String()
	} else if val.Type == "symbol" {
		return val.Value.(string)
	} else if val.Type == "list" {
		var str string
		for _, arg := range val.Value.([]LispVal) {
			str += arg.String() + " "
		}
		return "(" + str[:len(str)-1] + ")"
	} else if val.Type == "func" {
		return "<function>"
	} else {
		return "nil"
	}
}
