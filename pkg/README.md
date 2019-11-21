# Go-mdbm

The following is pre-build packages for easy install to your machine.

## List of packages

```bash
|-- docker
|   |-- centos
|   |   |-- 6
|   |   |   `-- Dockerfile
|   |   `-- 7
|   |       `-- Dockerfile
|   `-- ubuntu
|       |-- 12.04
|       |   `-- Dockerfile
|       |-- 14.04
|       |   `-- Dockerfile
|       |-- 16.04
|       |   `-- Dockerfile
|       `-- 18.04
|           `-- Dockerfile
|-- osx
|   `-- mdbm.rb
|-- rhel
|   |-- el6
|   |   |-- mdbm-4.12.4.0-1.el6.x86_64.rpm
|   |   |-- mdbm-debuginfo-4.12.4.0-1.el6.x86_64.rpm
|   |   |-- mdbm-devel-4.12.4.0-1.el6.x86_64.rpm
|   |   `-- mdbm-perl-4.12.4.0-1.el6.x86_64.rpm
|   |-- el7
|   |   |-- mdbm-4.12.4.0-1.el7.src.rpm
|   |   |-- mdbm-4.12.4.0-1.el7.x86_64.rpm
|   |   |-- mdbm-debuginfo-4.12.4.0-1.el7.x86_64.rpm
|   |   |-- mdbm-devel-4.12.4.0-1.el7.x86_64.rpm
|   |   `-- mdbm-perl-4.12.4.0-1.el7.x86_64.rpm
|   `-- el8
|       |-- cppunit-1.12.1-3.1.el8.src.rpm
|       |-- cppunit-1.12.1-3.1.el8.x86_64.rpm
|       |-- cppunit-debuginfo-1.12.1-3.1.el8.x86_64.rpm
|       |-- cppunit-debugsource-1.12.1-3.1.el8.x86_64.rpm
|       |-- cppunit-devel-1.12.1-3.1.el8.x86_64.rpm
|       |-- cppunit-doc-1.12.1-3.1.el8.x86_64.rpm
|       |-- doxygen-1.8.16-2.el8.src.rpm
|       |-- doxygen-1.8.16-2.el8.x86_64.rpm
|       |-- doxygen-debuginfo-1.8.16-2.el8.x86_64.rpm
|       |-- doxygen-debugsource-1.8.16-2.el8.x86_64.rpm
|       |-- doxygen-doxywizard-1.8.16-2.el8.x86_64.rpm
|       |-- doxygen-doxywizard-debuginfo-1.8.16-2.el8.x86_64.rpm
|       |-- doxygen-latex-1.8.16-2.el8.x86_64.rpm
|       |-- mdbm-4.12.4.0-1.el8.src.rpm
|       |-- mdbm-4.12.4.0-1.el8.x86_64.rpm
|       |-- mdbm-debuginfo-4.12.4.0-1.el8.x86_64.rpm
|       |-- mdbm-debugsource-4.12.4.0-1.el8.x86_64.rpm
|       |-- mdbm-devel-4.12.4.0-1.el8.x86_64.rpm
|       |-- mdbm-perl-4.12.4.0-1.el8.x86_64.rpm
|       `-- mdbm-perl-debuginfo-4.12.4.0-1.el8.x86_64.rpm
`-- ubuntu
    |-- mdbm-4.13.0-Bionic_Beaver.deb
    |-- mdbm-4.13.0-Bionic_Beaver.deb.md5
    |-- mdbm-4.13.0-Disco_Dingo.deb
    |-- mdbm-4.13.0-Disco_Dingo.deb.md5
    |-- mdbm-4.13.0-Precise_Pangolin.deb
    |-- mdbm-4.13.0-Precise_Pangolin.deb.md5
    |-- mdbm-4.13.0-Trusty_Tahr.deb
    |-- mdbm-4.13.0-Trusty_Tahr.deb.md5
    |-- mdbm-4.13.0-Xenial_Xerus.deb
    `-- mdbm-4.13.0-Xenial_Xerus.deb.md5
```

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
|RHEL6 (CentOS6.x)|64bit|mdbm-4.12.3.0-1.el6.x86_64.rpm|
|||mdbm-debuginfo-4.12.3.0-1.el6.x86_64.rpm|
|||mdbm-devel-4.12.3.0-1.el6.x86_64.rpm|
|||mdbm-perl-4.12.3.0-1.el6.x86_64.rpm|
|RHEL7 (CentOS7.x)|64bit|mdbm-4.12.3.0-1.el6.x86_64.rpm|
|||mdbm-debuginfo-4.12.3.0-1.el6.x86_64.rpm|
|||mdbm-devel-4.12.3.0-1.el6.x86_64.rpm|
|||mdbm-perl-4.12.3.0-1.el6.x86_64.rpm|
|RHEL8 (CentOS8.x)|64bit|mdbm-4.12.4.0-1.el8.x86_64.rpm|
|||mdbm-debuginfo-4.12.4.0-1.el8.x86_64.rpm|
|||mdbm-debugsource-4.12.4.0-1.el8.x86_64.rpm|
|||mdbm-devel-4.12.4.0-1.el8.x86_64.rpm|
|||mdbm-perl-4.12.4.0-1.el8.x86_64.rpm|
|||mdbm-perl-debuginfo-4.12.4.0-1.el8.x86_64.rpm|

### Installation (RHEL6~8, CentOS 6~8)

```shell
git clone https://github.com/torden/go-mdbm
rpm -Uvh pkg/rhel/rhel/elX/mdbm-4.12.3.0-1.elX.x86_64.rpm
rpm -Uvh pkg/rhel/rhel/elX/mdbm-devel-4.12.3.0-1.elX.x86_64.rpm
rpm -Uvh pkg/rhel/rhel/elX/mdbm-debuginfo-4.12.3.0-1.elX.x86_64.rpm
```


## OSX

### Using Homebrew  

```shell
/usr/bin/ruby -e "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/master/install)"
git clone https://github.com/torden/go-mdbm

brew install cd go-mdbm/pkg/osx/mdbm.rb

# OR (as soon)

brew install mdbm 
```

## BSD

as soon

