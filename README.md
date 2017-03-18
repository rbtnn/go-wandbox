
# go-wandbox

[wandbox](http://melpon.org/wandbox) CLI

## usage

    $make
    go build -o wandbox src/main.go
    $cat ~/a.vim
    echo "hi"
    $./wandbox -compiler vim-head -source ~/a.vim
    hi
    $./wandbox -compiler vim-head -code "echo 'hi'"
    hi
    $./wandbox -help
    Usage of ./wandbox:
      -code string
            code
      -compiler string
            compiler
      -list
            compiler list
      -source string
            source file

    $./wandbox -list
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
    ...
