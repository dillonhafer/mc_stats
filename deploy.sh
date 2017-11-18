set -e

source .env.deploy

echo 'Building server'

GOOS=linux GOARCH=386 go build

echo 'Compressing build'
tar czf mc_stats.tar.gz mc_stats

echo 'Builing frontend'
cd frontend && yarn build && cd ..

echo 'Compressing frontend'
tar czf frontend.tar.gz frontend/build

echo 'Uploading to server'
scp mc_stats.tar.gz frontend.tar.gz $server:$deploy_dir

echo 'Installing server'
ssh $server "cd $deploy_dir && tar xzf mc_stats.tar.gz && tar xzf frontend.tar.gz && rm mc_stats.tar.gz && rm frontend.tar.gz"

echo 'Cleaning up'
rm mc_stats.tar.gz frontend.tar.gz

echo 'Restarting server'
ssh $server 'systemctl restart mc-stats.service'