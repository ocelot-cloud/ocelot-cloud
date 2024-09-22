### Build From Scratch

You can build the image directly from source code using these commands:

```bash
cd ocelot-cloud
projectDir=$(pwd)
cd $projectDir/scripts
bash install.sh # Only tested on Debian 12 ("Bookworm")
cd $projectDir/src/ci-runner
go build
# Either you just build the image:
./ci-runner build
# or you build and directly deploy it:
./ci-runner deploy
```