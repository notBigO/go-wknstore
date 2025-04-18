# WebKnot Numbers (WKN)

A minimal REPL-based database that works with integer arrays.

## Overview

WKN is a lightweight database tool designed specifically for storing and manipulating integer arrays. It provides a simple REPL (Read-Eval-Print Loop) interface to create, manage, and perform operations on numerical data.

## Installation

```bash
# Clone the repository
git clone https://github.com/notBigO/wkn.git

# Build the application
cd wkn
go build -o wkn
```

## Usage

### Starting a New Database

```bash
./wkn new
```

This command creates a new `.wkn` database file in the current directory and starts the REPL.

### Loading an Existing Database

```bash
./wkn --db-path ./path_to_file.wkn
```

This command loads an existing `.wkn` file and starts the REPL with the loaded data.

## REPL Commands

Once in the REPL, you can use the following commands:

### Create a New Array

```
wkn> new array_name 1 2 3 4 5
```

Creates a new array with the given name and values.

### Display Arrays

```
wkn> show
```

Displays all arrays in the database.

```
wkn> show array_name
```

Displays a specific array.

### Merge Arrays

```
wkn> merge target_array source_array
```

Appends the values from `source_array` to `target_array`.

### Power Operation

```
wkn> pow arrayA.indexA arrayB.indexB
```

Calculates (arrayA[indexA] ^ arrayB[indexB]) % (10^9 + 7) using binary exponentiation.

### Delete an Array

```
wkn> del array_name
```

Removes the specified array from the database.

### Exit

```
wkn> exit
```

Exits the REPL.

## Examples

```
# Create two new arrays
wkn> new bases 2 3 5 7
CREATED (4)

wkn> new exponents 2 3 2 4
CREATED (4)

# View all arrays
wkn> show
bases: [2 3 5 7]
exponents: [2 3 2 4]

# Calculate 3^3 (bases[1]^exponents[1])
wkn> pow bases.1 exponents.1
27

# Merge arrays
wkn> merge bases exponents
MERGED

# View updated array
wkn> show bases
[2 3 5 7 2 3 2 4]

# Delete an array
wkn> del exponents
DELETED
```

## File Format

WKN stores data in a simple JSON format in a `.wkn` file. The database uses a lock file (`.wkn.lock`) for basic concurrency control.

## Features

- Simple, command-line based interface
- Persistent storage of integer arrays
- Basic concurrency control with file locking
- Mathematical operations (currently power with modulo)
- Array manipulation (creation, deletion, merging)
