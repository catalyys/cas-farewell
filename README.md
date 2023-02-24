# cas - farewell

Celeste Auto-Splitter for Linux.



## Pre-Version 1

All the times changed the location to `~/.config/casf/casf.json`.
To import your old pb and bule times, backup your pb and bule file before updating and import both of them.
```
./casf import --pb --file ./pb.json --run any # here you can name your run if you have more than just any%
./casf import --bule --file ./bule.json
```

All run are now saved in `~/.config/casf/casf.json` and support more than one run/route and even custom routes.



## Installing

Download the release or compile yourself.



### Compiling

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
| route    | configure route for runs                            | ./casf route show    |
| import   | import old pb or bule files                         | ./casf import  -bule -file ./bule.json |
| -i       | gives you more information                          | ./casf show best -i |
| -s       | gives you more information about your splits        | ./casf run -is    |
| -n       | changes chapter names to numbers                    | ./casf -n    |
| -z       | adds the side letter to chapter                     | ./casf -z   |
| -save, -saveslot| changes the saveslot [1, 2, 3] _3 is default_      | ./casf -save 1    |
| -route, -r| changes the route of the run      | ./casf -route anyB    |



### show arguments

| Argument | Usage                                 | example |
| -------- | ------------------------------------- | ------- |
| best     | shows personal best                   | ./casf show best --run any      |
| bule     | shows best splits                     | ./casf show bule      |
| routes   | shows all routes                      | ./casf show routes    |



### route arguments

| Argument | Usage                                 | example |
| -------- | ------------------------------------- | ------- |
| create   | creates a custom run                  | ./casf route create -name mycustomrun -route "1:a,2:b"     |
| show     | shows all run (same as `./casf show routes`) | ./casf route show      |
| remove   | remove a custom run                   | ./casf route delete -name mycustomrun    |

In the route create, the chapter can have following formats: `2:a`, `2:0`, `2A`, `20`

A custom route is created in the config file like this:

```json
:
  "customruns": {
    "mycustomrun": [
      "1:0",
      "2:1"
    ]
  }
:
```

You can also edit the run from here or delete it.



### Overlay

To run cas - farewell as an overlay you can simply use a terminal with 100% transparent background and mark as always on top. Then just resize it and move it to where you want your overlay to be.

I have done it via bash:

```bash
xfce4-terminal --geometry=50x13+0+130 --hide-borders --working-directory="$HOME/git/cas-farewell/" -e "./casf -is" -H
```

and then mark as always on top (ALT + Spacebar to open menu):
![](example/terminal.png)



### Settings

You can set some settings in the config file. Currently the default config file looks like this.

```json
:
  "settings": {
    "celeste_savefolder": "/home/olli/.local/share/Celeste/Saves/",
    "default_saveslot": "3",
    "default_run": "any",
    "flag_i": "false",
    "flag_s": "false",
    "flag_n": "false",
    "flag_z": "false"
  }
:
```



## Showcase

### Split

![](example/split.gif)

### New Run

![](example/autodelete.gif)


#### Credits

This Repository was built on top of the original cas from ~bfiedler.
The original Repo can be found at https://sr.ht/~bfiedler/cas/
