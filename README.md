
# go-wandbox

wandbox CLI

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

