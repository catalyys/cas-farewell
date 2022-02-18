# cas - farewell

Celeste Auto-Splitter for Linux.

## Installing

```
git clone https://github.com/DevilFreak/cas-farewell.git
cd cas-farewell/
go build
```

## running cas - farewell

You can just execute the compiled binary and it will start with the default configurations.

Here are all the things you can configure without changing the code.

| Argument | Usage                                 | example |
| -------- | ------------------------------------- | ------- |
| run      | start the application (same as executing with nothing)| ./cas run       |
| help     | shows the help                                      | ./cas help       |
| show     | shows you personal best or best splits              | ./cas show best / ./cas show splits    |
| -i       | gives you more information                          | ./cas -i show best    |
| -s       | gives you more information about your splits        | ./cas -si run    |
| -save, -savefile| changes the savefile slot [0, 1, 2]       | ./cas -save 0    |

## Route configuration

Modify `anyPercent` variable in `types.go`.

## Showcase

### Split

![](example/split.gif)

### New Run

![](example/autodelete.gif)


#### Credits

This Repository was built on top of the original cas from ~bfiedler.
The original Repo can be found at https://sr.ht/~bfiedler/cas/
