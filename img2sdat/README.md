# img2sdat
Convert filesystem ext4 image (.img) into Android sparse data image (.dat)

## Requirements
This binary requires Python 2.7 or newer installed on your system.

It currently supports Windows x86/x64, Linux x86/x64 & arm/arm64 architectures.



## Usage
```
img2sdat.py <system_img> [-o outdir] [-v version] [-p prefix]
```
- `<system_img>` = input system image
- `[-o outdir]` = output directory (current directory by default)
- `[-v version]` = transfer list version number (1 - 5.0, 2 - 5.1, 3 - 6.0, 4 - 7.0, will be asked by default, more info on xda thread)
- `[-p prefix]` = name of image (prefix.new.dat)


## 步骤
1. 将system目录转为system.img

```bash
.\make_ext4fs.exe -s -l 2097152000 -a system .\output\system.img system
```

2. 在`output`生成文件 `system.new.dat`, `system.patch.dat`, `system.transfer.list`

```bash
python ./img2sdat.py .\output\system.img -o output -v 4
```

3. 将`system.new.dat`压缩为`system.new.dat.br`
- 平衡方案
```bash
.\brotli.exe --quality=6 --output=.\output\system.new.dat.br .\output\system.new.dat
```

- 更高压缩（发布）
```bash
.\brotli.exe --quality=11 --output=.\output\system.new.dat.br .\output\system.new.dat
```


## Info
For more information about this binary, visit http://forum.xda-developers.com/android/software-hacking/how-to-conver-lollipop-dat-files-to-t2978952.
