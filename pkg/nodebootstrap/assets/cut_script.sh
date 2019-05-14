#!/usr/bin/env bash
max_pods="";



MAX_PODS=`get_max_pods`


echo $MAX_PODS

grep three /home/martina/kkkk | cut -f 2 -d ' '
function get_max_pods() {
grep three /home/martina/kkkk | while read type pods; do

  if  [[ "$pods" =~ ^[0-9]+$ ]] ; then
    echo $pods;
    return
  fi ;
done
}


