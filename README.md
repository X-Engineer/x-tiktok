# x-tiktok
第五届字节跳动后端青训营——抖音项目

# Contributing

1. `fork x-tiktok to your own namespace`
2. `git clone https://github.com/<your-username>/x-tiktok.git`
3. `git remote add upstream https://github.com/X-Engineer/x-tiktok.git`
4. `git fetch --all`
5. `git checkout -b <your-local-branch-name>`
6. `do some changes`
7. `git add .`
8. `git commit -m "your commit message"`
9. `git rebase origin/main (optional)`
10. `git push origin <your-branch-name>`
11. `create a pull request`
12. `wait for code review`
13. `merge your pull request`
14. `delete <your-branch-name>`
15. `git checkout main`
16. `git pull upstream main`
17. `Go back to step 5 to develop new features`

# deploy
1. `git clone https://github.com/X-Engineer/x-tiktok`
2. `cd x-tiktok`
3. config your database in `dao/db.go`
4. config your redis in `middleware/redis/redis.go`
5. config your rabbitmq in `middleware/rabbitmq/rabbitmq.go`
6. config your oss and jwt-token in `config/config.go`
7. sh run.sh