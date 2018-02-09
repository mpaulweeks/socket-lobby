# * * * * * ec2-user cd /home/ec2-user/socket-lobby && ./bash/cronjob.sh

git pull
status=`curl -s -o /dev/null -w "%{http_code}" -X POST localhost:5110/api/git`
echo $status
 if ! [[ $status == "200" ]]
 then
  echo 'restarting...'
  sleep 2
  ./bash/bg_socket.sh
fi
