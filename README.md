# 🧰 Toolkit v2
A lightweight utility library for Go projects, designed to simplify common tasks across JSON handling, file operations, validation, encryption, and temporal types. Built for clarity, composability, and production-grade reliability.

## ✨ Features

### ⚙️ Core Utilities
- [X] Read JSON
- [X] Write JSON
- [X] Error JSON takes an error and optionally a status code, and sends a JSON error message
- [X] Handle Error wraps ErrorJSON and writes to the logger on failure
- [X] Upload a file to a specified directory
- [X] Download a static file
- [X] Get a random string of length n
- [X] Post JSON to a remote service
- [X] Write XML
- [X] Create a directory, including all parent directories, if it does not already exist
- [X] Create a URL safe slug from a string
- [X] Validation utilities
- [X] Encrypt and Decrypt capability

### 🕒 Temporal Types
- [X] DateOnly struct (type) with JSON and SQL value & scan support
- [X] TimeOnly struct (type) with JSON and SQL value & scan support

### 📦 Installation

`go get -u github.com/fouched/toolkit/v2`

### 📅 Temporal Types Philosophy
The DateOnly and TimeOnly types are designed to:
- Represent nullable date/time values
- Integrate seamlessly with JSON, SQL, and domain logic
- Avoid zero-value ambiguity by using *time.Time internally
- Provide comparison, conversion, and formatting utilities
