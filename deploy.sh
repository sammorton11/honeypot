set -e

SERVER="server@address"

echo "Syncing files to the server..."
rsync -avz --exclude '.git' ./ SERVER:~/honeypot/

echo "Building and starting containers..."
ssh $SERVER "cd ~/honeypot && docker-compose up --build -d"

echo "Done!"