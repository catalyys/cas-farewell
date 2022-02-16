# cas

Celeste Auto-Splitter (well, kind-of) for Linux. Currently it is only possible
to split on chapters.

It automatically uses the bottom save file for timing runs. You can adjust this
behaviour by editing the `saveFile` variable in `main.go`.

## Compiling

`go build`

## Route configuration

Modify `anyPercent` variable in `types.go`.

## Import PB

Before you can split you should import your PB. Just call the program with the
path to your preferred save file (found in `$XDG_DATA_HOME/Celeste/Saves/`) and
it will initialize `pb.json` which contains your PB and `bule.json` which will
contain your best single splits.

## Showcase

### Split

![](example/split.gif)

### New Run

![](example/autodelete.gif)
