echo "Backing up..."
cd ~/go/src/github.com/Nv7-Github/Nv7haven/data
git add -A
git commit -m "Backup $(date +%m/%d/%Y)"
git push origin main