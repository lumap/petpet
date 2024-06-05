clear
echo "Starting..."

docker run \
--rm \
--env-file ./.env.prod \
-v ./src/:/app/src/ \
-p 8080:8080 \
$(docker build . -q)