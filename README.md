# Markov Chain Text Generator

A command-line Markov Chain text generator written in Go.

It reads input text from `stdin`, builds a word-based Markov chain using a configurable prefix length, and generates new text that mimics the style of the input.

## Features

- Builds a Markov Chain of configurable prefix length (`-l`)
- Allows setting a custom starting prefix (`-p`)
- Controls the maximum number of output words (`-w`)
- Simple and efficient implementation using only Go's standard library
- Handles input from files via `stdin`

## Usage

```bash
cat the_great_gatsby.txt | ./markovchain [options]
```

```
Option	Description
-w <N>	Number of words to generate (default: 100)
-l <N>	Prefix length (number of words in a key). Must be between 1 and 5. Default: 2
-p <S>	Starting prefix (must match prefix length)
--help	Show usage information
```

### Example
```
cat the_great_gatsby.txt | ./markovchain -w 10 -l 3 -p "to something funny"
```

### Output:
```
to something funny happened in the kitchen and I looked
```

## Building
```
git clone https://github.com/yourusername/markov-chain.git
cd markov-chain
go build -o markovchain .
```

## Project Structure
```
.
├── markov
│   ├── chain.go
│   └── flag.go
├── go.mod
├── main.go
├── markovchain
├── README.md
└── the_great_gatsby.txt
```

## Created by `zaaripzha`