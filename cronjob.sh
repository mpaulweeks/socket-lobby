git pull
git rev-parse HEAD > lobby-server/tmp/git.log
curl -s -o /dev/null -w "%{http_code}" localhost:5110/api/git
make prod-bg
