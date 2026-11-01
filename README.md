# tetoris
tetris, but in as a static binary and TUI!!!!1

Features:
- [x] all features related to moving a block, including rotating
- [x] all shapes of the original
- [x] basic score system
- [x] pausing
- [x] skip spamming down with an "instadown" key that teleports the block as low as it can go
- [x] deterministic games (seeded RNG)

## installation
### github releases
To install *tetoris*, download a binary below and place it somewhere within your $PATH, or execute it directly.

You need to grab a different binary depending on your system and architecture:
- **tetoris-linux-amd64** This is the version you need on Linux, choose if uncertain of architecture.
- **tetoris-darwin-arm64** This is the version you need on Apple Silicon macs (M1 and later)
- **tetoris-windows-amd64** This is the version you need on Windows machines, choose if uncertain of architecture.

- **tetoris-darwin-amd64** This is the version you need on Intel macs (everything before the M1)
- **tetoris-linux-arm64** This is the version for use on arm64 Linux eg. Raspberry Pis, Asahi Linux. Use linux-amd64 if your not sure.
- **tetoris-windows-arm64** Probably don't need this unless you are certain you do.

On macOS, you may get an error regarding Gatekeeper. You can get rid of that by running the following:
```sh
sudo xattr -rd com.apple.quarantine BINARY_PATH_HERE
```

On Linux and/or macOS, you may have to add the "execute" permission to the binary:
```sh
chmod +x BINARY_PATH_HERE
```

Then, you should be able to run it with `./tetoris-darwin-arm64` in the same directory as the binary, replacing the name with the one you downloaded.

### go install
Alternatively, you can install *tetoris* with `go install github.com/mifumofu/tetoris`, this will install a binary called `tetoris` to the bin/ directory in your $GOPATH.

You can run it with `tetoris` if you added `$GOPATH/bin` to your path.

## license
This code is licensed under the terms of the GPLv3
