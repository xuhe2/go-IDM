# build

```shell
make
```

- the binary file will be in `bin` directory

# run

use `-h` arg to get help docs

```shell
./bin/go-IDM -h
```

- `-t` to use threads
- `-n` to set file name
- `-p` to set path to save file
- `-f` to skip check file transfer method
- `--proxy` to set proxy, eg: `--proxy http://127.0.0.1:1080`
- `--memory` to save tmp file part in memory, default is false

## example

```shell
./bin/go-IDM https://d1.music.126.net/dmusic/NeteaseCloudMusic_Music_official_2.10.13.202675_32.exe
```
