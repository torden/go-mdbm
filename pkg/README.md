# Go-mdbm

The following is pre-build packages for easy install to your machine.

## Ubuntu

### List of Pre-build Packages

|OS/Release Ver.|Arch.|Pkg File|
|---|---|---|
|Ubuntu 17.10|64bit|mdbm-Artful_Aardvark.deb|
|Ubuntu 12.04.5 LTS|64bit|mdbm-Precise_Pangolin.deb|
|Ubuntu 14.04.5 LTS|64bit|mdbm-Trusty_Tahr.deb|

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


