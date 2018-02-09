# * * * * * ec2-user cd /home/ec2-user/socket-lobby && ./cronjob.sh

git pull
git rev-parse HEAD > lobby-server/tmp/git.log
status=`curl -s -o /dev/null -w "%{http_code}" localhost:5110/api/git`
if ! [[ $status == "200" ]]
then
  make prod-bg
fi
