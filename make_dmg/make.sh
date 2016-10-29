diskutil unmount egaas
GOARCH=amd64  CGO_ENABLED=1  go build -o make_dmg/egaas.app/Contents/MacOS/egaasbin
cd make_dmg
./make_dmg.sh -b background.png -i logo-big.icns -s "480:540" -c 240:400:240:200 -n egaas_osx64 "egaas.app"
cd ../