# Nemokama lietuviška televizija

[![Github All Releases](https://img.shields.io/github/downloads/erkexzcx/lietuviska-tv/total.svg)](https://github.com/erkexzcx/lietuviska-tv/releases)

Ši programa veikia kaip IPTV tarpinis serveris, kuris sugeneruoja šių kanalų `m3u` (_IPTV playlist_) compatible playlist'ą, kurį galima naudoti Kodi/VLC programose ir taip nemokamai žiūrėti lietuvišką TV internetu:

* LNK HD
* TV3 HD
* INFO TV HD
* Lietuvos rytas HD
* LRT HD
* LRT Plius HD

Viskas imama iš viešai prieinamų stream'ų internetu. Iš pačių LNK ir TV3 rankų, su žiupsneliu hackų :)

# Naudojimas

1. Perskaitote [#Troubleshooting](#Troubleshooting) ir įsitikinate, kad viskas jus tenkina.
2. Atsisiunčiate naujausią binary iš [releases](https://github.com/erkexzcx/lietuviska-tv/releases).
3. Paleidžiate atsisiųstą executable. Kad nereiktų jokių SystemD services rašyt, aš tiesiog naudoju `tmux` ir palieku veikti background'e:
```
./lietuviskatv_linux_armhf &
```

4. IPTV playlistas bus pasiekiamas per `http://<ipaddress>:8989/iptv`

**Newbies** - Jei norit pasileisti per/ant Windows, atsidarot powershell, pasileidžiat per jį atsisiųstą .exe failą (google kaip) ir tada naudojate `127.0.0.1:8989/iptv` ant to pačio device. Na o jei norit žiūrėt šią IPTV ant kito prietaiso (pvz Raspberry Pi su OSMC/LibreELEC), tuomet naudokit Windowsų kompiuterio IP adresą (google kaip susirast) (pvz `http://192.168.1.56:8989/iptv`). P.S. Binaries yra pateikti ir raspberiams, todėl nebūtina naudot Windwosų. *Juk Windows - sucks!*

# Troubleshooting

## TV3 atsilieka garsas

Nežinau kas konkrečiai dėl to kaltas, nes šitas issue galioja tik VLC playeriui (kiek bandžiau). Kodi neturi šitos problemos. Ant VLC reik spaust dešinį pelės klavišą ant rodomo video --> `tools` --> `Track synchronization` --> `Audio track synchronization` ir pakeisk į `-1`.

## Nerodo LNK

LNK rodo tokiu principu - kai kuriuo nors metu per LNK GO rodo kokią nors tiesioginę transliaciją, scriptas tuo metu gali nuparsinti `M3U8` nuorodą, per kurią galima žiūrėti LNK TV. Kai transliacija baigiasi - nuorodos nebelieka, tačiau scriptas iš atminties išsitraukia paskutinę galimą nuorodą, kuri dažniausiai veikia iki kitos dienos pietų.

## Reikia pagalbos dėl platformų pavadinimų. Pvz noriu įsirašyti ant RPI2

```
lietuviskatv_freebsd_x86_64 --> FreeBSD platformai, 64bit (pvz pfsense sistemai)
lietuviskatv_linux_aarch64 --> Linux, aarch46 (armv8) (rpi3 su 64bit OS, rpi4 su 64bit OS)
lietuviskatv_linux_arm --> Linux, arm (armv5 ir armv6) (rpi0, rpi1)
lietuviskatv_linux_armhf --> Linux, armhf (armv7) (rpi2, rpi3 su 32bit OS, rpi4 su 32bit OS)
lietuviskatv_linux_i386 --> Linux, 32bit
lietuviskatv_linux_x86_64 --> Linux, 64bit
lietuviskatv_windows_x86_64.exe --> Windows, 64bit
lietuviskatv_windows_i386.exe --> Windows, 32bit
```

## Trūksta norimos platformos

Jei nori pasileisti ant platformos, kurios nėra pateiktuose binaries (pvz OpenWRT routeris), teks susikompiliuoti pačiam. Pasiruošiant Linuxe golang'ą, atsisiunčiat projektą ir tada (pavyzdžiui MIPS softfloat platformai - kai kurie OpenWRT routeriai naudoja):
```
env GOOS=linux GOARCH=mips GOMIPS=softfloat go build -ldflags="-s -w" -o "dist/lietuviskatv_linux_mips_softfloat.exe" src/*.go
```
Daugiau info apie galimas architektūras ir galimus buildinimo parametrus https://golang.org/doc/install/source#environment
