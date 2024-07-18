for i in {1..100}; do
    curl "http://localhost:8081/get?key=key-$i"
done