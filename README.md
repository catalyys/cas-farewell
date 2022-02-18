# cas - farewell

Celeste Auto-Splitter for Linux.

## Compiling

`go build`

## running cas - farewell

You can just execute the compiled binary and it will start with the default configurations.

Here are all the things you can configure without changing the code.

| Argument | Usage                                 | example |
| -------- | ------------------------------------- | ------- |
| run | start the application (same as executing with nothing)         | ./cas run       |
| help     | shows the help                   | ./cas help       |
| show    | shows you personal best or best splits              | ./cas show best | ./cas show splits    |
| -i    | gives you more information | ./cas -i show best    |
| -s     | gives you more information about your splits                | ./cas -si run    |

## Route configuration

Modify `anyPercent` variable in `types.go`.

## Showcase

### Split

![](example/split.gif)

### New Run

![](example/autodelete.gif)
