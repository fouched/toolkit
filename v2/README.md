# Toolkit

A simple example of how to create a reusable Go module with commonly used tools.

The included tools are:

- [X] Read JSON
- [X] Write JSON
- [X] ErrorJSON takes an error and optionally a status code, and sends a JSON error message
- [X] HandleError wraps ErrorJSON and writes to the logger on failure
- [X] Uploads a file to a specified directory
- [X] Download a static file
- [X] Get a random string of length n
- [X] Post JSON to a remote service
- [X] Write XML
- [X] Create a directory, including all parent directories, if it does not already exist
- [X] Create a URL safe slug from a string
- [X] Validation utilities
- [X] Encrypt and Decrypt capability

## Installation

`go get -u github.com/fouched/toolkit/v2`

