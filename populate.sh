 #!/bash/zsh

 echo $RANDOM

 for shard in localhost:8080 localhost:8081; do 
    echo $RANDOM
    for i in {1..100}; do 
          curl "http://$shard/set?key=key-$RANDOM&value=value-$RANDOM"
    done
    
done 

# curl "http://localhost:8080/get?key=key-16389"    - shard 0
# curl "http://localhost:8081/get?key=key-17552"    - shard 1
# curl "http://localhost:8080/get?key=key-29338"   - shard 0
# curl "http://localhost:8081/get?key=key-16324"    - shard 1
