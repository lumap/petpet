clear
echo "Starting..."

docker run \
--rm \
-v ./src/:/app/src/ \
-p 8080:8080 \
$(docker build . -q)