## 当贝B3 PRO Armbian dtb

### 使用教程
1.下载Armbian的镜像（https://github.com/ophub/amlogic-s9xxx-armbian/releases）
| **版本名称**                                | **基础系统** | **发行版本**  | **描述/特点**                                                                                         |
| --------------------------------------- | -------- | --------- | ------------------------------------------------------------------------------------------------- |
| **Armbian\_trixie**                     | Debian   | Testing   | 基于 Debian **Trixie**，Debian 的测试版本，包含相对较新的软件包和特性，但可能不如稳定版可靠。                                       |
| **Armbian\_noble**                      | Ubuntu   | 22.04 LTS | 基于 **Ubuntu 22.04 LTS**，适用于长时间支持，提供稳定性和安全性，适合生产环境中的 ARM 开发板使用。                                    |
| **Armbian\_jammy**                      | Ubuntu   | 22.10     | 基于 **Ubuntu 22.10**，具有较新版本的 Ubuntu 软件包，但可能不如 LTS 稳定，适合需要更新软件和新特性的开发者。                             |
| **Armbian\_bullseye**                   | Debian   | 11        | 基于 Debian **Bullseye**，这是 Debian 的稳定版本，注重长期支持和系统稳定性，适用于大多数 ARM 开发板。                               |
| **Armbian\_bookworm**                   | Debian   | 12        | 基于 Debian **Bookworm**，Debian 的下一个稳定版（目前还在开发中），相比 Bullseye 会有更高版本的应用和新特性。                         |
| **Armbian\_HassIoSupervisor\_bookworm** | Debian   | 12        | 基于 Debian **Bookworm**，专门为 **Home Assistant Supervisor** 设计，适用于智能家居的 ARM 设备，提供 Home Assistant 支持。 |

- 选择版本后搜索s922x，下载指定版本，比如：`Armbian_25.08.0_amlogic_s922x_bullseye_6.1.147_server_2025.08.01.img.gz`
- 使用balenaEtcher刷入镜像到U盘中

2. 使用配置文件
- 下载Releases并解压，拷贝文件到U盘\dtb\amlogic目录下
- 修改U盘的根目录的文件`uEnv.txt`，将`FDT=/dtb/amlogic/xxx`改为`FDT=/dtb/amlogic/meson-g12b-dangbei-b3-pro.dtb`

### Github Actions

1. 运行工作流
- 推送标签触发
    ```bash
    git tag v1.0.0
    git push origin v1.0.0
    ```
- 初次：设置仓库权限
    - 点击仓库右上角的 Settings
    - 左侧菜单选择 Actions → General
    - 找到 Workflow permissions
    - 选择 Read and write permissions
    - 勾选 Allow GitHub Actions to create and approve pull requests

2. 下载生成文件，替换到U盘

3. 替换`uEnv.txt`文件

### 安装进emmc
1. 参照以下选项
    ```bash
    Please Input SoC Name(such as s9xxx): s922x
    Please Input DTB Name(such as meson-xxx.dtb): meson-g12b-dangbei-b3-pro.dtb
    Please Input UBOOT_OVERLOAD Name(such as u-boot-xxx.bin): u-boot-s905x2-s922.bin
    Please Input MAINLINE_UBOOT Name(such as xxx-u-boot.bin.sd.bin):u-boot-s905x2-s922.bin
    Please Input BOOTLOADER_IMG Name(such as xxx-bootloader.img):
    [ INFO ] Input Box ID: [ 0 ]
    [ INFO ] Model Name: [ GT-King-Pro,X88-King ]
    [ INFO ] FDTFILE: [ meson-g12b-dangbei-b3-pro.dtb ]
    [ INFO ] MAINLINE_UBOOT: [ u-boot-s905x2-s922.bin ]
    [ INFO ] BOOTLOADER_IMG:  [ ]
    [ INFO ] UBOOT_OVERLOAD: [ u-boot-s905x2-s922.bin ]
    [ INFO ] NEED_OVERLOAD: [ no ]
    ```
### 自编译

1. 安装环境
    ```bash
    sudo apt update
    apt install -y gcc-aarch64-linux-gnu libssl-dev
    ```
2. 编译
    ```bash
    cd u-boot
    export CROSS_COMPILE=aarch64-linux-gnu-
    make distclean
    make dangbei-b3-pro_defconfig
    make -j$(nproc)
    ```
3.生成boot文件
- 根目录下`u-boot.bin`
    | **文件名**              | **作用**                                   |
    | -------------------- | ---------------------------------------- |
    | `u-boot`             | U-Boot 可执行文件，包含 U-Boot 核心映像，用于引导嵌入式设备。   |
    | `u-boot-elf.lds`     | 链接脚本文件，定义如何将各个部分链接在一起，并指定内存布局。           |
    | `u-boot-nodtb.bin`   | 不包含设备树的 U-Boot 可执行映像。                    |
    | `u-boot.cfg`         | U-Boot 配置文件，包含环境变量和启动参数。                 |
    | `u-boot.dtb`         | 设备树二进制文件，包含硬件配置信息，通常传递给操作系统（如 Linux 内核）。 |
    | `u-boot.lds`         | 链接脚本文件，与 `u-boot-elf.lds` 类似，用于链接配置。     |
    | `u-boot.srec`        | SREC 格式的 U-Boot 映像文件，通常用于烧录到嵌入式设备。       |
    | `u-boot-dtb.bin`     | 包含设备树的 U-Boot 可执行文件，设备树与 U-Boot 一起加载。    |
    | `u-boot-elf.o`       | U-Boot 目标文件（object file），编译过程中生成的中间文件。   |
    | `u-boot.bin`         | 标准二进制 U-Boot 映像文件，用于烧录到嵌入式设备的存储器中。       |
    | `u-boot.cfg.configs` | U-Boot 构建配置文件，描述不同硬件平台和启动选项。             |
    | `u-boot.elf`         | U-Boot 的 ELF 格式可执行文件，包含调试信息，通常用于开发和调试。   |
    | `u-boot.map`         | U-Boot 链接时生成的符号映射文件，帮助开发者调试和优化代码。        |
    | `u-boot.sym`         | U-Boot 符号文件，包含所有符号（变量、函数名等）信息，用于调试和分析程序。 |

4. 打包下载文件
    ```bash
    tar --transform='s|^.*/||'  -czvf u-boot_files.tar.gz u-boot* ./arch/arm/dts/meson-g12b-dangbei-b3-pro.dtb
    ```

### 问题
1. 无法使用编译生成的boot文件，使用amlogic-boot-fip生成相应文件后，无任何报错，只是不断重启

2. 原厂固件中提取`bootloader`，生成`bl2` + `bl30` + `bl301` + `bl31`后重新编译仍然失败

### 链接

**[如何制作-u-boot-文件](https://github.com/ophub/amlogic-s9xxx-armbian/blob/main/documents/README.cn.md#1211-%E5%A6%82%E4%BD%95%E5%88%B6%E4%BD%9C-u-boot-%E6%96%87%E4%BB%B6)**

**[u-boot](https://github.com/unifreq/u-boot)**

**[amlogic-boot-fip](https://github.com/unifreq/amlogic-boot-fip)**
