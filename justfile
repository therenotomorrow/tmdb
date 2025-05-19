import "recipes/code.just"
import "recipes/test.just"

set dotenv-load := false
set shell := ["sh", "-cu"]

BIN := justfile_directory() / "bin"

[private]
default:
    just --list
