# Go-mdbm

The following is pre-build packages for easy install to your machine.

## Ubuntu

### List of Pre-build Packages

|OS/Release Ver.|Arch.|Pkg File|dpkg::Depends|dpkg::Suggests|
|---|---|---|---|---|
|Ubuntu 19.xx|64bit|mdbm-4.13.0-Disco_Dingo.deb|zlib1g, libssl1.1, libreadline8, libtinfo6, libstdc++6, libc6|perl-modules|
|Ubuntu 18.xx|64bit|mdbm-4.13.0-Bionic_Beaver.deb|zlib1g, libssl1.1, libreadline7, libtinfo5, libstdc++6, libc6|per-modules|
|Ubuntu 16.xx|64bit|mdbm-4.13.0-Xenial_Xerus.deb|zlib1g, libssl1.0.0, libreadline6, libtinfo5, libstdc++6, libc6|perl-modules|
|Ubuntu 14.xx|64bit|mdbm-4.13.0-Trusty_Tahr.deb|zlib1g, libssl1.0.0, libreadline6, libtinfo5, libstdc++6, libc6|perl-modules|
|Ubuntu 12.xx|64bit|mdbm-4.13.0-Precise_Pangolin.deb|zlib1g, libssl1.0.0, libreadline6, libtinfo5, libstdc++6, libc6|perl-modules|


### Installation

```shell
git clone https://github.com/torden/go-mdbm
dpkg -i pkg/ubuntu/mdbm-XXX.deb
```

## RedHat (RHEL)

### List of Pre-build Packages

|OS/Release Ver.|Arch.|Pkg File|
|---|---|---|
|RHEL6 ~ 7 (CentOS6.x ~7.x)|64bit|mdbm-4.12.3.0-1.el6.src.rpm|
|||mdbm-4.12.3.0-1.el6.x86_64.rpm|
|||mdbm-debuginfo-4.12.3.0-1.el6.x86_64.rpm|
|||mdbm-devel-4.12.3.0-1.el6.x86_64.rpm|
|||mdbm-perl-4.12.3.0-1.el6.x86_64.rpm|

### Installation (RHEL6, CentOS 6)

```shell
git clone https://github.com/torden/go-mdbm
rpm -Uvh pkg/rhel/rhel/el6/mdbm-4.12.3.0-1.el6.x86_64.rpm
rpm -Uvh pkg/rhel/rhel/el6/mdbm-devel-4.12.3.0-1.el6.x86_64.rpm
rpm -Uvh pkg/rhel/rhel/el6/mdbm-debuginfo-4.12.3.0-1.el6.x86_64.rpm
```


## OSX

as soon

## BSD

as soon


