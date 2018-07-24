#!/usr/bin/env bash
DBG="echo"
SRC="src"
VDR="vendor"
S=${PWD##*\/$SRC\/}
V=${PWD##*\/$VDR\/}
SL="${#S}"
VL="${#V}"
if [[ $SL -le $VL ]]; then
  PKG="$S"
else
  PKG="$V"
fi
PD=$(($(grep -o "/" <<< $PKG | wc -l | awk '{print $1}')+1))
if [[ `echo "$PKG" | cut -b 1` == '/' ]]; then
  echo "N/A"
  exit
fi
SPR="\.\.\/"
THS="\.\/"
FS=`ack -l "$SPR"`
echo "$FS" | while read f; do
  if [[ x$FS == x ]]; then
    echo "N|L[../]"
    break
  fi
  SDN="${f%\/*}"
  DIR="$SDN/"
  SFN="${f##*\/}"
  FD=$(grep -o "/" <<< $DIR | wc -l | awk '{print $1}')
  FP=$(dirname $DIR)
  KEY=`seq -f "$SPR" $FD | tr -d '\n'`
  VAL="$PKG/"
  VAL=`echo "$VAL" | sed -e 's/\./\\\\\./g' -e 's/\//\\\\\//g'`
  CMD=$(echo sed 's/"'$KEY'"/"'$VAL'"/g' "$f")
  RNP=~/.Trash/rnp.txt
  envbash "$RNP" && echo "$CMD > $f.rnp" >> "$RNP" && "$RNP" && rm "$RNP" && mv "$f.rnp" "$f"
done
FS=`ack -l "$THS"`
echo "$FS" | while read f; do
  if [[ x$FS == x ]]; then
    echo "N|L[./]"
    break
  fi
  SDN="${f%\/*}"
  DIR="$SDN/"
  SFN="${f##*\/}"
  FD=$(grep -o "/" <<< $DIR | wc -l | awk '{print $1}')
  FP=$(dirname $DIR)
  KEY=`seq -f "$THS" $FD | tr -d '\n'`
  VAL="$PKG/"
  VAL=`echo "$VAL" | sed -e 's/\./\\\\\./g' -e 's/\//\\\\\//g'`
  CMD=$(echo sed 's/"'$KEY'"/"'$VAL'"/g' "$f")
  RNP=~/.Trash/rnp.txt
  envbash "$RNP" && echo "$CMD > $f.rnp" >> "$RNP" && "$RNP" && rm "$RNP" && mv "$f.rnp" "$f"
done