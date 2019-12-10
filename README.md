[![License](https://img.shields.io/github/license/erkexzcx/lietuviska-tv)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/erkexzcx/lietuviska-tv)](https://goreportcard.com/report/github.com/erkexzcx/lietuviska-tv)
[![Github All Releases](https://img.shields.io/github/downloads/erkexzcx/lietuviska-tv/total.svg)](https://github.com/erkexzcx/lietuviska-tv/releases)

# Nemokama lietuviška televizija internetu

Ši programa veikia kaip tarpinis serveris tarp IPTV kliento (pvz VLC, Kodi) ir viešai prieinamų ir nemokamų lietuviškų IPTV stream'ų (pvz LNK, TV3). Yra galimi tokie kanalai:

* LNK HD (Perskaitykite [#Troubleshooting](#Troubleshooting))
* TV3 HD
* INFO TV HD
* Lietuvos rytas HD
* LRT HD
* LRT Plius HD

# Naudojimas

Atsisiunčiate naujausią binary iš [releases](https://github.com/erkexzcx/lietuviska-tv/releases). Tuomet programa paleidžiame terminale (Windows naudoja Powershell):
```
./lietuviskatv_<platform>_<architecture>
```
Ir tuomet IPTV playlist pasiekiamas per šią nuorodą: `http://<address>:8989/iptv`

P.S. Linux SystemD service sukursiu ateityje. Šiuo metu patariu naudoti `tmux` ir palikti veikti background'e.

# Troubleshooting

## LNK

Kai kuriuo nors metu per LNK GO yra rodoma gyvos transliacijos laida (pvz vakaro žinios), ši programa nuparsina rodomos transliacijos nuorodą, kuri yra ne kas kita kaip paties LNK stream'o nuoroda. Kai transliacija baigiasi - stream'o nuorodos nebelieka, tačiau programa iš atminties išsitraukia paskutinę galimą nuorodą, kuri veikia maždaug iki kitos dienos pietų.

Todėl jeigu šiuo metu nerodo LNK - palauk vakaro su šia įjungta programa. Ji nuparsins LNK nuorodą ir iš karto pradės rodyti.

## Atsilieka garsas

Gali būt, kad ant VLC atsilieka TV3 kanalo garsas. Jei taip nutiko - ant VLC reik spaust dešinį pelės klavišą (ant rodomo video) --> `tools` --> `Track synchronization` --> `Audio track synchronization` ir pakeisti į `-1`.

## Nesuprantu platformų ir/ar architektūrų

```
lietuviskatv_darwin_x86_64 --> MacOS, 64bit
lietuviskatv_freebsd_x86_64 --> FreeBSD platformai, 64bit (pvz pfsense sistemai)
lietuviskatv_linux_aarch64 --> Linux, aarch46 (armv8) (rpi3 su 64bit OS, rpi4 su 64bit OS)
lietuviskatv_linux_arm --> Linux, arm (armv5 ir armv6) (rpi0, rpi1)
lietuviskatv_linux_armhf --> Linux, armhf (armv7) (rpi2, rpi3 su 32bit OS, rpi4 su 32bit OS)
lietuviskatv_linux_i386 --> Linux, 32bit
lietuviskatv_linux_x86_64 --> Linux, 64bit
lietuviskatv_windows_x86_64.exe --> Windows, 64bit
lietuviskatv_windows_i386.exe --> Windows, 32bit
```

## Trūksta norimos platformos ir/ar architektūros

Jei norite pasileisti ant platformos ir/ar architektūros, kurios nėra pateiktuose binaries - reikia tai atlikti pačiam. Pasiruoškite Linux kompiuteryje Golang aplinką, atsisiųskite šį projektą ir tada (pavyzdžiui OpenWRT naudojamai `MIPS` `softfloat`):
```
env GOOS=linux GOARCH=mips GOMIPS=softfloat go build -ldflags="-s -w" -o "lietuviskatv_linux_mips_softfloat" src/*.go
upx --best "lietuviskatv_linux_mips_softfloat" # Daugiau nei per pusę sumažina sukompiliuoto binary dydį
```
Daugiau informacijos apie galimas platformas ir architektūras: https://golang.org/doc/install/source#environment
