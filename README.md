# cas - farewell

Celeste Auto-Splitter for Linux.

## Installing

```
git clone https://github.com/catalyys/cas-farewell.git
cd cas-farewell/
go build
```

## running cas - farewell

The default folder for Celeste is `~/.local/share/Celeste/Saves/` if you got it from Steam.
When you have the savefiles in a different folder you can change that in `main.go`.

You can then just execute the compiled binary and it will start with the default configurations.

Here are all the things you can configure without changing the code.

| Argument | Usage                                 | example |
| -------- | ------------------------------------- | ------- |
| run      | start the application (same as executing with nothing)| ./casf run       |
| help     | shows the help                                      | ./casf help       |
| show     | shows you personal best or best splits              | ./casf show splits    |
| -i       | gives you more information                          | ./casf show best -i |
| -s       | gives you more information about your splits        | ./casf run -is    |
| -n       | changes chapter names to numbers                    | ./casf -n    |
| -z       | adds the side letter to chapter                     | ./casf -z   |
| -save, -savefile| changes the savefile slot [0, 1, 2] _2 is default_      | ./casf -save 0    |
| -route, -r| changes the route of the run      | ./casf -route any%B    |


#### show arguments


| Argument | Usage                                 | example |
| -------- | ------------------------------------- | ------- |
| best     |  shows personal best                  | ./casf show best       |
| splits   | shows best splits                     | ./casf show splits       |
| routes   | shows all routes                      | ./casf show routes    |




### Overlay

To run cas - farewell as an overlay you can simply use a terminal with 100% transparent background and mark as always on top. Then just resize it and move it to where you want your overlay to be.

I have done it via bash:
```
xfce4-terminal --geometry=50x13+0+130 --hide-borders --working-directory="$HOME/git/cas-farewell/" -e "./casf -is" -H
```

and then mark as always on top (ALT + Spacebar to open menu):
![](example/terminal.png)

## Route configuration

To show all pre configured runs you can run `./casf show routes` to display all the routes.
To make a custum Run you need to modify the `anyPercent` variable in `types.go`.
An easier way is coming with a later version.

## Showcase

### Split

![](example/split.gif)

### New Run

![](example/autodelete.gif)


#### Credits

This Repository was built on top of the original cas from ~bfiedler.
The original Repo can be found at https://sr.ht/~bfiedler/cas/
