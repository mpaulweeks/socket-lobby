# * * * * * ec2-user cd /home/ec2-user/socket-lobby && ./bash/cronjob.sh

oldCommit=`cat lobby-server/tmp/git.log`
git pull
git log -n 1 --pretty=format:"%H" > lobby-server/tmp/git.log
newCommit=`cat lobby-server/tmp/git.log`
if ! [[ $oldCommit == $newCommit ]]
then
  curl http://localhost:5110/api/git
  sleep 2
  ./bash/bg_socket.sh
fi
