# fileprocess

**fileprocess** is a Go program that reads a file line by line, where each line contains a 16-bit unsigned integer representing milliseconds. Each line is processed concurrently with a configurable limit on the number of parallel operations. The program delays for the specified milliseconds, then prints the value.

## Features

- Reads a file containing numbers (each on a separate line).
- Processes lines concurrently, up to a specified limit.
- Delays for N milliseconds (where N is the value in the line) before printing the value.
- Preserves output timing relative to delay values, but not order of lines.

## Usage

```bash
fileprocess <NPROC> <FILE>
````

* `<NPROC>`: The maximum number of concurrent operations (goroutines).
* `<FILE>`: Path to the input file containing numbers (one per line).

## Example

Given a file `reverse-duplicates.txt`:

```
3000
3000
2000
2000
1000
1000
0
0
```

Running:

```bash
fileprocess 8 reverse-duplicates.txt
```

Produces output:

```
0
0
1000
1000
2000
2000
3000
3000
```

Each number is printed after its corresponding delay. The program completes in a little over 3 seconds, demonstrating concurrent execution.

## Installation

To build from source:

```bash
go build -o fileprocess main.go
```

Make sure you have [Go installed](https://golang.org/dl/).

## Notes

* Input must be valid integers in the range of a 16-bit unsigned number (0–65535).
* The program does not perform extensive error handling. Invalid input or arguments will cause it to exit with an error.

## License

This project is open source and available under the MIT License.
