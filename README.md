# What is Monkey?

Monkey is an interpreted, functional programming language written in Go. It was created for learning purposes and is perfect for experimenting with interpreters, compilers, and language design.

# Getting Started

This is a short Monkey tutorial that should take no more than 10 minutes. It assumes you already have the Monkey REPL installed.

## Monkey REPL

The Monkey REPL evaluates any Monkey expressions you enter and prints the result. Playing with Monkey in an interactive environment is a great way to explore what the language can do.

Try this:

```
>> "Hello World"
Hello World
```

What happened? Did we really just write the world’s shortest “Hello World” program?

Not exactly—the second line is simply the REPL showing the result of the expression. If you want to *print* “Hello World”, you need to use `puts`:

```
>> puts("Hello World")
Hello World
null
```

`puts` is Monkey’s basic print command. The `null` you see afterward is the return value of `puts`, since it always returns Monkey’s “nothing here” value.

---

## Your Free Calculator

The Monkey REPL already gives you a functional little calculator:

```
>> 3 + 2
5
>> 3 - 2
1
>> 3 * 2
6
>> 3 / 2
1
```

> [!IMPORTANT]
> Monkey does **not** support floating-point arithmetic. After all, no Monkey needs to store 3.14159265359 bananas!

---

## Monkey Do What Monkey Say

Let’s reduce typing by writing a function:

```
>> let sayHi = fn() { puts("Hello World"); };
```

This defines a function named `sayHi`. The code inside the braces is the function body.

Call it like this:

```
>> sayHi()
Hello World
null
```

Want to greet someone by name? Just add a parameter:

```
>> let sayHi = fn(name) { puts("Hello " + name); };
>> sayHi("Vinh")
Hello Vinh
null
```

> [!IMPORTANT]
> Monkey only supports string concatenation with the `+` operator.

---

## Arrays

You can declare and access array elements just like in many other languages. Monkey arrays can hold multiple types:

```
>> let arr = [1, "two", 3 == (2 + 1)]
>> arr[0]
1
>> arr[1]
two
>> arr[2]
true
```

---

## Hashes

Monkey supports hash literals with keys of different types:

```
>> let hash = {"foo": 1, 2: "bar", true: "baz"}
>> hash["foo"]
1
>> hash[2]
bar
>> hash[true]
baz
```

---

## If–Else Conditions

Monkey supports basic `if`–`else` expressions:

```
>> if (true) { puts("hello world"); } else { puts("goodbye world"); };
hello world
null
```

What about `else if`? Monkey doesn’t officially support it, but you can nest `if` blocks:

```
>> if (false) { puts("hello world") } else { if (true) { puts("monkey hello") } else { puts("goodbye world") } };
monkey hello
null
```

---

## Built-in Functions

### `len`

Returns the length of a string:

```
>> len("hello")
5
```

### `first`

Returns the first element of an array:

```
>> first([1, 2, 3])
1
```

### `last`

Returns the last element of an array:

```
>> last([1, 2, 3])
3
```

### `rest`

Returns a new array containing all elements *except* the first:

```
>> rest([1, 2, 3])
[2, 3]
```

### `push`

Returns a new array with the element appended. The original array is unchanged:

```
>> let arr = [1, 2, 3]
>> push(arr, 4)
[1, 2, 3, 4]
>> arr
[1, 2, 3]
```

# Additional Information
## Based on “Writing an Interpreter in Go”

This implementation of Monkey is heavily inspired by Thorsten Ball’s excellent book Writing an Interpreter in Go, which provides a step-by-step guide to building Monkey from scratch. The book serves as the foundational blueprint for the language’s structure and behavior.

## Built with Test-Driven Development (TDD)

Monkey was developed entirely using Test-Driven Development (TDD).
TDD played a crucial role in navigating the inherent complexity of designing and implementing an interpreted language. Writing tests first provided clarity, avoided unnecessary detours, and ensured that each new feature integrated cleanly with the rest of the system.

# Maintainer

Vinh Bui