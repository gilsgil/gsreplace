# gsreplace

A lightweight and efficient command-line tool for URL parameter fuzzing, designed for web application security testing.

## Overview

gsreplace reads URLs from standard input and generates new URLs by replacing parameter values with a specified fuzzing payload. This tool is useful for testing web applications for vulnerabilities such as XSS, SQL injection, and other input-based attacks.

## Features

- **Standard Mode**: Replaces each parameter individually with the fuzzing payload
- **Clusterbomb Mode**: Replaces all parameters simultaneously with the fuzzing payload
- **Append Mode**: Concatenates the fuzzing payload to existing parameter values
- **Duplicate Avoidance**: Prevents redundant URLs in clusterbomb mode
- **Path Ignoring**: Option to consider URLs with the same parameters but different paths as duplicates

## Installation

```bash
go get github.com/gilsgil/gsreplace

#OR

go install -v github.com/gilsgil/gsreplace@latest
```

Or clone the repository and build:

```bash
git clone https://github.com/gilsgil/gsreplace.git
cd gsreplace
go build
```

## Usage

```bash
cat urls.txt | gsreplace [options] <fuzz_word>
```

### Options:

- `-c`: Activates clusterbomb mode (replaces all parameters simultaneously)
- `-a`: Activates append mode (adds the fuzz word to existing parameter values)
- `-ignore-path`: Ignores path when considering duplicates (only in clusterbomb mode)

### Examples:

**Standard fuzzing:**
```bash
echo "https://example.com/?param1=value1&param2=value2" | gsreplace XSS
```
Output:
```
https://example.com/?param1=XSS&param2=value2
https://example.com/?param1=value1&param2=XSS
```

**Clusterbomb mode:**
```bash
echo "https://example.com/?param1=value1&param2=value2" | gsreplace -c XSS
```
Output:
```
https://example.com/?param1=XSS&param2=XSS
```

**Append mode:**
```bash
echo "https://example.com/?param1=value1&param2=value2" | gsreplace -a "<script>"
```
Output:
```
https://example.com/?param1=value1<script>&param2=value2
https://example.com/?param1=value1&param2=value2<script>
```

## Integrating with Other Tools

gsreplace follows Unix philosophy by reading from stdin and writing to stdout, making it easy to integrate with other tools:

```bash
cat urls.txt | gsreplace XSS | httpx -silent | nuclei -t xss.yaml
```

## Use Cases

- Testing for XSS vulnerabilities by replacing parameters with XSS payloads
- Checking for SQL injection by using SQL-specific fuzzing words
- Identifying parameter pollution issues
- Batch testing of multiple parameters across many URLs

## License

MIT License

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
