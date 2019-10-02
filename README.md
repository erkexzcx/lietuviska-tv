# lietuviska-tv

Ši programa sugeneruoja šių kanalų `m3u` (_IPTV playlist_) compatible playlist'ą, kurį galima naudoti Kodi/VLC programose ir taip nemokamai žiūrėti lietuvišką TV internetu.

* LNK HD
* TV3 HD
* INFO TV HD
* Lietuvos rytas HD
* LRT HD
* LRT Plius HD

Viskas imama iš viešai prieinamų stream'ų internetu. Iš pačių LNK ir TV3 rankų, su žiupsneliu hackų :)

# Instaliacija

1. Perskaitote [#Troubleshooting](#Troubleshooting) ir įsitikinate, kad viskas jus tenkina.
2. Atsisiunčiate sukompiliuotą executable iš [releases](https://github.com/erkexzcx/lietuviska-tv/releases).
3. Pasileidžiate programą. Kad nereiktų services rašyt, aš ant RPI naudoju `nohup ./livetv-linux-arm &`.
4. IPTV playlistas bus pasiekiamas per `http://<ipaddress>:8989/iptv`

# Troubleshooting

## Nieko nerodo (arba Error 403)

Veikia tik iš to paties public IP. Tai reiškia, kad tiek serveris, tiek jūs turite būti po tuo pačiu public IP adresu.

## TV3 atsilieka garsas

Nežinau kas konkrečiai dėl to kaltas, nes šitas issue galioja tik VLC playeriui (kiek bandžiau). Kodi neturi šitos problemos. Ant VLC reik spaust dešinį pelės klavišą ant rodomo video --> `tools` --> `Track synchronization` --> `Audio track synchronization` ir pakeisk į `-1`.

## Nerodo LNK

LNK rodo tokiu principu - kai kuriuo nors metu per LNK GO rodo kokią nors tiesioginią transliaciją, scriptas tuo metu gali nuparsinti `M3U8` nuorodą per kurią galima žiūrėti LNK TV. Kai transliacija baigiasi - nuorodos nebelieka, tačiau scriptas iš atminties išsitraukia paskutinę galimą nuorodą kuri dažniausiai veikia iki kitos dienos pietų.

## Neveikia ant Windows

Netestavau. Windows sucks - naudok Linux. Sukompilinau Windowsams, nes buvo tokia galimybė.
