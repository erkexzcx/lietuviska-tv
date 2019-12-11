[![License](https://img.shields.io/github/license/erkexzcx/lietuviska-tv)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/erkexzcx/lietuviska-tv)](https://goreportcard.com/report/github.com/erkexzcx/lietuviska-tv)
[![Github All Releases](https://img.shields.io/github/downloads/erkexzcx/lietuviska-tv/total.svg)](https://github.com/erkexzcx/lietuviska-tv/releases)

# Nemokama lietuviška televizija internetu

Ši programa veikia kaip tarpinis serveris tarp IPTV kliento (pvz VLC, Kodi) ir viešai prieinamų ir nemokamų lietuviškų IPTV stream'ų (pvz LNK, TV3).

Palaikomų kanalų sąrašas:
* BBC World News
* BTV
* CTC Baltija
* Esports
* Euronews
* FOX LT
* Fight sport
* INFO TV (D)
* Kidzone Plus
* LNK
* LNK HD (D)
* LRT
* LRT HD (D)
* LRT Plius
* LRT Plius (D)
* Lietuvos rytas
* Lietuvos rytas (D)
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

Atsisiunčiate naujausią binary iš [releases](https://github.com/erkexzcx/lietuviska-tv/releases). Tuomet programą paleidžiame terminale:
```
# Kad galėtume ant Linux/FreeBSD executinti (Windowsams nereikia šios komandos):
chmod +x lietuviskatv_<platform>_<architecture>

# Paleidimas (Windowsuose naudokit Powershell):
./lietuviskatv_<platform>_<architecture>
```
Ir tuomet IPTV playlist pasiekiamas per šią nuorodą: `http://<address>:8989/iptv` (jei ant to paties kompiuterio: `http://127.0.0.1:8989/iptv`). Šią nuorodą naudokit ant VLC arba Kodi su *Simple IPTV addon*.

P.S. Linux SystemD service sukursiu ateityje. Šiuo metu patariu naudoti `tmux` ir palikti veikti background'e.

# FAQ

## Ką reiškia "(D)"

Kai kurie kanalai turi nekintančią transliacijos nuorodą, o kiti - dinaminę (ji nuolat kinta). **(D)** reiškia, kad nuoroda yra dinaminė ir toks kanalas gali ne visada veikti. Pamėginkite jį įsijungti vėliau - turėtų rodyti.

## Kaip ši programa veikia?

Programa atlieka kelias funkcijas:
1. Pateikia `M3U` IPTV kanalų playlist.
2. Sugeneruoja IPTV kanalų URL (kai kurie yra kintantys arba ne visada galimi).
3. Veikia kaip tarpininkas (savotiškas proxy serveris) tarp IPTV kliento ir prieinamų IPTV kanalų adresų. Visas IPTV srautas keliauja per šią programą.

Antras punktas labiausiai atspindi šios programos esmę. :)

Šiuo metu esu radęs apie 30 neapribotų IPTV kanalų, kurių `M3U8` adresas nekinta. Visi tie adresai yra įvesti programos kode (AKA hardcoded) ir jie tiesiog "yra". Likę IPTV kanalai (su **(D)** ženklu) yra nuolat kintantys, ir juos reikia išgauti programos pagalba. Pavyzdžiui kanalas *LNK HD (D)* yra internetu prieinamas tik tuomet, kai yra gyvai rodomos LNK vakaro žinios, todėl ši programa būtent tada išgauna rodomos transliacijos nuorodą ir ją laiko atmintyje, o kai nuorodos internete nebelieka - programa ją paima iš savo atminties. Tokia nuoroda dažniausiai veikia iki kitos dienos pietų.

Galbūt pastebėjote, kad užkrovus `<address>:8989/iptv` visų kanalų nuorodos yra adresuotos į tokį patį adresą, kuriuo yra pasiekiama programa (paminėkim žodį *proxy*). Tai yra dėl keletos priežąsčių:
1. Programa kas 10-15 min atnaujina dinaminių kanalų nuorodas, tačiau Kodi to nežino, kad IPTV kanalo nuoroda atsinaujino (ir niekada nesužino - tiesiog taip veikia). Kad priversti Kodi sužinoti naują nuorodą, reikia perkrauti arba pareloadinti *Simple IPTV addon*, antraip kanalas po kiek laiko nebus rodomas. Su šio proxinimo pagalba (URL perrašymu), iš Kodi perspektyvos, IPTV kanalo nuoroda yra visada vienoda, o pati programa ją fone nuolat atnaujina.
2. Kodi nenorėjo normaliai elgtis su pilna TV3 nuoroda - per Kodi net nerodydavo TV3 kanalo, o telefone rodydavo. Su šio proxinimo pagalba (URL perrašymu), šios problemos nebeliko.
3. Sugeneruota dinaminio kanalo nuoroda būna su tam tikru sesijos ID, kurį sugeneruoja pats IPTV kanalo serveris. Kiekvieną kart, kai IPTV klientas kreipiasi į serverį, yra tikrinama ir prie sesijos ID pririštas public IP (jeigu jis kitoks - kanalas nerodomas ir gaunamas HTTP 403 error). Su šio proxinimo pagalba (URL perrašymu), IPTV kanalo serverio pusės matomas IP adresas visada bus vienodas, todėl dabar šią programą galima talpinti išoriniame serveryje (pvz Google cloud).

## Kai kurių kanalų nerodo

Jei prie kanalo nėra prirašyta **(D)** - greičiausiai ir nerodys. Pamėgink vėliau - gal pradės.

## Ant VLC atsilieka garsas

Gali būt, kad ant VLC atsilieka **TV3** kanalo garsas. Jei taip nutiko - ant VLC reik spaust dešinį pelės klavišą (ant rodomo video) --> `tools` --> `Track synchronization` --> `Audio track synchronization` ir pakeisti į `-1`.

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

Jei norite pasileisti ant platformos ar architektūros, kurios nėra pateiktuose binaries - reikia pačiam sukūrti binary. Ant Linux įsirašykite Go ([taip](https://golang.org/doc/install) arba [taip](https://www.digitalocean.com/community/tutorials/how-to-install-go-on-debian-9), nes official repos esanti versija yra per sena), atsisiųskite šį projektą ir tada (pavyzdžiui OpenWRT naudojamai `MIPS` `softfloat`):
```
env GOOS=linux GOARCH=mips GOMIPS=softfloat go build -ldflags="-s -w" -o "lietuviskatv_linux_mips_softfloat" src/*.go
upx --best "lietuviskatv_linux_mips_softfloat" # Daugiau nei per pusę sumažina sukompiliuoto binary dydį
```
Daugiau informacijos apie galimas platformas ir architektūras: https://golang.org/doc/install/source#environment

## Žinau kanalą, kurį galima žiūrėti internetu, tačiau jo nėra tavo programoje

Pakelk naują issue šiam projektui surašydamas visas detales kur kas ir kaip. Pridėsiu į projektą.
