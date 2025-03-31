# FixedPoint Library

The `FixedPoint` library provides a robust implementation for fixed-point arithmetic in Go. It includes a `Context` type for managing
precision, rounding, and signal handling, as well as operations for arithmetic, comparison, and parsing.

## Features

- **Precision Control**: Supports configurable precision for fixed-point values.
- **Rounding Modes**: Includes multiple rounding modes (e.g., `RoundHalfUp`, `RoundHalfEven`).
- **Signal Handling**: Detects and handles errors like overflow, invalid operations, and conversion syntax issues.
- **Arithmetic Operations**: Provides addition, subtraction, multiplication, division, negation, and absolute value operations.
- **Comparison**: Supports equality and relational comparisons.
- **Parsing**: Converts strings into fixed-point values, including special values like `NaN` and `Infinity`.

## Installation

To use the library, add it to your Go module:

```bash
go get github.com/trippwill/go-currency/fixedpoint
```

## Usage

### Creating a Context

You can create a `Context` with default values or customize it:

```go
ctx := fixedpoint.BasicContext() // Basic context with default precision and rounding
extendedCtx := fixedpoint.ExtendedContext() // Extended context with maximum precision
customCtx, err := fixedpoint.NewContext(10, fixedpoint.RoundHalfUp, fixedpoint.SignalOverflow)
if err != nil {
    log.Fatal(err)
}
```

### Parsing FixedPoint Values

```go
value := ctx.Parse("123.45")
fmt.Println(value)
```

### Arithmetic Operations

```go
result := ctx.Add(a, b)
result = ctx.Sub(a, b)
result = ctx.Mul(a, b)
result = ctx.Div(a, b)
```

### Comparison

```go
if ctx.LessThan(a, b) {
    fmt.Println("a is less than b")
}
```

### Signal Handling

```go
result := ctx.Must(ctx.Add(a, b)) // Panics if traps are triggered
result = ctx.Trap(func(ctx *fixedpoint.Context, a fixedpoint.FixedPoint) fixedpoint.FixedPoint {
    fmt.Println("Handling trap")
    return a
}, a)
```

### Aggregation

```go
result := fixedpoint.All(ctx.Add, a, b, c, d)
```

## Contributing

Contributions are welcome! Please open an issue or submit a pull request.

## License

This project is licensed under the MIT License. See the `LICENSE` file for details.
