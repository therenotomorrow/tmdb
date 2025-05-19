BUILD_PATH := BIN / "tmdb"
BUILD_MAIN := invocation_directory() / "cmd" / "tmdb" / "main.go"

[doc('Build the tmdb binary with an optional `token` embedded')]
build token="":
    go build -o {{BUILD_PATH}} -ldflags "-X 'main.TMDBToken={{token}}'" {{BUILD_MAIN}}
