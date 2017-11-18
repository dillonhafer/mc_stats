set -e

source .env.deploy
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m' # No Color

function notice() {
  printf "${GREEN}---->${NC} $1\n"
}

echo "Starting Full Deploy to $server"
notice 'Building server'

GOOS=linux GOARCH=386 go build

notice 'Compressing build'
tar czf mc_stats.tar.gz mc_stats

notice 'Builing frontend'
cd frontend && yarn build > /dev/null && cd ..

notice 'Compressing frontend'
tar czf frontend.tar.gz frontend/build

notice 'Uploading to server'
scp -q mc_stats.tar.gz frontend.tar.gz $server:$deploy_dir

notice 'Installing server'
ssh $server "cd $deploy_dir && tar xzf mc_stats.tar.gz && tar xzf frontend.tar.gz && rm mc_stats.tar.gz && rm frontend.tar.gz"

notice 'Cleaning up'
rm mc_stats.tar.gz frontend.tar.gz

notice 'Restarting server'
ssh $server 'systemctl restart mc-stats.service'

echo '      ✨  Done.'