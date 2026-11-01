#!/bin/sh
# Build script to generate binaries for a list of platforms

set -e

if [[ $1 = "--dry-run" ]]; then
  dryrun=1
fi

platforms=("windows" "darwin" "linux")
archs=("amd64" "arm64")
outdir="out"
name=$(basename "$PWD")

# clear out dir
rm -rf "$outdir" # f doesn't error if it doesn't exist
mkdir $outdir

# build binaries
for platform in "${platforms[@]}"; do
  for arch in "${archs[@]}"; do
    export GOOS=$platform
    export GOARCH=$arch
    export CGO_ENABLED=0

    outpath="out/$name-$platform-$arch"
    cmd="go build -o '$outpath' -ldflags='-s -w' -trimpath"
    echo "=> $cmd"
    
    if [ "$dryrun" = "" ]; then
      go build -o "$outpath" -ldflags="-s -w" -trimpath
      file "$outpath"
      du -h "$outpath"
    fi
  done
done
