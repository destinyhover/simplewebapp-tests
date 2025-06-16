# simplewebapp-tests

This is a personal fork of [Jon Bodner](https://github.com/learning-go-book-2e/simplewebapp)'s Go web application, extended for educational purposes.

## 🔧 What's added

- Unit tests for:
  - `parser.go`
  - `WriteData.go`
  - `DataProcessor.go`
- Basic fuzz test data (removed from tracked files)
- `.gitignore` cleanup for build artifacts and coverage files

## 🚀 How to run tests

```bash
go test ./...
