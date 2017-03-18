
# go-wandbox

[wandbox](http://melpon.org/wandbox) CLI

## usage

    $make
    go build -o wandbox src/main.go
    $cat ~/a.vim
    echo "hi"
    $./wandbox -c vim-head -f ~/a.vim
    hi
    $./wandbox --help
    Usage of ./wandbox:
      -c string
            compiler
      -f string
            source file
      -list
            compiler list
    $./wandbox --list
    gcc-head-c
    gcc-6.3.0-c
    gcc-6.2.0-c
    gcc-6.1.0-c
    gcc-5.4.0-c
    gcc-5.3.0-c
    gcc-5.2.0-c
    gcc-5.1.0-c
    gcc-4.9.3-c
    gcc-4.9.2-c
    gcc-4.9.1-c
    gcc-4.9.0-c
    gcc-4.8.5-c
    gcc-4.8.4-c
    gcc-4.8.3-c
