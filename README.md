[![License](https://img.shields.io/github/license/erkexzcx/lietuviska-tv)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/erkexzcx/lietuviska-tv)](https://goreportcard.com/report/github.com/erkexzcx/lietuviska-tv)
[![Github All Releases](https://img.shields.io/github/downloads/erkexzcx/lietuviska-tv/total.svg)](https://github.com/erkexzcx/lietuviska-tv/releases)

# Nemokama lietuviška televizija internetu

Ši programa veikia kaip tarpinis serveris tarp IPTV kliento (pvz VLC, Kodi) ir viešai prieinamų ir nemokamų lietuviškų IPTV stream'ų (pvz LNK, TV3).

**Telegram pokalbių kambarys**: https://t.me/lietuviskaIPTV

Palaikomų kanalų sąrašas:
* BBC World News
* BTV
* CTC Baltija
* Esports
* Euronews
* FOX LT
* Fight sport
* INFO TV HD
* Kidzone Plus
* LNK
* LNK HD
* LRT
* LRT HD
* LRT Plius
* LRT Plius HD
* Lietuvos rytas
* Lietuvos rytas HD
* MTV Hits
* My Hits
* NBA
* NTV Mir Baltic LT
* NatioNAL Geographic
* Nick Junior
* Nickelodeon LT
* Ohota y Ribalka
* PBK LT
* Ren TV Baltic LT
* Setanta Sports
* Sony
* Sony Turbo
* Sport 1
* TV1
* TV1000
* TV1000 Action
* TV1000 Ruskoje kino
* TV3
* TV3 Film
* TV3 HD
* TV3 Sport
* TV3 Sport 2
* TV6
* TV8
* Viasat Explore
* Viasat History
* Viasat Nature

# Naudojimas

## Windows

1. Atsisiunčiate naujausią binary (.exe failą) iš [releases](https://github.com/erkexzcx/lietuviska-tv/releases/latest) (pavyzdžiui `lietuviskatv_windows_i386.exe`) ir išsisaugojate kompiuteryje.

2. Kompiuteryje atsidarykite failų naršyklę (file explorer) ir nueikite į tą aplanką (folder), kuriame yra atsisiųstas .exe failas. Tuomet laikykite nuspaudę SHIFT mygtuką ir neatleidę spauskite dešinį pelės klavišą ant balto fono (tame pačiame file explorer/failų naršyklėje). Iššokusiame lange atsiras pasirinkimas *Open PowerShell Window here*. Jį paspauskite ir atsidarys PowerShell aplikacija (terminalas).

3. PowerShell aplikacijoje parašykite (nukopijuokite) komandą būtent taip ir spauskite enter:
```
.\lietuviskatv_windows_i386.exe
```
Ir paleista programa ekrane parašys `Started!`. Jeigu pamatėte šį teksta - programą sėkmingai paleidote. Palikite ją veikti (neuždarykite PowerShell programos).

4. Atidarykite VLC programą, joje pasirinkite `Medua` --> `Open Network Stream...` ir į matomą laukelį įveskite adresą `http://127.0.0.1:8989/iptv`. Spauskite Enter klavišą ir IPTV pradės rodyti. VLC programoje bus matomas visų kanalų sąrašas ir galėsite pasikeisti matomą kanalą.

## Linux, MacOS ir FreeBSD

Atsisiunčiate naujausią binary iš [releases](https://github.com/erkexzcx/lietuviska-tv/releases/latest) ir paleidžiate terminale.

```
# Linux, MacOS ir FreeBSD:
chmod +x lietuviskatv_<platform>_<architecture>
./lietuviskatv_<platform>_<architecture>
```

Ir tuomet IPTV playlist pasiekiamas per šią nuorodą: `http://<address>:8989/iptv` (jei ant to paties kompiuterio: `http://127.0.0.1:8989/iptv`). Šią nuorodą naudokit ant VLC arba Kodi su *Simple IPTV addon*.

P.S. Linux SystemD service sukursiu ateityje. Šiuo metu patariu naudoti `tmux` ir palikti veikti background'e.





# FAQ

## Kaip ši programa veikia?

Programa atlieka kelias funkcijas:
1. Pateikia `M3U` IPTV kanalų playlist.
2. Sugeneruoja keletos IPTV kanalų URL (kai kurie yra kintantys).
3. Veikia kaip tarpininkas (savotiškas proxy serveris) tarp IPTV kliento ir prieinamų IPTV kanalų adresų. Visas IPTV srautas keliauja per šią programą.

Šiuo metu esu radęs apie 30+ neapribotų IPTV kanalų, kurių `M3U8` adresas nekinta. Visi tie adresai yra įvesti programos kode (AKA hardcoded) ir jie tiesiog "yra". Likę keli kanalai yra nuolat kintantys, ir juos reikia išgauti programos pagalba.

Galbūt pastebėjote, kad užkrovus `<address>:8989/iptv` visų kanalų nuorodos yra adresuotos į tokį patį adresą, kuriuo yra pasiekiama programa (paminėkim žodį *proxy*). Tai yra dėl keletos priežąsčių:
1. Programa kas 2 valandas atnaujina dinaminių kanalų nuorodas, tačiau Kodi to nežino, kad IPTV kanalo nuoroda atsinaujino (ir niekada nesužino - tiesiog taip veikia). Kad priversti Kodi sužinoti naują nuorodą, reikia perkrauti arba pareloadinti *Simple IPTV addon*, antraip kanalas po kiek laiko nebus rodomas. Su šio proxinimo pagalba (URL perrašymu), iš Kodi perspektyvos, IPTV kanalo nuoroda yra visada vienoda, o pati programa ją fone nuolat atnaujina.
2. Kodi nenorėjo normaliai elgtis su pilna TV3 nuoroda - per Kodi net nerodydavo TV3 kanalo, o telefone rodydavo. Su šio proxinimo pagalba (URL perrašymu), šios problemos nebeliko.
3. Sugeneruota dinaminio kanalo nuoroda būna su tam tikru sesijos ID, kurį sugeneruoja pats IPTV kanalo serveris. Kiekvieną kart, kai IPTV klientas kreipiasi į serverį, yra tikrinama ir prie sesijos ID pririštas public IP (jeigu jis kitoks - kanalas nerodomas ir gaunamas HTTP 403 error). Su šio proxinimo pagalba (URL perrašymu), IPTV kanalo serverio pusės matomas IP adresas visada bus vienodas, todėl dabar šią programą galima talpinti išoriniame serveryje (pvz Google cloud).

## Kai kurių kanalų nerodo

Kai kurių kanalų tiesiog nerodo. Jeigu pastesbėsi, kad apskritai neberodo daugumos arba nustojo rodyti konkretų kanalą - kelk naują [issue](https://github.com/erkexzcx/lietuviska-tv/issues).

## Ant VLC atsilieka garsas

Gali būt, kad ant VLC atsilieka **TV3** kanalo garsas. Jei taip nutiko - ant VLC reik spaust dešinį pelės klavišą (ant rodomo video) --> `tools` --> `Track synchronization` --> `Audio track synchronization` ir pakeisti į `-1`.

## Nesuprantu platformų ir/ar architektūrų

```
lietuviskatv_darwin_x86_64 --> MacOS, 64bit
lietuviskatv_freebsd_x86_64 --> FreeBSD platformai, 64bit (pvz pfsense sistemai)
lietuviskatv_linux_aarch64 --> Linux, aarch64 (armv8) (rpi3 su 64bit OS, rpi4 su 64bit OS)
lietuviskatv_linux_arm --> Linux, arm (armv5 ir armv6) (rpi0, rpi1)
lietuviskatv_linux_armhf --> Linux, armhf (armv7) (rpi2, rpi3 su 32bit OS, rpi4 su 32bit OS)
lietuviskatv_linux_i386 --> Linux, 32bit
lietuviskatv_linux_x86_64 --> Linux, 64bit
lietuviskatv_windows_x86_64.exe --> Windows, 64bit
lietuviskatv_windows_i386.exe --> Windows, 32bit
```

## Žinau kanalą, kurį galima žiūrėti internetu, tačiau jo nėra tavo programoje

Pakelk naują [issue](https://github.com/erkexzcx/lietuviska-tv/issues) šiam projektui surašydamas visas detales kur kas ir kaip. Pridėsiu į projektą.

## Trūksta norimos platformos ir/ar architektūros

Jei norite pasileisti ant platformos ar architektūros, kurios nėra pateiktuose binaries - reikia pačiam sukompiliuoti binary. Ant Linux įsirašykite Go ([štai taip](https://www.digitalocean.com/community/tutorials/how-to-install-go-on-debian-9), nes official repos esanti versija dažniausiai būna per sena), atsisiųskite šį projektą ir tada (pavyzdžiui OpenWRT naudojamai `MIPS` `softfloat`):
```
# Binary kompiliavimo komanda:
env GOOS=linux GOARCH=mips GOMIPS=softfloat go build -ldflags="-s -w" -o "lietuviskatv_linux_mips_softfloat" src/*.go

# Jei nori daugiau nei per pusę sumažinti binary dydį, tačiau binary gali
# nepasileisti arba nesusispausti
upx --best "lietuviskatv_linux_mips_softfloat" 
```
Daugiau informacijos apie galimas platformas ir architektūras: https://golang.org/doc/install/source#environment
