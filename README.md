# lietuviska-tv

Tai yra tiesiog scriptas, kuris sugeneruoja šių kanalų `m3u` (_IPTV playlist_) compatible playlist'ą:
* LNK HD
* TV3 HD
* INFO TV HD
* Lietuvos rytas HD

Viskas imama iš viešai prieinamų live streamų internetu. Iš pačių LNK ir TV3 rankų, su žiupsneliu hackų :)

# Instaliacija

1. Susitvarkot servą ir kad jame veiktų `http` su `php`.
2. Susitvarkot, kad veiktų `go` komanda. Dazniausiai uztenka isirasyti `go` package.
3. Susitvarkot, kad server veiktų `crontab`'as.
4. Paleidziat zemiau komandas:
```
git clone https://github.com/erkexzcx/lietuviska-tv.git
sudo mv lietuviska-tv /opt/livetv
sudo chown -R <tavo_useris>:<http_useris> /opt/livetv
cd /opt/livetv
go build livetv.go
./livetv generate
```
5. Per komandą `EDITOR=vim crontab -e` prirašom tokią eilutę:
```
*/5 * * * * cd /opt/livetv && ./livetv generate
```
6. Pasidarot `index.php` su tokiu turiniu:
```
<?php

$output = shell_exec('cd /opt/livetv && ./livetv show');
echo $output;
```
7. Atidarot `index.php` per naršyklę ir jei matosi panašus vaizdas - viskas veikia!
```
#EXTM3U
#EXTINF:-1 group-title="LT" tvg-id="" tvg-logo="https://www.telia.lt/documents/20184/3686852/tv3-on-white.png", TV3 HD
https://cdn7.tvplayhome.lt/live/eds/TV3_LT_HD/HLS_encr/TV3_LT_HD.m3u8

#EXTM3U
#EXTINF:-1 group-title="LT" tvg-id="" tvg-logo="https://www.telia.lt/documents/20184/3686852/INFO-LOGO-HD.png", INFO TV HD
https://live.lnk.lt/lnk_live/tiesiogiai/playlist.m3u8?tokenstarttime=0&tokenendtime=1565079301&tokenhash=ZM9LPePBOPv3wc7lipJfJU5IB6H_fhmHajoSb9rfY8q6RyTAPHYp4Guoz-fgVV_7fB4M-le2oKQPTTQtxtDVng==

#EXTM3U
#EXTINF:-1 group-title="LT" tvg-id="" tvg-logo="https://www.telia.lt/documents/20184/3686852/LNK-LOGO-HD.png", LNK HD
https://live.lnk.lt/lnk_live/lnk/playlist.m3u8?tokenstarttime=0&tokenendtime=1565075957&tokenhash=tkyceQAZZYmNTbc2u_oIo-dxz1g5jmM8EOGGgTMsNjfBheoyegxSxo2fKMPkECHnQ5Gna-5KZooOsftDmMUuIw==
```
8. KODI --> _Simple IPTV_ plugine nurodote savo server pakurto PHP failo URL. EPG galima susigeneruoti patiems su [WebGrabber++](http://www.webgrabplus.com/).

Jeigu kas neveikia arba niekaip neišeina pasidaryti - kelk naują [issue](https://github.com/erkexzcx/lietuviska-tv/issues) :)
