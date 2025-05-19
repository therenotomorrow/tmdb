import "recipes/build.just"
import "recipes/code.just"

set dotenv-load := false
set shell := ["sh", "-cu"]

BIN := justfile_directory() / "bin"
