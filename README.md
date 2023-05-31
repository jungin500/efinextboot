# EFINextBoot
Reimplementation of [chengxuncc/booToLinux](https://github.com/chengxuncc/booToLinux) with zero-dependency in mind
Using `bcdedit` to select the next time boot of UEFI: ```bcdedit /set {fwbootmgr} bootsequence {GUID}```

## What's new?
- Faster launch speed
  - Removed powershell dependency
- Fix broken non-ASCII characters on Multilingual OS
  - Simple fix with CMD builtin `chcp` command

## Download
Download prebuilt binary: [Releases](https://github.com/jungin500/efinextboot/releases)

## Build 
```dos
go build
```

## Build with Administrator privileges
```dos
go build github.com/akavel/rsrc
rsrc.exe -manifest efinextboot.exe.manifest -o efinextboot.syso
go build
```
Double click to run.

## License
EFINextBoot is licensed under the MIT license.
