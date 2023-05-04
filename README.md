# EFINextBoot
Reimplementation of [chengxuncc/booToLinux](https://github.com/chengxuncc/booToLinux) with zero-dependency in mind
Using `bcdedit` to select the next time boot of UEFI: ```bcdedit /set {fwbootmgr} bootsequence {GUID}```

## Download
If you don't want to build by yourself, you can download prebuild binary here: [release](https://github.com/jungin500/efinextboot/releases)

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