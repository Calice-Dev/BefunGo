package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"strings"
)

const MAX_SIZE_X = 128
const MAX_SIZE_Y = 128

type stackNode struct {
	data byte
	next *stackNode
}

type stack struct {
	top *stackNode
}

func (stack *stack) isEmpty() bool {
	return stack.top == nil
}

func (stack *stack) Push(data byte) {
	newNode := new(stackNode)
	newNode.data = data
	newNode.next = stack.top
	stack.top = newNode
}

func main() {
	c, _ := readInput("code")
	Interpret(c)
}

func stackCreate() *stack {
	newStack := new(stack)
	newStack.top = nil
	return newStack
}

func (stack *stack) Pop() byte {
	if stack.isEmpty() {
		return 0
	}
	node := stack.top
	v := node.data
	stack.top = node.next
	return v
}

func (stack *stack) Peek() byte {
	if stack.isEmpty() {
		return 0
	}
	return stack.top.data
}

func readInput(fileName string) ([]string, error) {
	// reads the file as bytes in inputText, stores error in err
	inputText, err := ioutil.ReadFile(fileName)
	// if there are any errors with reading the file returns
	if err != nil {
		errorMessage := fmt.Sprintf("file not located: %v", fileName)
		return []string{}, errors.New(errorMessage)
	}
	lines := strings.Split(string(inputText), "\n")
	for i, l := range lines {
		if len(l) < MAX_SIZE_X {
			n := MAX_SIZE_Y - len(l)
			s := strings.Repeat(" ", n)
			lines[i] = l + s
		}
	}
	return lines, nil
}

func getCommandAt(x, y int, code []string) byte {
	if y >= len(code) {
		y = 0
	}
	if x >= len(code[y]) {
		x = 0
	}
	return code[y][x]
}

// https://stackoverflow.com/questions/24893624/how-to-replace-a-letter-at-a-specific-index-in-a-string-in-go
func replaceAtIndex(in string, r byte, i int) string {
	out := []byte(in)
	out[i] = r
	s := string(out)
	if len(s) < MAX_SIZE_X {
		n := MAX_SIZE_X - len(s)
		e := strings.Repeat(" ", n)
		s = s + e
	}
	return s
}

func runCommand(command byte, stack *stack, dir, x, y *int, code []string, stringMode *bool) {
	if command == '"' {
		*stringMode = !(*stringMode)
		return
	}
	if *stringMode {
		stack.Push(command)
		return
	}
	switch command {
	case ' ':
		return
	case '@':
		return
	case '>':
		*dir = 0
		return
	case '<':
		*dir = 1
		return
	case 'v':
		*dir = 2
		return
	case '^':
		*dir = 3
		return
	case '?':
		*dir = rand.Int() % 4
		return
	case '+':
		a, b := stack.Pop(), stack.Pop()
		stack.Push(a + b)
		return
	case '-':
		a, b := stack.Pop(), stack.Pop()
		stack.Push(b - a)
		return
	case '*':
		a, b := stack.Pop(), stack.Pop()
		stack.Push(b * a)
		return
	case '/':
		a, b := stack.Pop(), stack.Pop()
		stack.Push(b / a)
		return
	case '%':
		a, b := stack.Pop(), stack.Pop()
		stack.Push(b % a)
		return
	case '!':
		a := stack.Pop()
		if int(a) == 0 {
			stack.Push(byte(1))
		} else {
			stack.Push(byte(0))
		}
		return
	case '`':
		a, b := stack.Pop(), stack.Pop()
		if b > a {
			stack.Push(byte(1))
		} else {
			stack.Push(byte(0))
		}
		return
	case ':':
		a := stack.Peek()
		stack.Push(a)
		return
	case '$':
		stack.Pop()
		return
	case '_':
		a := stack.Pop()
		if a == 0 {
			*dir = 0
			return
		}
		*dir = 1
		return
	case '|':
		a := stack.Pop()
		if a == 0 {
			*dir = 2
			return
		}
		*dir = 3
		return
	case '\\':
		a, b := stack.Pop(), stack.Pop()
		stack.Push(a)
		stack.Push(b)
		return
	case '#':
		updatePosition(x, y, *dir)
		return
	case 'p':
		y, x, v := stack.Pop(), stack.Pop(), stack.Pop()
		code[y] = replaceAtIndex(code[y], v, int(x))
		return
	case 'g':
		y, x := stack.Pop(), stack.Pop()
		stack.Push(code[y][x])
		return
	case ',':
		c := stack.Pop()
		fmt.Printf("%c", c)
		return
	case '.':
		c := stack.Pop()
		fmt.Printf("%d ", c)
		return
	}
	if command >= '0' && command <= '9' {
		command -= '0'
		stack.Push(command)
	}
}

func updatePosition(x, y *int, dir int) {
	switch dir {
	case 0:
		*x++
		*x = *x % MAX_SIZE_X
	case 1:
		*x--
		if *x < 0 {
			*x = MAX_SIZE_X - 1
		}
	case 2:
		*y++
		*y = *y % MAX_SIZE_Y
	case 3:
		*y--
		if *y < 0 {
			*y = MAX_SIZE_Y - 1
		}
	}
}

func PrintCode(code []string) {
	fmt.Println("------------")
	for i, v := range code {
		fmt.Println(i, " ", v)
	}
}

func PrintStack(stack *stack) {
	fmt.Println("Current stack:\n------")
	node := stack.top
	for node != nil {
		fmt.Println(node.data)
		node = node.next
	}
}

func Interpret(code []string) {
	stack := stackCreate()
	x, y := 0, 0
	currentDir := 0
	currentCommand := getCommandAt(x, y, code)
	stringMode := false
	for currentCommand != '@' {
		runCommand(currentCommand, stack, &currentDir, &x, &y, code, &stringMode)
		updatePosition(&x, &y, currentDir)
		currentCommand = getCommandAt(x, y, code)
		//PrintCode(l)
	}
}
